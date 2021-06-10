package service

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/music-transcription-service/internal/pkg/limiter"
	"github.com/airenas/music-transcription-service/internal/pkg/utils"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type (
	// FileSaver saves the file with the provided name
	FileSaver interface {
		Save(name string, reader io.Reader) (string, error)
	}

	// Encoder encodes file, returns file name
	Transcriber interface {
		Convert(nameIn string) (string, error)
	}

	//Data is service operation data
	Data struct {
		Port int

		Saver   FileSaver
		Worker  Transcriber
		Limiter *limiter.Count

		readFunc func(string) ([]byte, error)
	}
)

//StartWebServer starts the HTTP service and listens for the convert requests
func StartWebServer(data *Data) error {
	goapp.Log.Infof("Starting HTTP music transcription service at %d", data.Port)
	portStr := strconv.Itoa(data.Port)
	data.readFunc = ioutil.ReadFile
	e := initRoutes(data)
	e.Server.Addr = ":" + portStr
	e.Server.ReadTimeout = 2 * time.Minute
	e.Server.WriteTimeout = 2 * time.Minute
	e.Server.IdleTimeout = 2 * time.Minute

	w := goapp.Log.Writer()
	defer w.Close()
	l := log.New(w, "", 0)
	gracehttp.SetLogger(l)

	return gracehttp.Serve(e.Server)
}

var promMdlw *prometheus.Prometheus

func init() {
	promMdlw = prometheus.NewPrometheus("mts", nil)
}

func initRoutes(data *Data) *echo.Echo {
	e := echo.New()
	promMdlw.Use(e)

	e.POST("/transcribe", transcribe(data))
	e.GET("/live", live(data))

	goapp.Log.Info("Routes:")
	for _, r := range e.Routes() {
		goapp.Log.Infof("  %s %s", r.Method, r.Path)
	}
	return e
}

type output struct {
	MusicXML string `json:"musicXML,omitempty"`
	Error    string `json:"error,omitempty"`
}

func transcribe(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		closef, err := data.Limiter.Acquire()
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "service too busy")
		}
		defer closef()

		defer goapp.Estimate("Service method")()

		form, err := c.MultipartForm()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "no multipart form data")
		}
		defer cleanFiles(form)

		files, ok := form.File["file"]
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "no file")
		}
		if len(files) > 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "multiple files")
		}

		file := files[0]
		ext := filepath.Ext(file.Filename)
		ext = strings.ToLower(ext)
		if !checkFileExtension(ext) {
			return echo.NewHTTPError(http.StatusBadRequest, "wrong file type: "+ext)
		}

		id := uuid.New().String()
		fileName := id + ext

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "can't read file")
		}
		defer src.Close()

		est := goapp.Estimate("Saving")
		fileNameIn, err := data.Saver.Save(fileName, src)
		if err != nil {
			goapp.Log.Error(err)
			return errors.Wrap(err, "can not save file")
		}
		defer deleteFile(fileNameIn)
		est()

		est = goapp.Estimate("Transcribe")
		fileNameOut, err := data.Worker.Convert(fileNameIn)
		res := &output{}
		if err != nil {
			var errTr *utils.ErrTranscribe
			if errors.As(err, &errTr) {
				res.Error = errTr.Msg
				return c.JSON(http.StatusOK, res)
			}
			goapp.Log.Error(err)
			return errors.Wrap(err, "can not transcribe file")
		}
		defer deleteFile(fileNameOut)
		est()

		fd, err := data.readFunc(fileNameOut)
		if err != nil {
			goapp.Log.Error(err)
			return errors.Wrap(err, "Can not read file")
		}

		res.MusicXML = base64.StdEncoding.EncodeToString(fd)

		return c.JSON(http.StatusOK, res)
	}
}

func checkFileExtension(ext string) bool {
	return ext == ".wav"
}

func deleteFile(file string) {
	os.RemoveAll(file)
}

func live(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{"service":"OK"}`))
	}
}

func cleanFiles(f *multipart.Form) {
	if f != nil {
		f.RemoveAll()
	}
}
