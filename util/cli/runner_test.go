package cli_test

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	buf := bytes.NewBuffer(nil)
	in := make(chan []byte)
	wg := sync.WaitGroup{}

	var r *cli.Runner
	r = &cli.Runner{
		Context:      cli.NewContext(buf, in, nil),
		StartMessage: "Start Message\n",
		Prompt:       ">",
		Commands: lerr.Must(handler.Cmds([]*handler.Command{
			{
				Name:   "test",
				Usage:  "this is a test",
				Action: func() {},
			}, {
				Name: "exit",
				Action: func() {
					r.Exit = true
				},
			}, {
				Name: "person",
				Action: func(p *Person) {
					assert.Equal(t, p.Name, "Adam")
					assert.Equal(t, p.Age, 39)
					wg.Done()
				},
			},
		})),
	}

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

		assert.Equal(t, "Start Message\n>", read())
		wg.Add(1)
		in <- []byte("person")
		assert.Equal(t, "(person:Name) ", read())
		in <- []byte("Adam")
		assert.Equal(t, "(person:Age) ", read())
		in <- []byte("39")

		in <- []byte("exit")
		wg.Done()
	}()

	err := timeout.After(25, func() {
		r.Run()
		wg.Wait()

		wg.Add(1)
		r.Static([]string{"person", "Name:Adam", "Age:39"})
		wg.Wait()
	})
	assert.NoError(t, err)

}
