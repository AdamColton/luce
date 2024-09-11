package lexec

import (
	"bytes"
	"fmt"
)

type Mock struct {
	Stdout, Stderr, Stdin *bytes.Buffer
	Err                   error
}

func NewMock() *Mock {
	out := bytes.NewBuffer(nil)
	return &Mock{
		Stdout: out,
		Stderr: out,
		Stdin:  bytes.NewBuffer(nil),
	}
}

func (m *Mock) ShellCmd(format string, args ...any) *bytes.Buffer {
	fmt.Fprintf(m.Stdin, format, args...)
	return m.Stdout
}

func (m *Mock) GetInput() string {
	str := m.Stdin.String()
	m.Stdin.Reset()
	return str
}

func (m *Mock) Run() error {
	return m.Err
}

func (m *Mock) New() Commander {
	return m
}
