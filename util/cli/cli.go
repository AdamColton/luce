package cli

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"unsafe"

	"github.com/adamcolton/luce/util/reflector"
)

type Context interface {
	Write(b []byte) (int, error)
	WriteString(str string) (int, error)
	WriteStrings(strs ...string) (int, error)
	ReadString() string
	Input(prompt string, v any) bool
	PopulateStruct(cmd string, s any) bool
	Parser() reflector.Parser[string]
}

func NewContext(w io.Writer, in <-chan []byte, parser reflector.Parser[string]) Context {
	if parser == nil {
		parser = Parser
	}
	c := &context{
		Writer: w,
		in:     in,
		parser: parser,
	}

	// TODO: move to luceio?
	ws, ok := w.(io.StringWriter)
	if ok {
		c.writeString = ws.WriteString
	} else {
		c.writeString = func(str string) (int, error) {
			return w.Write([]byte(str))
		}
	}

	return c
}

type context struct {
	io.Writer
	in          <-chan []byte
	writeString func(str string) (int, error)
	parser      reflector.Parser[string]
}

func (c *context) WriteString(str string) (int, error) {
	return c.writeString(str)
}

func (c *context) WriteStrings(strs ...string) (int, error) {
	sum := 0
	for _, s := range strs {
		d, _ := c.WriteString(s)
		sum += d
	}
	return sum, nil
}

func (c *context) ReadString() string {
	str := string(<-c.in)
	return strings.TrimSpace(str)
}

var cancel = string([]rune{24})

var Parser = reflector.Parser[string]{
	reflector.Type[*string]():  reflector.Parsers.String,
	reflector.Type[*float64](): reflector.Parsers.Float64,
	reflector.Type[*int]():     reflector.Parsers.Int,
}

// Input provides a prompt and populates a value. Currently supports string or
// int.
func (c *context) Input(prompt string, v any) bool {
	for {
		c.WriteString(prompt)
		str := c.ReadString()
		if str == cancel {
			return false
		}
		err := c.parser.Parse(v, str)
		if err == nil {
			return true
		}
		c.WriteString("could not parse\n")
	}
}

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
