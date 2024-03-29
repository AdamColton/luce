package luceio

import (
	"bytes"
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
