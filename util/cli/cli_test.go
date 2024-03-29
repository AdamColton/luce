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
