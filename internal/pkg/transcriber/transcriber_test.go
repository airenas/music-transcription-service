package transcriber

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/airenas/music-transcription-service/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var testTr *Worker
var testParams []string

func initTest(t *testing.T) {
	var err error
	testParams = nil
	testTr, err = NewWorker("trApp {{INPUT}} {{OUTPUT}}")
	testTr.convertFunc = func(cmd []string) error {
		testParams = cmd
		return nil
	}
	assert.Nil(t, err)
}

func TestFile(t *testing.T) {
	initTest(t)
	testTr.convertFunc = func(cmd []string) error {
		testParams = cmd
		return errors.New("olia")
	}
	d, err := testTr.Convert("/dir/1.wav", "ins")
	assert.Equal(t, "", d)
	assert.NotNil(t, err)
}

func TestFile_Fail(t *testing.T) {
	initTest(t)
	d, err := testTr.Convert("/dir/1.wav", "ins")
	assert.Equal(t, []string{"trApp", "/dir/1.wav", "/dir/1.musicxml"}, testParams)
	assert.Equal(t, "/dir/1.musicxml", d)
	assert.Nil(t, err)
}

func TestMapError(t *testing.T) {
	err := mapError(errors.New("olia"), "stderr msg")
	assert.Equal(t, "stderr msg", err.Error())
	var te *utils.ErrTranscribe
	assert.True(t, errors.As(err, &te))
	err = mapError(&exec.ExitError{ProcessState: &os.ProcessState{}}, "")
	assert.Equal(t, "exit status 0", err.Error())
	assert.False(t, errors.As(err, &te))
}

func TestRunCmd(t *testing.T) {
	err := runCmd([]string{"ls", "-la"}, time.Second)
	assert.Nil(t, err)
	err = runCmd([]string{"badcmddd"}, time.Second)
	assert.NotNil(t, err)
	err = runCmd([]string{"bash", "-c", "echo aaa 1>&2; exit 1"}, time.Second)
	assert.Equal(t, "aaa\n", err.Error())
	var tErr *utils.ErrTranscribe
	assert.True(t, errors.As(err, &tErr))
	err = runCmd([]string{"bash", "-c", "echo aaa 1>&2; exit 0"}, time.Second)
	assert.Equal(t, "aaa\n", err.Error())
	assert.True(t, errors.As(err, &tErr))
	err = runCmd([]string{"bash", "-c", "echo aaa; exit 1"}, time.Second)
	assert.False(t, errors.As(err, &tErr))
}

func TestRunCmd_Timeout(t *testing.T) {
	err := runCmd([]string{"sleep", "1"}, time.Millisecond*50)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
}

func TestPrepareParams(t *testing.T) {
	assert.Equal(t, []string{"app"}, prepareParams("app", "1", "2", ""))
	assert.Equal(t, []string{"app", "1"}, prepareParams("app {{INPUT}}", "1", "2", ""))
	assert.Equal(t, []string{"app", "2", "1"}, prepareParams("app {{OUTPUT}} {{INPUT}}", "1", "2", ""))
	assert.Equal(t, []string{"app", "2", "1=2"}, prepareParams("app {{OUTPUT}} {{INPUT}}={{OUTPUT}}", "1", "2", ""))
	assert.Equal(t, []string{"app", "2", "1=2", "flute"}, prepareParams("app {{OUTPUT}} {{INPUT}}={{OUTPUT}} {{INSTRUMENT}}", "1", "2", "flute"))
}
