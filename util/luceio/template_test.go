package luceio_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	data := struct {
		Test string
	}{
		Test: "testing",
	}
	tmpl := template.Must(template.New("test").Parse(`base template{{define "core"}}My name is {{.Test}}{{end}}`))
	twt := luceio.NewTemplateTo(tmpl, "", data)
	buf := &bytes.Buffer{}
	twt.WriteTo(buf)
	assert.Equal(t, "base template", buf.String())
	buf.Reset()
	twt.Name = "core"
	n, err := twt.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, "My name is testing", buf.String())
	assert.Equal(t, int64(buf.Len()), n)
}

func TestTemplateWrapper(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(`base{{define "core"}}core {{.}}{{end}}`))
	tw := luceio.TemplateWrapper{TemplateExecutor: tmpl}

	buf := &bytes.Buffer{}
	n, err := tw.TemplateTo("core", "text").WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, "core text", buf.String())
	assert.Equal(t, int64(9), n)

	buf.Reset()
	_, err = tw.TemplateTo("", nil).WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, "base", buf.String())
}

func TestTemplateToError(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(`{{.Missing}}`))
	n, err := luceio.NewTemplateTo(tmpl, "", struct{}{}).WriteTo(&bytes.Buffer{})
	assert.Error(t, err)
	assert.Zero(t, n)
}
