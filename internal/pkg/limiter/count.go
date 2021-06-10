package limiter

import (
	"time"

	"github.com/pkg/errors"
)

var ErrAcquireTimeout = errors.New("limiter count acquire timeout")

//Count limiter struct
type Count struct {
	limitCh     chan bool
	waitTimeout time.Duration
}

//NewCount creates new count limiter
func NewCount(max int, waitTmeout time.Duration) (*Count, error) {
	res := Count{}
	res.waitTimeout = waitTmeout
	if max < 1 || max > 100 {
		return nil, errors.New("max must be in [1, 100]")
	}
	res.limitCh = make(chan bool, max)
	return &res, nil
}

//Acquire gets the access
func (l *Count) Acquire() (func(), error) {
	select {
	case <-time.After(l.waitTimeout):
		return nil, ErrAcquireTimeout
	case l.limitCh <- true:
		return l.returnLock, nil
	}
}

func (l *Count) returnLock() {
	<-l.limitCh
}
