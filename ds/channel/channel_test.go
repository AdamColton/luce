package channel_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	ch := make(chan string)

	expected := "test"
	go func() { ch <- expected }()
	got, err := channel.Timeout(5, ch)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	got, err = channel.Timeout(5, ch)
	assert.Equal(t, channel.ErrTimeout, err)
	assert.Equal(t, "", got)
}

func TestSlice(t *testing.T) {
	s := slice.New([]int{3, 1, 4})
	ch := channel.Slice(s, nil)
	s.Iter().For(func(i int) {
		assert.Equal(t, i, <-ch)
	})

	s = slice.New([]int{1, 5, 9, 2, 6, 5, 3, 5})
	go channel.Slice(s, ch)
	s.Iter().For(func(i int) {
		assert.Equal(t, i, <-ch)
	})
}
