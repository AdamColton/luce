package json32_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/serial/wrap/json/json32"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name string
}

func (*person) TypeID32() uint32 {
	return 123
}

type strSlice []string

func (strSlice) TypeID32() uint32 {
	return 789
}

type signal struct{}

func (signal) TypeID32() uint32 {
	return 456
}

type handlerObj chan string

func (ho handlerObj) HandleSignal(s signal) {
	ho <- "signal"
}

func (ho handlerObj) HandleStrSlice(s strSlice) {
	ho <- strings.Join(s, "|")
}

func (ho handlerObj) HandleFoo(p *person) {
	ho <- p.Name
}

func (ho handlerObj) ErrHandler(err error) {
	ho <- "Error: " + err.Error()
}

func (ho handlerObj) HandleFooErr(p *person) error {
	return errors.New(p.Name)
}

func TestSendReceive(t *testing.T) {
	bCh := make(chan []byte)

	s := json32.NewSender(bCh)
	r := json32.NewReceiver(bCh)

	s.RegisterType(strSlice(nil))
	r.RegisterType(strSlice(nil))

	done := make(chan bool)
	go func() {
		r.Run()
		done <- true
	}()

	go s.Send(strSlice{"this", "is", "a", "test"})
	select {
	case <-time.After(time.Millisecond * 5):
		t.Error("Timeout: failed to send")
	case v := <-r.Out:
		assert.Equal(t, strSlice{"this", "is", "a", "test"}, v)
	}

	close(bCh)
	select {
	case <-time.After(time.Millisecond * 5):
		t.Error("Timeout: failed to close")
	case <-done:
	}
}

func TestListeners(t *testing.T) {
	strCh := make(chan string)

	tests := map[string]struct {
		handler  interface{}
		send     type32.TypeIDer32
		expected string
	}{
		"*person": {
			handler: func(p *person) {
				strCh <- p.Name
			},
			send:     &person{Name: "person test"},
			expected: "person test",
		},
		"signal": {
			handler: func(s signal) {
				strCh <- "signal test"
			},
			send:     signal{},
			expected: "signal test",
		},
		"strSlice": {
			handler: func(s strSlice) {
				strCh <- strings.Join(s, "|")
			},
			send:     strSlice{"a", "b", "c"},
			expected: "a|b|c",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bCh := make(chan []byte)
			s := json32.NewSender(bCh)
			s.RegisterType(tc.send)

			l, err := serialbus.NewListener(bCh, s.ReaderDeserializer(json.Deserialize), type32.NewTypeMap(), nil, tc.handler)
			assert.NoError(t, err)

			done := make(chan bool)
			go func() {
				l.Run()
				done <- true
			}()

			go s.Send(tc.send)
			assert.Equal(t, tc.expected, <-strCh)

			close(bCh)
			select {
			case <-time.After(time.Millisecond * 5):
				t.Error("Timeout: failed to close")
			case <-done:
			}
		})
	}
}

func TestRegisterHandlers(t *testing.T) {
	ho := make(handlerObj)

	bCh := make(chan []byte)
	s := json32.NewSender(bCh)

	err := serial.RegisterTypes(s,
		signal{},
		strSlice(nil),
		(*person)(nil),
	)
	assert.NoError(t, err)

	h, err := json32.NewHandler(bCh, ho)
	assert.NoError(t, err)

	done := make(chan bool)
	go func() {
		h.Run()
		done <- true
	}()

	s.Send(signal{})
	assert.Equal(t, "signal", <-ho)

	s.Send(strSlice{"3", "1", "4"})
	assert.Equal(t, "3|1|4", <-ho)

	s.Send(&person{Name: "RegisterHandlers"})
	got := make(map[string]bool)
	timeout.After(100, func() {
		got[<-ho] = true
		got[<-ho] = true
	})
	// Order of these messages is not guarenteed
	assert.True(t, got["Error: RegisterHandlers"])
	assert.True(t, got["RegisterHandlers"])
	close(bCh)
	select {
	case <-time.After(time.Millisecond * 5):
		t.Error("Timeout: failed to close")
	case <-done:
	}
}
