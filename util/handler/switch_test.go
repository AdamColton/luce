package handler_test

import (
	"fmt"
	"strconv"
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

type handlerObj struct {
	name         string
	includeFloat bool
}

func (ho handlerObj) StringHandler(s string) string {
	return ho.name + s
}

func (ho handlerObj) StringUsage() string {
	return "combine ho.name and s"
}

func (ho handlerObj) IntHandler(i int) string {
	return strconv.Itoa(i)
}

func (ho handlerObj) Float64Handler(f float64) string {
	return fmt.Sprint(f)
}

func (ho handlerObj) Float64Usage() (string, bool) {
	return "convert float64 to string", ho.includeFloat
}

func (ho handlerObj) IAmNotAHandler(i int, s string) {
}

func TestMagic(t *testing.T) {
	hm := handler.NewSwitch(10)
	ho := handlerObj{
		name: "test",
	}

	handler.DefaultRegistrar.Register(hm, ho)

	a, err := hm.Handle(" foo")
	assert.NoError(t, err)
	assert.Equal(t, "test foo", a)

	cmds := handler.DefaultRegistrar.Commands(ho).
		Sort(func(i, j handler.Command) bool {
			return i.Name < j.Name
		})
	assert.Equal(t, "int", cmds[0].Name)
	assert.Equal(t, "string", cmds[1].Name)
	assert.Equal(t, ho.StringUsage(), cmds[1].Usage)
}
