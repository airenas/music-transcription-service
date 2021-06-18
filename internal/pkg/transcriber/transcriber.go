package transcriber

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/music-transcription-service/internal/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

//Worker wrapper for file conversion
type Worker struct {
	cmdPath     string
	convertFunc func([]string) error
}

//NewWorker return new transcribe wrapper
func NewWorker(cmdPath string) (*Worker, error) {
	res := Worker{}
	if cmdPath == "" {
		return nil, errors.New("no cmd path")
	}
	res.cmdPath = cmdPath
	res.convertFunc = func(p []string) error { return runCmd(p, time.Minute*2) }
	goapp.Log.Infof("Cmd: %s", cmdPath)
	return &res, nil
}

//Convert returns name of new converted file
func (e *Worker) Convert(nameIn, instrument string) (string, error) {
	resName := getNewFile(nameIn)
	params := prepareParams(e.cmdPath, nameIn, resName, instrument)
	err := e.convertFunc(params)
	if err != nil {
		return "", err
	}
	return resName, nil
}

func runCmd(cmdArr []string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	goapp.Log.Infof("Cmd: %s", strings.Join(cmdArr, " "))
	cmd := exec.CommandContext(ctx, cmdArr[0], cmdArr[1:]...)
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	resChan := make(chan error, 1)
	go func() {
		resChan <- cmd.Run()
	}()

	var err error
	select {
	case e := <-resChan:
		err = e
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), fmt.Sprintf("Timeout. Out: %s\nErr: %s", outBuffer.String(), errBuffer.String()))
	}

	es := errBuffer.String()
	if err != nil || es != "" {
		goapp.Log.Errorf("Cmd err: %s\nout %s:", es, outBuffer.String())
		return mapError(err, es)
	}
	return nil
}

func getNewFile(file string) string {
	f := filepath.Base(file)
	ext := filepath.Ext(f)
	d := filepath.Dir(file)
	return filepath.Join(d, fmt.Sprintf("%s.%s", f[:len(f)-len(ext)], "musicxml"))
}

func mapError(err error, es string) error {
	if es != "" {
		return utils.NewErrTranscribe(es)
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		c := exitErr.ExitCode()
		if c == -1 {
			if exitErr.ProcessState != nil {
				s := exitErr.ProcessState.String()
				goapp.Log.Error(s)
				return errors.Wrapf(err, "Err: %s", s)
			}
		}
	}
	return err
}

func prepareParams(cmd, fIn, fOut, ins string) []string {
	res := []string{}
	iCmd := strings.ReplaceAll(cmd, "{{INPUT}}", fIn)
	iCmd = strings.ReplaceAll(iCmd, "{{OUTPUT}}", fOut)
	iCmd = strings.ReplaceAll(iCmd, "{{INSTRUMENT}}", ins)
	res = append(res, strings.Split(iCmd, " ")...)
	return res
}
