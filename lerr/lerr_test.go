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

func TestNotEqual(t *testing.T) {
	a, b := 5, 5
	err := lerr.NewNotEqual(a == b, a, b)
	assert.NoError(t, err)
	b = 3
	err = lerr.NewNotEqual(a == b, a, b)
	assert.Equal(t, "Expected 5 got 3", err.Error())
}

func TestSliceErrs(t *testing.T) {
	a := []int{1, 2, 0, 4}
	b := []int{1, 2, 3}
	fn := func(i int) error {
		return lerr.NewNotEqual(a[i] == b[i], a[i], b[i])
	}

	err := lerr.NewSliceErrs(len(a), len(b), fn)

	assert.Equal(t, "Lengths do not match: Expected 4 got 3\n\t2: Expected 0 got 3", err.Error())

	a = []int{1, 2, 3}
	err = lerr.NewSliceErrs(len(a), len(b), fn)
	assert.NoError(t, err)

	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	b = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	err = lerr.NewSliceErrs(len(a), len(b), fn)
	assert.Equal(t, "Lengths do not match: Expected 11 got 12\n\t0: Expected 1 got 2\n\t1: Expected 2 got 3\n\t2: Expected 3 got 4\n\t3: Expected 4 got 5\n\t4: Expected 5 got 6\n\t5: Expected 6 got 7\n\t6: Expected 7 got 8\n\t7: Expected 8 got 9\n\t8: Expected 9 got 10\nOmitting 2 more", err.Error())

	se := lerr.SliceErrs{}
	se = se.AppendF(1, "%s is a test %d", "this", 123)
	assert.Equal(t, "\t1: this is a test 123", se.Error())

	a = []int{1, 2, 3}
	b = []int{1, 2, 3, 4, 5}
	err = lerr.NewSliceErrs(len(a), -1, fn)
	assert.NoError(t, err)
}
