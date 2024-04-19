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

func TestUnpack(t *testing.T) {
	d := delim.Delimiter("\n")
	u := d.Unpacker()
	got := u.Unpack([]byte("this"))
	assert.Len(t, got, 0)
	got = u.Unpack([]byte(" is"))
	assert.Len(t, got, 0)
	got = u.Unpack([]byte(" a"))
	assert.Len(t, got, 0)
	got = u.Unpack([]byte(" test\nFOOOOO"))
	expected := [][]byte{
		[]byte("this is a test"),
	}
	assert.Equal(t, expected, got)

	got = u.Unpack([]byte(" 1\n2\n3\n456\n"))
	expected = [][]byte{
		[]byte("FOOOOO 1"),
		[]byte("2"),
		[]byte("3"),
		[]byte("456"),
	}
	assert.Equal(t, expected, got)
}
