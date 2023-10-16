package handler_test

import (
	"sync"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestSwitchErrors(t *testing.T) {
	s := handler.NewSwitch(0)
	a, err := s.Handle("nothing is registered")
	assert.Nil(t, a)
	assert.Equal(t, handler.ErrNoHandler, err)

	// Only funcs and channels can be registered, and a func has to be one that
	// Handler allows.
	assert.Equal(t, handler.ErrRegisterInterface, s.RegisterInterface("string"))
	assert.Error(t, s.RegisterInterface(func(a, b int) {}))
}

func TestHandlers(t *testing.T) {
	s, err := handler.Handlers(
		func(name string) string { return "hello " + name },
		func(i int) int { return i * 2 },
	)
	assert.NoError(t, err)
	a, err := s.Handle("Ada")
	assert.NoError(t, err)
	assert.Equal(t, "hello Ada", a)
	a, err = s.Handle(21)
	assert.NoError(t, err)
	assert.Equal(t, 42, a)

	// It stops at the first handler that can't be registered.
	s, err = handler.Handlers(func(i int) int { return i }, 5, func(s string) {})
	assert.Equal(t, handler.ErrRegisterInterface, err)
	a, err = s.Handle(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, a)
	_, err = s.Handle("late")
	assert.Equal(t, handler.ErrNoHandler, err)
}

func TestSwitchNoArgument(t *testing.T) {
	s := handler.NewSwitch(1)
	assert.NoError(t, s.RegisterInterface(func() string { return "called" }))
	a, err := s.Handle(nil)
	assert.NoError(t, err)
	assert.Equal(t, "called", a)

	// There is only one, so a second replaces the first.
	assert.NoError(t, s.RegisterInterface(func() string { return "replaced" }))
	a, err = s.Handle(nil)
	assert.NoError(t, err)
	assert.Equal(t, "replaced", a)
}

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
}
