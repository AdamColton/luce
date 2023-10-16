package cli

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"unsafe"

	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/parsers"
)

// Context is what a command line needs from a terminal, or whatever stands in
// for one: somewhere to write and a source of lines to read. It is an io.Writer
// and an io.Reader, so it can be passed on wherever those are wanted.
type Context interface {
	io.Writer
	io.Reader
	// WriteString writes a string.
	WriteString(str string) (int, error)
	// WriteStrings writes each string in turn and stops at the first error. It
	// returns the number of bytes written.
	WriteStrings(strs ...string) (int, error)
	// ReadString waits for a line and returns it with the white space around it
	// removed. If a value is sent on the cancel channel first it returns "". A nil
	// cancel channel is never sent on. It returns "" without waiting once the input
	// is closed.
	ReadString(cancel <-chan bool) string
	// Input writes the prompt, reads a line and parses it into v, which is a
	// pointer to any type that the Parser handles. If the line can't be parsed it
	// says so and asks again. It returns false if the input is cancelled with
	// ctrl+x or is closed.
	Input(prompt string, v any) bool
	// PopulateStruct calls Input for each field of the struct that s points to.
	// The prompt is made from cmd and the name of the field, and the field's
	// prompt tag if it has one. It returns false as soon as one is cancelled.
	PopulateStruct(cmd string, s any) bool
	// Parser is what Input and PopulateStruct use to parse text into values.
	Parser() reflector.Parser[string]
	// Closed is true once a read found that the input is closed. Nothing can be
	// read after that.
	Closed() bool
}

// NewContext creates a Context that writes to w and reads lines from in. If
// parser is nil, Parser is used.
func NewContext(w io.Writer, in <-chan []byte, parser reflector.Parser[string]) Context {
	if parser == nil {
		parser = Parser
	}
	return &context{
		Writer: w,
		in:     in,
		parser: parser,
	}
}

type context struct {
	io.Writer
	in     <-chan []byte
	parser reflector.Parser[string]
	buf    []byte
	closed bool
}

func (c *context) WriteString(str string) (int, error) {
	return io.WriteString(c.Writer, str)
}

func (c *context) WriteStrings(strs ...string) (int, error) {
	sum := 0
	for _, s := range strs {
		d, err := c.WriteString(s)
		sum += d
		if err != nil {
			return sum, err
		}
	}
	return sum, nil
}

// recv takes the next value off the input, and notes that the input is closed.
func (c *context) recv() []byte {
	bs, ok := <-c.in
	if !ok {
		c.closed = true
	}
	return bs
}

func (c *context) ReadString(cancel <-chan bool) string {
	var str string
	if cancel == nil {
		str = string(c.recv())
	} else {
		select {
		case bs, ok := <-c.in:
			c.closed = c.closed || !ok
			str = string(bs)
		case <-cancel:
			return ""
		}
	}
	return strings.TrimSpace(str)
}

func (c *context) Closed() bool {
	return c.closed
}

// Read reads what remains of the last value that was received, or waits for the
// next. It returns io.EOF once the input is closed.
func (c *context) Read(p []byte) (n int, err error) {
	n = len(c.buf)
	if n == 0 {
		c.buf = c.recv()
		if c.closed {
			return 0, io.EOF
		}
		n = len(c.buf)
	}
	if lnp := len(p); lnp < n {
		n = lnp
	}
	copy(p, c.buf[:n])
	c.buf = c.buf[n:]
	return
}

var (
	// allows for ctrl+x to cancel an operation
	cancel = string([]rune{24}) // 24 is ascii cancel

	// Parser is the default Parser for a Context. It parses strings, float64, int,
	// int64 and bool.
	Parser = reflector.Parser[string]{}
)

func init() {
	reflector.ParserAdd(Parser, parsers.String)
	reflector.ParserAdd(Parser, parsers.Float64)
	reflector.ParserAdd(Parser, parsers.Int)
	reflector.ParserAdd(Parser, parsers.Int64)
	reflector.ParserAdd(Parser, parsers.Bool)
}

func (c *context) Input(prompt string, v any) bool {
	for {
		c.WriteString(prompt)
		str := c.ReadString(nil)
		if str == cancel || c.closed {
			return false
		}
		err := c.parser.Parse(v, str)
		if err == nil {
			return true
		}
		c.WriteString("could not parse\n")
	}
}

// PopulateStruct panics if s is not a pointer to a struct. Fields that are not
// exported are populated too.
func (c *context) PopulateStruct(cmd string, s interface{}) bool {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		panic("Require pointer to struct")
	}
	v = v.Elem()
	ln := v.NumField()
	t := v.Type()
	for i := 0; i < ln; i++ {
		f := v.Field(i)
		sf := t.Field(i)
		prompt := sf.Tag.Get("prompt")
		p := reflect.NewAt(sf.Type, unsafe.Pointer(f.UnsafeAddr()))
		if prompt == "" {
			prompt = fmt.Sprintf("(%s:%s) ", cmd, sf.Name)
		} else {
			prompt = fmt.Sprintf("(%s:%s %s) ", cmd, sf.Name, prompt)
		}
		if !c.Input(prompt, p.Interface()) {
			return false
		}
	}
	return true
}

func (c *context) Parser() reflector.Parser[string] {
	return c.parser
}
