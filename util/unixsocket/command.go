package unixsocket

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"unsafe"
)

// Command that a socket can accept.
type Command struct {
	Name   string
	Action func(*Context)
	Usage  string
}

// Context is passed into a Command
type Context struct {
	conn        net.Conn
	rawStr      string
	Socket      *Socket
	shouldClose bool
	Args        []string
	in          <-chan string
}

// Write to socket
func (c *Context) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

// WriteString to socket
func (c *Context) WriteString(str string) (int, error) {
	return c.conn.Write([]byte(str))
}

// Printf wraps fmt.Printf
func (c *Context) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(c, format, a...)
}

// Error will print an error if there is one. Returns a bool indicating if there
// was an error.
func (c *Context) Error(err error) bool {
	if err == nil {
		return false

	}
	c.conn.Write([]byte(err.Error()))
	return true
}

// String is what was passed into the Command
func (c *Context) String() string {
	return c.rawStr
}

// Close will cause the client end of the socket to close
func (c *Context) Close() {
	c.shouldClose = true
}

// Read a line from the socket connection
func (c *Context) Read() string {
	return <-c.in
}

var cancel = string([]rune{24})

// Input provides a prompt and populates a value. Currently supports string or
// int.
func (c *Context) Input(prompt string, v interface{}) bool {
	c.WriteString(prompt)
	str := <-c.in
	if str == cancel {
		return false
	}

	switch v := v.(type) {
	case *string:
		*v = str
	case *int:
		*v = c.getInt(prompt, str)
	}
	return true
}

func (c *Context) getInt(prompt, str string) int {
	var (
		i   int
		err error
	)
	for {
		i, err = strconv.Atoi(str)
		if err == nil {
			break
		}
		c.WriteString("Expected int; could not parse\n")
		c.WriteString(prompt)
		str = <-c.in
	}
	return i
}

// PopulateStruct will provide a prompt for each field
func (c *Context) PopulateStruct(cmd string, s interface{}) bool {
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
