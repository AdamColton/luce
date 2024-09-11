package lexec

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Commander interface {
	ShellCmd(format string, args ...any) *bytes.Buffer
	Run() error
	New() Commander
}

type Command struct {
	*exec.Cmd
}

func New() *Command {
	return &Command{}
}

func (c *Command) New() Commander {
	return New()
}

func (c *Command) Wrapped() any {
	return c.Cmd
}

var ShellPath = "/bin/sh"

func (c *Command) ShellCmd(format string, args ...any) *bytes.Buffer {
	cmdStr := fmt.Sprintf(format, args...)
	c.Cmd = exec.Command(ShellPath, "-c", cmdStr)
	buf := bytes.NewBuffer(nil)
	c.Cmd.Stdout = buf
	c.Cmd.Stderr = buf
	return buf
}
