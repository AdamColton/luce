package lfile_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

// config is shaped like the configuration of a service that serves templates.
type config struct {
	Addr      string
	Templates struct {
		Re string
		lfile.PathLength
	}
}

const configJSON = `{"Addr":":8080","Templates":{"Re":"\\.html$","PathLength":2}}`

func TestJsonConfigFS(t *testing.T) {
	files := fstest.MapFS{
		"config.json": {Data: []byte(configJSON)},
		"other.json":  {Data: []byte(`{"Addr":":9090"}`)},
	}
	want := func(addr string) (c config) {
		c.Addr = addr
		return
	}

	// with no environment variable the default is the file
	var c config
	assert.NoError(t, lfile.JsonConfigFS(files, "", "config.json", &c))
	assert.Equal(t, ":8080", c.Addr)
	assert.Equal(t, `\.html$`, c.Templates.Re)
	assert.Equal(t, lfile.PathLength(2), c.Templates.PathLength)

	// a variable that is set names the file
	t.Setenv("LFILE_TEST_CONFIG", "other.json")
	c = config{}
	assert.NoError(t, lfile.JsonConfigFS(files, "LFILE_TEST_CONFIG", "config.json", &c))
	assert.Equal(t, want(":9090"), c)

	// one that is empty or not set falls back to the default
	t.Setenv("LFILE_TEST_CONFIG", "")
	c = config{}
	assert.NoError(t, lfile.JsonConfigFS(files, "LFILE_TEST_CONFIG", "config.json", &c))
	assert.Equal(t, ":8080", c.Addr)
	c = config{}
	assert.NoError(t, lfile.JsonConfigFS(files, "LFILE_TEST_UNSET", "config.json", &c))
	assert.Equal(t, ":8080", c.Addr)
}

func TestJsonConfigFSErrors(t *testing.T) {
	core := newCoreOnly(fstest.MapFS{
		"bad.json": {Data: []byte(`{"Addr":`)},
		"ok.json":  {Data: []byte(configJSON)},
	})

	var c config
	err := lfile.JsonConfigFS(core, "", "missing.json", &c)
	assert.ErrorIs(t, err, fs.ErrNotExist)

	// the file is closed, however it ends
	assert.Error(t, lfile.JsonConfigFS(core, "", "bad.json", &c))
	assert.NoError(t, lfile.JsonConfigFS(core, "", "ok.json", &c))
	assert.Equal(t, 2, *core.opened)
	assert.Equal(t, 2, *core.closed)
}

func TestJsonConfig(t *testing.T) {
	dir := t.TempDir()
	def := filepath.Join(dir, "default.json")
	env := filepath.Join(dir, "env.json")
	assert.NoError(t, os.WriteFile(def, []byte(configJSON), 0o644))
	assert.NoError(t, os.WriteFile(env, []byte(`{"Addr":":9090"}`), 0o644))

	var c config
	assert.NoError(t, lfile.JsonConfig("", def, &c))
	assert.Equal(t, ":8080", c.Addr)

	t.Setenv("LFILE_TEST_CONFIG", env)
	c = config{}
	assert.NoError(t, lfile.JsonConfig("LFILE_TEST_CONFIG", def, &c))
	assert.Equal(t, ":9090", c.Addr)

	err := lfile.JsonConfig("", filepath.Join(dir, "missing.json"), &c)
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
