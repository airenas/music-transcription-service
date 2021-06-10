package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	l, err := NewCount(1, time.Second)
	assert.Nil(t, err)
	assert.NotNil(t, l)
}

func TestNew_Fail(t *testing.T) {
	l, err := NewCount(0, time.Second)
	assert.NotNil(t, err)
	assert.Nil(t, l)
}

func TestLimit(t *testing.T) {
	l, _ := NewCount(1, time.Second)
	assert.NotNil(t, l)
	cf, err := l.Acquire()
	assert.Nil(t, err)
	assert.NotNil(t, cf)
	assert.Equal(t, 1, len(l.limitCh))
	cf()
	assert.Equal(t, 0, len(l.limitCh))
}

func TestLimit_Timeout(t *testing.T) {
	l, _ := NewCount(1, time.Millisecond)
	assert.NotNil(t, l)
	cf, err := l.Acquire()
	assert.Nil(t, err)
	assert.NotNil(t, cf)
	cf1, err := l.Acquire()
	assert.Nil(t, cf1)
	assert.NotNil(t, err)
	cf()
	cf, err = l.Acquire()
	assert.Nil(t, err)
	assert.NotNil(t, cf)
	defer cf()
}
