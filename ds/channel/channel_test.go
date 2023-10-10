package channel_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	ch := make(chan string)

	expected := "test"
	go func() { ch <- expected }()
	got, err := channel.TimeoutMS(5, ch)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	got, err = channel.TimeoutMS(5, ch)
	assert.Equal(t, channel.ErrTimeout, err)
	assert.Equal(t, "", got)
}
