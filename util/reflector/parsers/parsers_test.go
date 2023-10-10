package parsers_test

import (
	"testing"

	"github.com/adamcolton/luce/util/reflector/parsers"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	var s string
	expected := "TestString"
	err := parsers.String(&s, expected)
	assert.NoError(t, err)
	assert.Equal(t, expected, s)
}

func TestFloat64(t *testing.T) {
	var f float64
	err := parsers.Float64(&f, "6.283")
	assert.NoError(t, err)
	assert.Equal(t, 6.283, f)

	err = parsers.Float64(&f, "not a float64")
	assert.Error(t, err)
}

func TestInt(t *testing.T) {
	var i int
	err := parsers.Int(&i, "54321")
	assert.NoError(t, err)
	assert.Equal(t, 54321, i)

	err = parsers.Int(&i, "not an int")
	assert.Error(t, err)
}

func TestInt64(t *testing.T) {
	var i int64
	err := parsers.Int64(&i, "54321")
	assert.NoError(t, err)
	assert.Equal(t, int64(54321), i)

	err = parsers.Int64(&i, "not an int")
	assert.Error(t, err)
}

func TestBool(t *testing.T) {
	var b bool
	err := parsers.Bool(&b, "y")
	assert.NoError(t, err)
	assert.True(t, b)
	err = parsers.Bool(&b, "n")
	assert.NoError(t, err)
	assert.False(t, b)
}
