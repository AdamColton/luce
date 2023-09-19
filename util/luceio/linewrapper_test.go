package luceio_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestBasicWrapping(t *testing.T) {
	text := []byte("aggrandize epistolography playwoman unreformable wretched supinate reassort relent kurchicine lithophyllous trilingual inventiveness historicoprophetic Bereshith musal unempty Lagothrix symbological zechin soundlessly arylate fetterbush probationism pluriseptate")

	buf := bytes.NewBuffer(make([]byte, 0, 300))
	w := luceio.NewLineWrappingWriter(buf)
	w.Write(text)
	expected := `aggrandize epistolography playwoman unreformable wretched supinate reassort
relent kurchicine lithophyllous trilingual inventiveness historicoprophetic
Bereshith musal unempty Lagothrix symbological zechin soundlessly arylate
fetterbush probationism pluriseptate`
	assert.Equal(t, expected, buf.String())

	buf.Reset()
	luceio.WrapWidth = 25
	w = luceio.NewLineWrappingWriter(buf)
	w.Write(text)
	luceio.WrapWidth = 80
	expected = `aggrandize epistolography
playwoman unreformable
wretched supinate
reassort relent
kurchicine lithophyllous
trilingual inventiveness
historicoprophetic
Bereshith musal unempty
Lagothrix symbological
zechin soundlessly
arylate fetterbush
probationism
pluriseptate`
	assert.Equal(t, expected, buf.String())

	buf.Reset()
	w = luceio.NewLineWrappingWriter(
		luceio.LineWrapperContextWriter{
			Writer: buf,
			Width:  35,
			Pad:    "// ",
		},
	)
	w.WritePadding()
	w.Write(text)
	expected = `// aggrandize epistolography
// playwoman unreformable wretched
// supinate reassort relent
// kurchicine lithophyllous
// trilingual inventiveness
// historicoprophetic Bereshith
// musal unempty Lagothrix
// symbological zechin soundlessly
// arylate fetterbush probationism
// pluriseptate`
	assert.Equal(t, expected, buf.String())
}

func TestContinueWrapping(t *testing.T) {
	buf := &bytes.Buffer{}
	w := luceio.NewLineWrappingWriter(
		luceio.LineWrapperContextWriter{
			Writer: buf,
			Width:  80,
		},
	)
	w.Write([]byte("aggrandize epistolography playwoman unreformable"))
	w.Write([]byte(" wretched supinate reassort relent kurchicine lithophyllous"))
	expected := `aggrandize epistolography playwoman unreformable wretched supinate reassort
relent kurchicine lithophyllous`
	assert.Equal(t, expected, buf.String())
}

func TestHandleNewLine(t *testing.T) {
	text := []byte("aggrandize epistolography playwoman unreformable wretched supinate reassort relent kurchicine lithophyllous\ntrilingual inventiveness historicoprophetic Bereshith musal unempty Lagothrix symbological zechin soundlessly arylate fetterbush probationism pluriseptate")

	buf := &bytes.Buffer{}
	w := luceio.NewLineWrappingWriter(buf)
	w.Write(text)
	expected := `aggrandize epistolography playwoman unreformable wretched supinate reassort
relent kurchicine lithophyllous
trilingual inventiveness historicoprophetic Bereshith musal unempty Lagothrix
symbological zechin soundlessly arylate fetterbush probationism pluriseptate`
	assert.Equal(t, expected, buf.String())
}

func TestWrapSumWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	sw := luceio.NewSumWriter(luceio.LineWrapperContextWriter{
		Writer: buf,
		Width:  10,
		Pad:    "> ",
	})
	w := luceio.NewLineWrappingWriter(sw)
	n, err := w.Write([]byte("aaa bbb ccc ddd"))
	assert.NoError(t, err)
	assert.Equal(t, "aaa bbb\n> ccc ddd", buf.String())

	// The SumWriter is shared, so it counts what the LineWrappingWriter writes.
	assert.Equal(t, int64(n), sw.Sum)
}

type failWriter struct {
	err error
}

func (fw failWriter) Write([]byte) (int, error) {
	return 0, fw.err
}

func TestWriteError(t *testing.T) {
	errTest := errors.New("test error")
	w := luceio.NewLineWrappingWriter(failWriter{errTest})

	_, err := w.Write([]byte("abc"))
	assert.Equal(t, errTest, err)

	// Once there is an error, nothing more is written.
	n, err := w.Write([]byte("def"))
	assert.Equal(t, errTest, err)
	assert.Zero(t, n)
}

func TestWriteNewline(t *testing.T) {
	buf := &bytes.Buffer{}
	w := luceio.NewLineWrappingWriter(
		luceio.LineWrapperContextWriter{
			Writer: buf,
			Width:  80,
			Pad:    "// ",
		},
	)
	n, err := w.WriteNewline()
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	w.Write([]byte("abc"))
	assert.Equal(t, "\n// abc", buf.String())
}
