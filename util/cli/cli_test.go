package cli_test

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/cli"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	in := make(chan []byte)
	c := cli.NewContext(buf, in, nil)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		read := func() string {
			str := ""
			for ; str == ""; time.Sleep(time.Millisecond) {
				str = buf.String()
			}
			buf.Reset()
			return str
		}

		str := read()
		assert.Equal(t, "(Person:Name) ", str)
		in <- []byte("Adam")

		str = read()
		assert.Equal(t, "(Person:Age) ", str)
		in <- []byte("39")

		str = read()
		assert.Equal(t, "testing", str)
		wg.Done()
	}()

	type Person struct {
		Name string
		Age  int
	}
	p := &Person{}
	ok := c.PopulateStruct("Person", p)
	assert.True(t, ok)
	assert.Equal(t, "Adam", p.Name)
	assert.Equal(t, 39, p.Age)

	c.WriteStrings("test", "ing")
	assert.Equal(t, cli.Parser, c.Parser())
	wg.Wait()
}

type simpleWriter struct {
	data []byte
}

func (sw *simpleWriter) Write(data []byte) (n int, err error) {
	sw.data = append(sw.data, data...)
	return len(data), nil
}

func TestReadStringCancel(t *testing.T) {
	w := &simpleWriter{}
	in := make(chan []byte)
	c := cli.NewContext(w, in, nil)

	c.WriteString("foobar")
	assert.Equal(t, "foobar", string(w.data))

	cancel := make(chan bool)
	out := make(chan string)
	fn := func() {
		out <- c.ReadString(cancel)
	}
	go fn()
	in <- []byte("test")
	assert.Equal(t, "test", <-out)

	go fn()
	cancel <- true
	assert.Equal(t, "", <-out)
}

func TestRead(t *testing.T) {
	w := &simpleWriter{}
	in := make(chan []byte)
	c := cli.NewContext(w, in, nil)

	go func() { in <- []byte("testing") }()

	r := make([]byte, 4)
	ln, err := c.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, 4, ln)
	assert.Equal(t, "test", string(r))

	go func() { in <- []byte(" foobar") }()

	ln, err = c.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, 3, ln)
	assert.Equal(t, "ing", string(r[:ln]))

	ln, err = c.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, 4, ln)
	assert.Equal(t, " foo", string(r[:ln]))

	ln, err = c.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, 3, ln)
	assert.Equal(t, "bar", string(r[:ln]))
}
