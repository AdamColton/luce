package serialbus_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	bCh := make(chan []byte)
	sCh := serialbus.String(bCh)

	expect := []string{"Apple", "Banana", "Cantaloupe"}
	go func() {
		for _, s := range expect {
			bCh <- []byte(s)
		}
		close(bCh)
	}()

	for _, s := range expect {
		assert.Equal(t, s, <-sCh)
	}

	// The channel is closed after the one it reads.
	_, open := <-sCh
	assert.False(t, open)
}
