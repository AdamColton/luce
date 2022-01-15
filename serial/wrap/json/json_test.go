package json

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"

	"github.com/adamcolton/luce/serial/wrap/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSerializeDeserialize(t *testing.T) {
	testutil.SerialFuncsRoundTrip(t, Serialize, Deserialize)
}

type bufCloser struct {
	*bytes.Buffer
}

func (bufCloser) Close() error {
	return nil
}

func TestSerializerDeserializer(t *testing.T) {
	s := NewSerializer("", "  ")
	d := Deserializer{}
	testutil.SerialFuncsRoundTrip(t, s.WriteTo, d.ReadFrom)

	type Name struct {
		Name string
	}

	s.Prefix = "x"
	expected := []byte("{\nx  \"Name\": \"Adam\"\nx}\n")
	buf := make([]byte, 0, len(expected))
	got, err := s.Serialize(Name{
		Name: "Adam",
	}, buf)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	n := &Name{}
	d.Deserialize(n, []byte("{\n  \"Name\": \"Lauren\"\n}\n"))
	assert.Equal(t, "Lauren", n.Name)

	restoreOpen, restoreCreate := osOpen, osCreate
	defer func() {
		osOpen, osCreate = restoreOpen, restoreCreate
	}()

	bc := bufCloser{bytes.NewBuffer(nil)}
	osCreate = func(path string) (io.WriteCloser, error) {
		assert.Equal(t, "path-create", path)
		return bc, nil
	}
	osOpen = func(path string) (io.ReadCloser, error) {
		assert.Equal(t, "path-open", path)
		return bc, nil
	}

	s.Prefix = ""
	err = s.Save(n, "path-create")
	assert.NoError(t, err)
	assert.Equal(t, "{\n  \"Name\": \"Lauren\"\n}\n", bc.String())
	gotName := &Name{}
	err = d.Load(gotName, "path-open")
	assert.NoError(t, err)
	assert.Equal(t, "Lauren", gotName.Name)
}

func TestSaveLoad(t *testing.T) {
	type Name struct {
		Name string
	}
	path := filepath.Join(t.TempDir(), "name.json")

	s := NewSerializer("", "  ")
	assert.NoError(t, s.Save(Name{Name: "Adam"}, path))

	got := &Name{}
	assert.NoError(t, Deserializer{}.Load(got, path))
	assert.Equal(t, "Adam", got.Name)
}

func TestErrors(t *testing.T) {
	dir := t.TempDir()

	// The directory does not exist, so the file can't be created.
	err := NewSerializer("", "").Save(1, filepath.Join(dir, "missing", "x.json"))
	assert.Error(t, err)

	// The file does not exist.
	var i int
	err = Deserializer{}.Load(&i, filepath.Join(dir, "missing.json"))
	assert.Error(t, err)

	// A channel can't be encoded as JSON.
	_, err = Serializer{}.Serialize(make(chan int), nil)
	assert.Error(t, err)
}
