package file

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testS *Saver
var testWr *testWC

func initTest(t *testing.T) {
	var err error
	testS, err = NewSaver("./")
	testWr = &testWC{}
	testS.createF = func(fn string) (io.WriteCloser, error) { return testWr, nil }
	assert.Nil(t, err)
}

func TestNew(t *testing.T) {
	initTest(t)
}

func TestNew_Fail(t *testing.T) {
	s, err := NewSaver("")
	assert.NotNil(t, err)
	assert.Nil(t, s)
}

func TestSave(t *testing.T) {
	initTest(t)
	f, err := testS.Save("f.wav", strings.NewReader("olia"))
	assert.Nil(t, err)
	assert.Equal(t, "f.wav", f)
	assert.Equal(t, "olia", testWr.b.String())
	assert.True(t, testWr.cl)
}

func TestSave_Fail(t *testing.T) {
	initTest(t)
	testS.createF = func(fn string) (io.WriteCloser, error) { return nil, errors.New("err") }
	f, err := testS.Save("f.wav", strings.NewReader("olia"))
	assert.NotNil(t, err)
	assert.Equal(t, "", f)
}

type testWC struct {
	b  strings.Builder
	cl bool
}

func (w *testWC) Close() error {
	w.cl = true
	return nil
}

func (w *testWC) Write(p []byte) (n int, err error) {
	return w.b.Write(p)
}
