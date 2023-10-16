package cli_test

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/adamcolton/luce/util/cli"
	"github.com/stretchr/testify/assert"
)

// safeBuffer is a bytes.Buffer that one goroutine can write while another reads.
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *safeBuffer) WriteString(s string) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.WriteString(s)
}

func (b *safeBuffer) Read(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Read(p)
}

func (b *safeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func (b *safeBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf.Reset()
}

// newTestContext creates a Context that writes to the buffer and reads the
// lines, which are already waiting on a channel that is not closed.
func newTestContext(lines ...string) (cli.Context, *bytes.Buffer) {
	in := make(chan []byte, len(lines))
	for _, line := range lines {
		in <- []byte(line)
	}
	buf := bytes.NewBuffer(nil)
	return cli.NewContext(buf, in, nil), buf
}

// failWriter fails after it has written 'ok' times.
type failWriter struct {
	ok int
}

func (fw *failWriter) Write(data []byte) (int, error) {
	if fw.ok == 0 {
		return 0, errors.New("write failed")
	}
	fw.ok--
	return len(data), nil
}

func TestWriteStrings(t *testing.T) {
	ctx, buf := newTestContext()
	n, err := ctx.WriteStrings("ab", "cde", "")
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "abcde", buf.String())

	// It stops at the first error.
	ctx = cli.NewContext(&failWriter{ok: 1}, nil, nil)
	n, err = ctx.WriteStrings("ab", "cd", "ef")
	assert.EqualError(t, err, "write failed")
	assert.Equal(t, 2, n)
}

func TestInput(t *testing.T) {
	ctx, buf := newTestContext("abc", "  42  ")
	var age int
	assert.True(t, ctx.Input("age? ", &age))
	assert.Equal(t, 42, age)

	// It asks again after a line that can't be parsed.
	assert.Equal(t, "age? could not parse\nage? ", buf.String())

	// ctrl+x cancels.
	ctx, _ = newTestContext("\x18")
	age = 0
	assert.False(t, ctx.Input("age? ", &age))
	assert.Equal(t, 0, age)
}

func TestPopulateStructPrompts(t *testing.T) {
	type Req struct {
		Name   string `prompt:"who?"`
		Count  int
		hidden bool
	}

	ctx, buf := newTestContext("Ada", "3", "y")
	var r Req
	assert.True(t, ctx.PopulateStruct("greet", &r))

	// The prompt has the command and the field, and the prompt tag if there is
	// one. Fields that are not exported are asked for as well.
	assert.Equal(t, "(greet:Name who?) (greet:Count) (greet:hidden) ", buf.String())
	assert.Equal(t, Req{Name: "Ada", Count: 3, hidden: true}, r)

	// It stops when one is cancelled.
	ctx, _ = newTestContext("Grace", "\x18", "unused")
	r = Req{}
	assert.False(t, ctx.PopulateStruct("greet", &r))
	assert.Equal(t, Req{Name: "Grace"}, r)

	assert.PanicsWithValue(t, "Require pointer to struct", func() {
		ctx.PopulateStruct("greet", Req{})
	})
}

func TestClosedInput(t *testing.T) {
	in := make(chan []byte, 1)
	in <- []byte("last")
	close(in)
	ctx := cli.NewContext(&simpleWriter{}, in, nil)
	assert.False(t, ctx.Closed())

	// The last value is still there to be read, and then the input is closed.
	buf := make([]byte, 8)
	n, err := ctx.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, "last", string(buf[:n]))
	assert.False(t, ctx.Closed())

	n, err = ctx.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)
	assert.True(t, ctx.Closed())
}

func TestClosedReadString(t *testing.T) {
	for name, cancel := range map[string]<-chan bool{"no cancel": nil, "cancel": make(chan bool)} {
		t.Run(name, func(t *testing.T) {
			in := make(chan []byte)
			close(in)
			ctx := cli.NewContext(&simpleWriter{}, in, nil)
			assert.Equal(t, "", ctx.ReadString(cancel))
			assert.True(t, ctx.Closed())

			// Input gives up rather than asking for ever.
			var n int
			assert.False(t, ctx.Input("? ", &n))
		})
	}
}
