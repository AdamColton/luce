package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/lexec"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

type BashCmd struct {
	Format  string
	Auto    bool
	Args    slice.Slice[any]
	running bool
	buf     *bytes.Buffer
}

func (bc *BashCmd) Run(cmdr lexec.Commander) {
	if bc.running {
		return
	}
	bc.running = true
	bc.buf = cmdr.ShellCmd(bc.Format, bc.Args...)
	go func() {
		cmdr.Run()
		bc.running = false
	}()
}

func (bc *BashCmd) Running() bool {
	return bc.running
}

func (s *Server) listBashCmds(w http.ResponseWriter, r *http.Request, d *struct {
	Session  *lusess.Session
	Redirect string
}) {
	if !d.Session.User().In("admin") {
		d.Redirect = "/"
		return
	}

	fmt.Fprintln(w, "--bash commands--", "<br>")

	for _, c := range s.BashCommands {
		fmt.Fprintln(w, c.Format, "<br>")
	}
}
