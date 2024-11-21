package bus_test

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/json"
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
	done := make(chan bool)
	bCh := make(chan []byte)
	iCh := make(chan interface{})

	tm := type32.NewTypeMap()

	s := &serialbus.Sender{
		TypeSerializer: tm.WriterSerializer(json.Serialize),
		Chan:           bCh,
	}
	r := &serialbus.Receiver{
		In:               bCh,
		Out:              iCh,
		TypeDeserializer: tm.ReaderDeserializer(json.Deserialize),
		TypeRegistrar:    tm,
		ErrHandler: func(err error) {
			assert.NoError(t, err)
		},
	}
	r.RegisterType(strSlice(nil))
	go func() {
		r.Run()
		done <- true
	}()

	go func() {
		err := s.Send(strSlice{"this", "is", "a", "test"})
		assert.NoError(t, err)
	}()
	assert.Equal(t, strSlice{"this", "is", "a", "test"}, <-iCh)

	close(bCh)
	select {
	case <-time.After(time.Millisecond * 5):
		t.Error("Timeout: failed to close")
	case <-done:
	}
}

func TestMultiSender(t *testing.T) {
	type ch struct {
		b    chan []byte
		i    chan interface{}
		done chan bool
		r    bus.Receiver
	}

	tm := type32.NewTypeMap()
	sender := serialbus.NewMultiSender(tm.WriterSerializer(json.Serialize))
	chs := make([]*ch, 5)
	for i := range chs {
		iCh := make(chan interface{})
		bCh := make(chan []byte)
		done := make(chan bool)
		r := &serialbus.Receiver{
			In:               bCh,
			Out:              iCh,
			TypeDeserializer: tm.ReaderDeserializer(json.Deserialize),
			TypeRegistrar:    tm,
		}
		chs[i] = &ch{
			b:    bCh,
			i:    iCh,
			done: done,
			r:    r,
		}
		r.RegisterType(strSlice(nil))
		go func() {
			r.Run()
			done <- true
		}()

		if !assert.NoError(t, sender.Add(strconv.Itoa(i), bCh)) {
			return
		}
	}

	msg := strSlice{"this", "is", "a", "test"}
	assert.NoError(t, sender.Send(msg, "0"))
	assert.Equal(t, msg, <-chs[0].i)

	msg = strSlice{"twas", "brillig"}
	sender.Send(msg, "3", "1", "4")
	assert.Equal(t, msg, <-chs[4].i)
	assert.Equal(t, msg, <-chs[1].i)
	assert.Equal(t, msg, <-chs[3].i)

	msg = strSlice{"calling", "all", "channels"}
	sender.Send(msg)
	for _, c := range chs {
		assert.Equal(t, msg, <-c.i)
	}

	for _, c := range chs {
		close(c.b)
		select {
		case <-time.After(time.Millisecond * 5):
			t.Error("Timeout: failed to close")
		case <-c.done:
		}
		close(c.i)
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
			tm := type32.NewTypeMap()
			done := make(chan bool)
			bCh := make(chan []byte)
			s := &serialbus.Sender{
				TypeSerializer: tm.WriterSerializer(json.Serialize),
				Chan:           bCh,
			}
			l, err := serialbus.NewListener(bCh, tm.ReaderDeserializer(json.Deserialize), tm, nil, tc.handler)
			assert.NoError(t, err)
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
	tm := type32.NewTypeMap()
	ho := make(handlerObj)

	bCh := make(chan []byte)
	done := make(chan bool)
	s := &serialbus.Sender{
		TypeSerializer: tm.WriterSerializer(json.Serialize),
		Chan:           bCh,
	}
	r := &serialbus.Receiver{
		In:               bCh,
		TypeDeserializer: tm.ReaderDeserializer(json.Deserialize),
		TypeRegistrar:    tm,
	}
	l, err := bus.NewListener(r, nil, nil)
	assert.NoError(t, err)
	bus.RegisterHandlerType(l, ho)

	go func() {
		l.Run()
		done <- true
	}()

	s.Send(signal{})
	assert.Equal(t, "signal", <-ho)

	s.Send(strSlice{"3", "1", "4"})
	assert.Equal(t, "3|1|4", <-ho)

	// This will trigger both HandleFoo and HandleFooErr. HandleFooErr will
	// return an error which will be sent to the ErrHandler. The order in which
	// HandleFoo and HandleFooErr are called is not determinate.
	s.Send(&person{Name: "RegisterHandlers"})
	got := make([]string, 2)
	for i := range got {
		select {
		case <-time.After(time.Millisecond * 5):
			t.Error("Timeout: failed to receive from handlerObj")
		case x := <-ho:
			got[i] = x
		}
	}
	if assert.Len(t, got, 2) {
		sort.Strings(got)
		assert.Equal(t, "Error: RegisterHandlers", got[0])
		assert.Equal(t, "RegisterHandlers", got[1])
	}

	close(bCh)
	select {
	case <-time.After(time.Millisecond * 5):
		t.Error("Timeout: failed to close")
	case <-done:
	}
}