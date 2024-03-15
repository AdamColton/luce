package cmprtest

import (
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/stretchr/testify/assert"
)

// TODO: I think I'm going to take this out

// Assert wraps an instance of assert.Assertions but will replace calls to
// Equal with cmprtest when the type is float64 or fulfills AssertEqual. This
// allows calls to be made without passing in testing.T each time.
type Assert struct {
	*assert.Assertions
	assert.TestingT
}

// New creates an instance of Assert.
func New(t assert.TestingT) *Assert {
	return &Assert{
		Assertions: assert.New(t),
		TestingT:   t,
	}
}

// Equal will call cmprtest.EqualInDelta if expected is a float64 or if it
// fulfills AssertEqualizer.
func (g *Assert) Equal(expected, actual interface{}, msg ...interface{}) bool {
	if _, isAssert := expected.(cmpr.AssertEqualizer); isAssert {
		return EqualInDelta(g.TestingT, expected, actual, Small, msg...)
	}
	if _, isFloat := expected.(float64); isFloat {
		return EqualInDelta(g.TestingT, expected, actual, Small, msg...)
	}
	return g.Assertions.Equal(expected, actual, msg...)
}
