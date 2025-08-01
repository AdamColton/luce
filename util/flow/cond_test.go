package flow_test

import (
	"testing"

	"github.com/adamcolton/luce/util/flow"
	"github.com/stretchr/testify/assert"
)

func TestNilCheck(t *testing.T) {
	type foo struct {
		Name string
	}
	var f *foo
	c := func() *foo {
		return &foo{Name: "test"}
	}

	f = flow.NilCheck(f, c)
	assert.Equal(t, "test", f.Name)

	f.Name = "not nil"
	f = flow.NilCheck(f, c)
	assert.Equal(t, "not nil", f.Name)
}

func TestTern(t *testing.T) {
	assert.Equal(t, 3, flow.Tern(false, 2, 3))
	assert.Equal(t, 2, flow.Tern(true, 2, 3))
}
