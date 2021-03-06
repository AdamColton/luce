package luceio

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestStringWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	var sw StringWriter = buf
	_, err := sw.WriteString("testing")
	assert.NoError(t, err)
	assert.Equal(t, "testing", buf.String())
}

func TestStringsWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	sw := NewSumWriter(buf)
	_, err := sw.WriteStrings("this", "is", "a", "test")
	assert.NoError(t, err)
	assert.Equal(t, "thisisatest", buf.String())
}

func TestSumWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	sw := NewSumWriter(buf)
	sw.WriteString("test")
	sw.WriteInt(1)
	sw.WriteRune('-')
	sw.Write([]byte("test2"))
	i, err := sw.WriterTo(StringWriterTo("-test3"))
	assert.NoError(t, err)
	assert.Equal(t, int64(6), i)
	assert.NoError(t, sw.Err)
	expectErr := fmt.Errorf("test error")
	sw.Err = expectErr
	sw.WriteString("test3")
	i, err = sw.Rets()
	assert.Equal(t, expectErr, err)
	assert.Equal(t, "test1-test2-test3", buf.String())
	assert.Equal(t, int64(17), i)

}

func TestMany(t *testing.T) {
	buf := &bytes.Buffer{}
	wm := WriterToMerge(StringWriterTo("test1"), StringWriterTo("test2"))
	wm = WriterToMerge(wm, StringWriterTo("test3"))

	tos := wm.(WriterTos)
	assert.Len(t, tos, 3)
	n, err := MultiWrite(buf, tos, ":")
	assert.NoError(t, err)
	assert.Equal(t, int64(17), n)
	assert.Equal(t, "test1:test2:test3", buf.String())
}

func TestTemplate(t *testing.T) {
	data := struct {
		Test string
	}{
		Test: "testing",
	}
	tmpl := template.Must(template.New("test").Parse(`base template{{define "core"}}My name is {{.Test}}{{end}}`))
	twt := NewTemplateTo(tmpl, "", data)
	buf := &bytes.Buffer{}
	twt.WriteTo(buf)
	assert.Equal(t, "base template", buf.String())
	buf.Reset()
	twt.Name = "core"
	n, err := twt.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, "My name is testing", buf.String())
	assert.Equal(t, int(n), len(buf.Bytes()))
}

func TestJoin(t *testing.T) {
	assert.Equal(t, "a, b, c", Join("a", "b", "c", ", "))
	assert.Equal(t, "a", Join("a", ", "))
	assert.Equal(t, "", Join(", "))
}
