package lerr_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestLenMismatch(t *testing.T) {
	min, max, err := lerr.NewLenMismatch(15, 10)
	assert.Equal(t, 10, min)
	assert.Equal(t, 15, max)
	assert.Equal(t, "Lengths do not match: Expected 15 got 10", err.Error())
	assert.Equal(t, err, lerr.LenMismatch(15, 10))

	min, max, err = lerr.NewLenMismatch(5, 5)
	assert.Equal(t, 5, min)
	assert.Equal(t, 5, max)
	assert.NoError(t, err)
}
