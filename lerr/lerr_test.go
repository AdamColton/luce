package lerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testErr = Str("TestError")

func TestError(t *testing.T) {
	assert.Error(t, testErr)
	assert.Equal(t, testErr.Error(), string(testErr))
}

func TestPanic(t *testing.T) {
	Panic(nil)

	defer func() {
		assert.Equal(t, recover(), testErr)
	}()

	Panic(testErr)
}

func TestCtx(t *testing.T) {
	ctx := Wrap(nil, "No Error")
	assert.NoError(t, ctx)

	ctx = Wrap(testErr, "Should Err %d time", 1)
	assert.Error(t, ctx)
	assert.Equal(t, "Should Err 1 time: TestError", ctx.Error())
}

func TestMany(t *testing.T) {
	var m Many
	assert.Nil(t, m.Get())

	err1 := Str("Error 1")
	m = m.Add(err1)
	assert.Equal(t, err1, m.Get())
	m = m.Add(nil)
	assert.Equal(t, err1, m.Get())
	m = m.Add(Str("Error 2"))
	assert.Equal(t, m, m.Get())

	assert.Equal(t, "Error 1\nError 2", m.Error())
}
