package bus

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type person struct {
	Name string
}

func (*person) TypeID32() uint32 {
	return 123
}

func TestListenerMux(t *testing.T) {
	done := make(chan bool)
	strCh := make(chan string)
	// hdlr will first put person.Name on the string channel then return an error
	// which will be picked up by the ErrHandler and will also be placed on the
	// string channel.
	personHdlr := func(f *person) error {
		strCh <- f.Name
		return errors.New("test error")
	}
	errHdlr := func(err error) {
		strCh <- err.Error()
	}
	personChan := make(chan *person)

	ifcCh := make(chan interface{})
	ml, err := NewListenerMux(ifcCh, errHdlr, personHdlr, personChan)
	assert.NoError(t, err)
	go func() {
		ml.Run()
		done <- true
	}()

	ifcCh <- &person{Name: "this is a test"}

	assert.Equal(t, "this is a test", <-strCh)
	assert.Equal(t, "test error", <-strCh)
	assert.Equal(t, "this is a test", (<-personChan).Name)

	close(ifcCh)
	select {
	case <-time.After(time.Millisecond * 5):
		t.Error("Timeout: failed to close")
	case <-done:
	}
}
