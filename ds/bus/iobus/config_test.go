package iobus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigMakeErrCh(t *testing.T) {
	c := Config{}
	assert.Nil(t, c.makeErrCh())

	c.MakeErrCh = true
	assert.NotNil(t, c.makeErrCh())
}
