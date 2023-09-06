package lerr_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

const testErr = lerr.Str("TestError")

func TestError(t *testing.T) {
	assert.Error(t, testErr)
	assert.Equal(t, testErr.Error(), string(testErr))
}

func TestPanic(t *testing.T) {
	lerr.Panic(nil)

	lerr.Panic(testErr, testErr)

	defer func() {
		assert.Equal(t, recover(), testErr)
	}()

	lerr.Panic(testErr)
}

func TestCtx(t *testing.T) {
	ctx := lerr.Wrap(nil, "No Error")
	assert.NoError(t, ctx)

	ctx = lerr.Wrap(testErr, "Should Err %d time", 1)
	assert.Error(t, ctx)
	assert.Equal(t, "Should Err 1 time: TestError", ctx.Error())
}

func TestMany(t *testing.T) {
	var m lerr.Many
	assert.Nil(t, m.Get())

	err1 := lerr.Str("Error 1")
	m = m.Add(err1)
	assert.Equal(t, err1, m.Get())
	m = m.Add(nil)
	assert.Equal(t, err1, m.Get())
	m = m.Add(lerr.Str("Error 2"))
	assert.Equal(t, m, m.Get())

	assert.Equal(t, "Error 1\nError 2", m.Error())
}

func TestLog(t *testing.T) {
	assert.False(t, lerr.Log(testErr, testErr))
	called := false
	lerr.LogTo = func(err error) {
		assert.Equal(t, testErr, testErr)
		called = true
	}

	assert.True(t, lerr.Log(testErr))
	assert.True(t, called)
}
func TestTypeMismatch(t *testing.T) {
	expected := `Types do not match: expected "string", got "float64"`
	err := lerr.NewTypeMismatch("test", 1.0)
	assert.Equal(t, expected, err.Error())

	err = lerr.TypeMismatch("test", 1.0)
	assert.Equal(t, expected, err.Error())

	err = lerr.NewTypeMismatch("test", "foo")
	assert.NoError(t, err)
}

func TestLenMismatch(t *testing.T) {
	min, max, err := lerr.NewLenMismatch(5, 5)
	assert.Equal(t, min, 5)
	assert.Equal(t, max, 5)
	assert.NoError(t, err)

	min, max, err = lerr.NewLenMismatch(3, 5)
	assert.Equal(t, min, 3)
	assert.Equal(t, max, 5)
	assert.Equal(t, "Lengths do not match: Expected 3 got 5", err.Error())

	min, max, err = lerr.NewLenMismatch(5, 3)
	assert.Equal(t, min, 3)
	assert.Equal(t, max, 5)
	assert.Equal(t, "Lengths do not match: Expected 5 got 3", err.Error())
}
