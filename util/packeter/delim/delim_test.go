package delim_test

import (
	"testing"

	"github.com/adamcolton/luce/util/packeter/delim"
	"github.com/stretchr/testify/assert"
)

func TestPack(t *testing.T) {
	d := delim.Delimiter("\n")
	got := d.Pack([]byte("this\nis\na\ntest"))
	expected := [][]byte{
		[]byte("this\n"),
		[]byte("is\n"),
		[]byte("a\n"),
		[]byte("test\n"),
	}
	assert.Equal(t, expected, got)

	got = d.Pack([]byte("this\nis\na\ntest\n"))
	assert.Equal(t, expected, got)

	got = d.Pack([]byte("fooooo"))
	expected = [][]byte{
		[]byte("fooooo\n"),
	}
	assert.Equal(t, expected, got)

	got = d.Pack(nil)
	expected = [][]byte{
		[]byte("\n"),
	}
	assert.Equal(t, expected, got)
}
