package linject_test

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/adamcolton/luce/util/linject"
	"github.com/stretchr/testify/assert"
)

func TestCallPrintsInjectError(t *testing.T) {
	fi := linject.Initilizers{failingInjector{err: errors.New("boom")}}
	fn := fi.Apply(func(s string, data *struct{}) string { return s }).Interface().(func(string) string)

	r, w, err := os.Pipe()
	assert.NoError(t, err)
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	got := fn("hello")
	os.Stdout = stdout
	w.Close()
	out, err := io.ReadAll(r)
	assert.NoError(t, err)

	assert.Equal(t, "hello", got, "the function is still called")
	assert.Contains(t, string(out), "boom")
}
