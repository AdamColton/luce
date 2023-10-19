package timeout_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestMockTimeMessy(t *testing.T) {
	mt := timeout.NewMockTime()

	i := 1000
	wg := sync.WaitGroup{}
	wg.Add(1000)

	a := atomic.Int32{}
	a.Add(1000)
	go func() {
		for ; i > 0; i-- {
			time.Sleep(time.Microsecond)
			go func() {
				mt.Sleep(time.Millisecond * 100)
				wg.Done()
			}()
		}
	}()

	done := false
	go func() {
		for !done {
			time.Sleep(time.Microsecond * 100)
			go func() {
				mt.Tick(1)
			}()
		}
	}()

	err := timeout.After(10000, &wg)
	done = true
	assert.NoError(t, err)

}

func TestMockTimeClean(t *testing.T) {
	mt := timeout.NewMockTime()

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			mt.Sleep(100)
			wg.Done()
		}()
	}

	time.Sleep(time.Millisecond * 10) // give goroutines a chance to catch up

	mt.Tick(101)

	timeout.After(10, &wg)
}
