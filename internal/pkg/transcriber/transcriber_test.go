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
	d, err := testTr.Convert("/dir/1.wav")
	assert.Equal(t, "", d)
	assert.NotNil(t, err)
}

func TestFile_Fail(t *testing.T) {
	initTest(t)
	d, err := testTr.Convert("/dir/1.wav")
	assert.Equal(t, []string{"trApp", "/dir/1.wav", "/dir/1.musicxml"}, testParams)
	assert.Equal(t, "/dir/1.musicxml", d)
	assert.Nil(t, err)
}

func TestMapError(t *testing.T) {
	err := mapError(errors.New("olia"), func() string { return "err" })
	assert.Equal(t, "Output: err: olia", err.Error())

	err = mapError(&exec.ExitError{ProcessState: &os.ProcessState{}}, func() string { return "err" })
	assert.Equal(t, "Some other error", err.Error())
}

func TestRunCmd(t *testing.T) {
	err := runCmd([]string{"ls", "-la"}, time.Second)
	assert.Nil(t, err)
	err = runCmd([]string{"badcmddd"}, time.Second)
	assert.NotNil(t, err)
	err = runCmd([]string{"bash", "-c", "exit 1"}, time.Second)
	assert.Equal(t, "Error 1", err.Error())
	var tErr *utils.ErrTranscribe
	assert.True(t, errors.As(err, &tErr))
}

func TestRunCmd_Timeout(t *testing.T) {
	err := runCmd([]string{"sleep", "1"}, time.Millisecond*50)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
}

func TestPrepareParams(t *testing.T) {
	assert.Equal(t, []string{"app"}, prepareParams("app", "1", "2"))
	assert.Equal(t, []string{"app", "1"}, prepareParams("app {{INPUT}}", "1", "2"))
	assert.Equal(t, []string{"app", "2", "1"}, prepareParams("app {{OUTPUT}} {{INPUT}}", "1", "2"))
	assert.Equal(t, []string{"app", "2", "1=2"}, prepareParams("app {{OUTPUT}} {{INPUT}}={{OUTPUT}}", "1", "2"))
}

