package core_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/server/core"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

type cmdr struct {
	*core.Server
	cli.Helper
	ec *cli.ExitClose
}

func (c *cmdr) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(c)
	handler.AddAlias(cmds,
		"help", "h",
	)
	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)

	return lerr.Must(handler.Cmds(cs))
}

func (c *cmdr) EC() *cli.ExitClose {
	return c.ec
}

func (c *cmdr) Handlers(rnr *cli.Runner) []any {
	return []any{
		rnr.ExitRespHandler,
		rnr.CloseRespHandler,
		rnr.HelpRespHandler,
	}
}

func TestCore(t *testing.T) {
	cfg := core.Config{
		Addr: ":53452",
	}
	srv := cfg.NewServer()
	expected := []byte("this is a test")
	srv.Router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	})

	closed := make(chan bool)
	go func() {
		srv.Run()
		closed <- true
	}()

	resp, err := http.Get("http://localhost" + cfg.Addr + "/test")
	assert.NoError(t, err)
	got, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	srv.CliHandler = func(ec *cli.ExitClose) cli.Commander {
		return &cmdr{
			Server: srv,
			ec:     ec,
		}
	}
	buf := bytes.NewBuffer(nil)
	in := make(chan []byte)
	ctx := cli.NewContext(buf, in, nil)
	didExit := make(chan bool, 1)
	onExit := func() {
		didExit <- true
	}
	go srv.Cli(ctx, onExit)

	err = timeout.After(100, func() {
		for {
			str := buf.String()
			if str != "" {
				assert.Equal(t, str, "> ")
				buf.Reset()
				return
			}
			time.Sleep(time.Millisecond)
		}
	})
	assert.NoError(t, err)
	in <- []byte("help")
	err = timeout.After(100, func() {
		for {
			str := buf.String()
			if str != "" {
				assert.Contains(t, str, "h, help")
				buf.Reset()
				break
			}
			time.Sleep(time.Millisecond)
		}
		err = srv.Close()
		assert.NoError(t, err)
		timeout.After(5, closed)
	})
	assert.NoError(t, err)

	timeout.After(10, didExit)
}
