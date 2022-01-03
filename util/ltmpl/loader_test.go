package ltmpl

import (
	"bytes"
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

// failFS can't read a file.
type failFS struct {
	fstest.MapFS
	err error
}

func (f failFS) ReadFile(string) ([]byte, error) { return nil, f.err }

func TestHTMLLoader(t *testing.T) {
	files := fstest.MapFS{
		"foo.bar": {Data: []byte("foo.bar - TEMPLATE")},
		"bar.bar": {Data: []byte("bar.bar - TEMPLATE")},
	}
	paths := lfile.Paths{"foo.bar", "bar.bar"}

	l := HTMLLoader{
		Trimmer:        lfile.PathLength(3),
		IteratorSource: paths.FS(files),
	}
	tmpl, err := l.Load(nil)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	tmpl.ExecuteTemplate(buf, "foo.bar", nil)
	assert.Equal(t, "foo.bar - TEMPLATE", buf.String())

	buf.Reset()
	tmpl.ExecuteTemplate(buf, "bar.bar", nil)
	assert.Equal(t, "bar.bar - TEMPLATE", buf.String())

	expected := fmt.Errorf("Test Error")
	l.IteratorSource = paths.FS(failFS{files, expected})
	tmpl, err = l.Load(nil)
	assert.Nil(t, tmpl)
	assert.Equal(t, expected, err)
}
