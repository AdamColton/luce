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
	twt.WriteTo(buf)
	assert.Equal(t, "My name is testing", buf.String())
}
