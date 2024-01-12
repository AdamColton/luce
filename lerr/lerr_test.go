package lerr_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

const testErr = lerr.Str("TestError")

func TestError(t *testing.T) {
	assert.Error(t, testErr)
	assert.Equal(t, testErr.Error(), string(testErr))
}

func ExampleStr() {
	const ErrExample = lerr.Str("example error")
	fmt.Println(ErrExample)
	// Output: example error
}

func TestPanic(t *testing.T) {
	lerr.Panic(nil)
	lerr.Panic(io.EOF, io.EOF)

	defer func() {
		assert.Equal(t, recover(), testErr)
	}()

	lerr.Panic(testErr)
}

func ExamplePanic() {
	var err error

	// won't panic on nil error
	lerr.Panic(err)

	err = io.EOF
	// won't panic when err in except args
	lerr.Panic(err, io.EOF)

	defer func() {
		fmt.Println(recover())
	}()

	err = lerr.Str("this will panic")
	lerr.Panic(err)

	// Output: this will panic
}

func TestMust(t *testing.T) {
	noErr := func() (int, error) {
		return 5, nil
	}
	hasErr := func() (int, error) {
		return 10, testErr
	}

	i := lerr.Must(noErr())
	fmt.Println("No error, got:", i)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic:", r)
		}
	}()

	i = lerr.Must(hasErr())
	fmt.Println("This line is never reached", i)
	// Output: No error, got: 5
	// Panic: TestError
}

func ExampleMust() {

}

func TestWrap(t *testing.T) {
	w := lerr.Wrap(nil, "No Error")
	assert.NoError(t, w)

	w = lerr.Wrap(testErr, "Should Err %d time", 1)
	assert.Error(t, w)
	assert.Equal(t, "Should Err 1 time: TestError", w.Error())
}

func ExampleWrap() {
	w := lerr.Wrap(nil, "No Error")
	fmt.Println(w)

	innerError := lerr.Str("TestError")
	w = lerr.Wrap(innerError, "Should Err %d time", 1)
	fmt.Println(w)
	// Output:
	// <nil>
	// Should Err 1 time: TestError
}

func TestMany(t *testing.T) {
	var m lerr.Many
	assert.NoError(t, m.Cast())
	m = m.Add(lerr.Str("Error 1"))
	m = m.Add(nil)
	m = m.Add(lerr.Str("Error 2"))

	assert.Equal(t, "Error 1\nError 2", m.Cast().Error())

	m = m[:0]
	assert.NoError(t, m.Cast())
}

func ExampleMany() {
	var m lerr.Many
	var err error

	// nil error not added to Many
	m = m.Add(err)
	// <nil>
	fmt.Println(m.Cast())

	fmt.Println("---")

	// when many contains a single error, only that is returned from cast
	err = lerr.Str("first error")
	m = m.Add(err)
	fmt.Println(m.Cast())

	fmt.Println("---")

	err = lerr.Str("second error")
	m = m.Add(err)
	fmt.Println(m.Cast())

	// Output:
	// <nil>
	// ---
	// first error
	// ---
	// first error
	// second error
}

func TestHandlerFunc(t *testing.T) {

	var called error
	fn := func(err error) {
		called = err
	}
	got, err := lerr.HandlerFunc(fn)
	assert.NoError(t, err)
	expected := lerr.Str("test1")
	assert.True(t, got.Handle(expected))
	assert.Equal(t, expected, called)

	assert.False(t, got.Handle(nil))

	ch := make(chan error, 1)
	got, err = lerr.HandlerFunc(ch)
	assert.NoError(t, err)
	expected = lerr.Str("test2")
	assert.True(t, got.Handle(expected))
	assert.Equal(t, expected, <-ch)

	var chIn chan<- error = ch
	got, err = lerr.HandlerFunc(chIn)
	assert.NoError(t, err)
	expected = lerr.Str("test3")
	assert.True(t, got.Handle(expected))
	assert.Equal(t, expected, <-ch)

	got, err = lerr.HandlerFunc(func(string) {})
	assert.Equal(t, lerr.ErrHandlerFunc, err)
	assert.Nil(t, got)

	got, err = lerr.HandlerFunc(nil)
	assert.Nil(t, got)
	assert.Nil(t, err)
	got.Handle(testErr) // Make sure this doesn't panic
}

func TestLog(t *testing.T) {
	restore := lerr.LogTo
	defer func() {
		lerr.LogTo = restore
	}()

	buf := bytes.NewBuffer(nil)
	lerr.LogTo = func(err error) {
		if err != nil {
			buf.WriteString(err.Error())
		}
	}

	te := lerr.Str("test error")
	assert.True(t, lerr.Log(te))
	assert.Equal(t, te.Error(), buf.String())

	assert.False(t, lerr.Log(te, te))
}

func TestManyFirst(t *testing.T) {
	te := lerr.Str("test error")
	to := lerr.Str("test other")
	err := lerr.NewMany(nil, te, nil, to, nil).First()
	assert.Equal(t, te, err)

	err = lerr.NewMany(nil, nil, nil).First()
	assert.NoError(t, err)
}

func TestBadWhence(t *testing.T) {
	err := lerr.Whence(0)
	assert.NoError(t, err)

	err = lerr.Whence(-1)
	assert.Error(t, err)

	err = lerr.Whence(3)
	assert.Equal(t, "lerr: Seek whence value should be io.SeekStart (0), io.SeekCurrent (1) or io.SeekEnd (2), got:3", err.Error())
}
