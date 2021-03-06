package service

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/airenas/music-transcription-service/internal/pkg/limiter"
	"github.com/airenas/music-transcription-service/internal/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	tData  *Data
	tSaver *testSaver
	tCoder *testCoder
	tEcho  *echo.Echo
	tReq   *http.Request
	tRec   *httptest.ResponseRecorder
)

func initTest(t *testing.T) {
	tSaver = &testSaver{name: "test.wav"}
	tCoder = &testCoder{res: "olia"}
	tData = newTestData(tSaver, tCoder)
	tEcho = initRoutes(tData)
	tReq = newTestRequest("file.wav")
	tRec = httptest.NewRecorder()
}

func TestLive(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodGet, "/live", nil)

	e := initRoutes(tData)
	e.ServeHTTP(tRec, req)
	assert.Equal(t, http.StatusOK, tRec.Code)
	assert.Equal(t, `{"service":"OK"}`, tRec.Body.String())
}

func TestTranscribe(t *testing.T) {
	initTest(t)

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusOK, tRec.Code)
	assert.Equal(t, `{"musicXML":"dGVzdA=="}`+"\n", tRec.Body.String())
	assert.Equal(t, "flute", tCoder.ins)
}

func TestCTranscribe_FailData(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest("POST", "/transcription", strings.NewReader("aa"))

	tEcho.ServeHTTP(tRec, req)

	assert.Equal(t, http.StatusBadRequest, tRec.Code)
}

func TestCTranscribe_FailLimit(t *testing.T) {
	initTest(t)
	tData.Limiter, _ = limiter.NewCount(1, time.Millisecond)
	cf, _ := tData.Limiter.Acquire()
	defer cf()
	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusForbidden, tRec.Code)
}

func TestTranscribe_FailType(t *testing.T) {
	initTest(t)
	req := newTestRequest("a.txt")

	tEcho.ServeHTTP(tRec, req)

	assert.Equal(t, http.StatusBadRequest, tRec.Code)
}

func TestTranscribe_FailSaver(t *testing.T) {
	initTest(t)

	tSaver.err = errors.New("olia")

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

func TestTranscribe_FailConvert(t *testing.T) {
	initTest(t)

	tCoder.err = errors.New("olia")

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

func TestTranscribe_FailWithError(t *testing.T) {
	initTest(t)

	tCoder.err = utils.NewErrTranscribe("olia")

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusOK, tRec.Code)
	assert.Equal(t, `{"error":"olia"}`+"\n", tRec.Body.String())
}

func TestTranscribe_FailRead(t *testing.T) {
	initTest(t)

	tData.readFunc = func(string) ([]byte, error) { return nil, errors.New("olia") }

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

type testSaver struct {
	name string
	err  error
	data bytes.Buffer
}

func (s *testSaver) Save(name string, reader io.Reader) (string, error) {
	io.Copy(&s.data, reader)
	return s.name, s.err
}

type testCoder struct {
	err  error
	name string
	ins  string
	res  string
}

func (s *testCoder) Convert(name, ins string) (string, error) {
	s.name = name
	s.ins = ins
	return s.res, s.err
}

func newTestData(s FileSaver, e Transcriber) *Data {
	l, _ := limiter.NewCount(1, time.Second)
	return &Data{Saver: s, Worker: e, readFunc: func(string) ([]byte, error) { return []byte("test"), nil },
		Limiter: l}
}

func newTestRequest(file string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if file != "" {
		part, _ := writer.CreateFormFile("file", file)
		_, _ = io.Copy(part, strings.NewReader("body"))
	}
	writer.WriteField("instrument", "flute")
	writer.Close()
	req := httptest.NewRequest("POST", "/transcription", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
