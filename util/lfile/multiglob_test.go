package lfile

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiglob(t *testing.T) {
	restore := Glob
	defer func() { Glob = restore }()

	mockdir := map[string][]string{
		"foo*":  {"foo", "fooer", "foo.bar"},
		"*.bar": {"foo.bar", "bar.bar"},
	}
	Glob = func(pattern string) ([]string, error) {
		matches, found := mockdir[pattern]
		if found {
			return matches, nil
		}
		return nil, filepath.ErrBadPattern
	}

	got, err := MultiGlob{"foo*", "*.bar"}.Filenames()
	assert.NoError(t, err)
	sort.Strings(got)
	expected := []string{"bar.bar", "foo", "foo.bar", "fooer"}
	assert.Equal(t, expected, got)

	_, err = MultiGlob{"]]] == bad pattern == [[["}.Filenames()
	assert.Equal(t, filepath.ErrBadPattern, err)
}

func TestMultiGlobIter(t *testing.T) {
	restoreGlob := Glob
	restoreReadFile := ReadFile
	defer func() { Glob, ReadFile = restoreGlob, restoreReadFile }()

	ReadFile = func(filename string) ([]byte, error) {
		return []byte(filename), nil
	}

	mockdir := map[string][]string{
		"foo*":  {"foo", "fooer", "foo.bar"},
		"*.bar": {"foo.bar", "bar.bar"},
	}
	Glob = func(pattern string) ([]string, error) {
		matches, found := mockdir[pattern]
		if found {
			return matches, nil
		}
		return nil, filepath.ErrBadPattern
	}

	c := 0
	for i, done := (MultiGlob{"foo*", "*.bar"}).Iter(true); !done; done = i.Next() {
		c++
		assert.Equal(t, i.Filename, string(i.Data))
	}
	assert.Equal(t, 4, c)

	i, done := MultiGlob{"]]] == bad pattern == [[["}.Iter(true)
	assert.True(t, done)
	assert.Equal(t, filepath.ErrBadPattern, i.Err)

}
