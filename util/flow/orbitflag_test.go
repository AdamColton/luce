package flow_test

import (
	"testing"

	"github.com/adamcolton/luce/util/flow"
	"github.com/stretchr/testify/assert"
)

func TestOrBitFlag(t *testing.T) {
	either := flow.NewOrFlag(uint8(0b001), uint8(0b110))
	assert.False(t, either.Check(0))
	assert.True(t, either.Check(0b001))
	assert.False(t, either.Check(0b010), "a flag with several bits needs all of them")
	assert.True(t, either.Check(0b110))
	assert.True(t, either.Check(0b111))

	assert.False(t, flow.NewOrFlag[uint8]().Check(0b111), "no flags")
	assert.True(t, flow.NewOrFlag[uint8](0).Check(0), "a flag of 0 is always set")
}
