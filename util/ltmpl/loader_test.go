package ltmpl

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func TestHTMLLoader(t *testing.T) {
	restore := lfile.ReadFile
	defer func() { lfile.ReadFile = restore }()
	lfile.ReadFile = func(filename string) ([]byte, error) {
		return []byte(filename + " - TEMPLATE"), nil
	}

	l := HTMLLoader{
		Trimmer:        lfile.PathLength(3),
		IteratorSource: lfile.Paths{"foo.bar", "bar.bar"},
	}
	tmpl, err := l.Load()
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	tmpl.ExecuteTemplate(buf, "foo.bar", nil)
	assert.Equal(t, "foo.bar - TEMPLATE", buf.String())

	buf.Reset()
	tmpl.ExecuteTemplate(buf, "bar.bar", nil)
	assert.Equal(t, "bar.bar - TEMPLATE", buf.String())

	expected := fmt.Errorf("Test Error")
	lfile.ReadFile = func(filename string) ([]byte, error) {
		return nil, expected
	}
	tmpl, err = l.Load()
	assert.Nil(t, tmpl)
	assert.Equal(t, expected, err)
}
