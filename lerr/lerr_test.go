package lerr

import (
	"testing"

	"github.com/testify/assert"
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
	assert.Equal(t, "Should Err 1 time\nTestError", ctx.Error())
}
