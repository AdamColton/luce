package handler_test

import (
	"sync"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestSwitch(t *testing.T) {
	s := handler.NewSwitch(10)

	var hmi handler.Switcher = s
	assert.NotNil(t, hmi)

	strCh := make(chan string)
	err := s.RegisterInterface(func(s string) int {
		strCh <- s
		return 123
	})
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		assert.Equal(t, "test", <-strCh)
		wg.Done()
	}()

	a, err := s.Handle("test")
	assert.NoError(t, err)
	assert.Equal(t, 123, a)
	timeout.After(5, &wg)

	intCh := make(chan int)
	s.RegisterInterface(intCh)
	wg.Add(1)
	go func() {
		s.Handle(31415)
		wg.Done()
	}()
	err = timeout.After(30000, func() {
		assert.Equal(t, 31415, <-intCh)
	})
	assert.NoError(t, err)
	wg.Wait()

	testErr := lerr.Str("test error")
	s.RegisterInterface(func(b bool) error {
		return testErr
	})
	a, err = s.Handle(true)
	assert.Nil(t, a)
	assert.Equal(t, testErr, err)

	testErr = lerr.Str("multi return")
	err = s.RegisterInterface(func(s float64) (int, error) {
		if s > 0 {
			return 456, nil
		}
		return 789, testErr
	})
	assert.NoError(t, err)

	a, err = s.Handle(1.0)
	assert.Equal(t, 456, a)
	assert.NoError(t, err)

	a, err = s.Handle(-1.0)
	assert.Equal(t, testErr, err)
	assert.Equal(t, 789, a)

	// h, err := handler.New(func() string {
	// 	return "hello"
	// }, "sayHi")
	// assert.NoError(t, err)
	// s.RegisterHandler(h)

	// a, err = s.Handle("sayHi")
	// assert.NoError(t, err)
	// assert.Equal(t, "hello", a)
}
