package lerr_test

import (
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
