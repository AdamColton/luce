package lfile

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiglob(t *testing.T) {
	restore := setupForTestMultiglob()
	defer restore()

	got, err := MultiGlob{"foo*", "*.bar"}.Paths()
	assert.NoError(t, err)
	sort.Strings(got)
	expected := []string{"bar.bar", "foo", "foo.bar", "fooer"}
	assert.Equal(t, expected, got)

	_, err = MultiGlob{"]]] == bad pattern == [[["}.Paths()
	assert.Equal(t, filepath.ErrBadPattern, err)
}

func setupForTestMultiglob() func() {
	restoreGlob := Glob

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

	return func() { Glob = restoreGlob }
}

func TestMultiGlobIter(t *testing.T) {
	restore := setupForTestMultiGlobIter()
	defer restore()

	c := 0
	for i, done := (MultiGlob{"foo*", "*.bar"}).Iterator(); !done; done = i.Next() {
		c++
		assert.Equal(t, i.Path(), string(i.Data()))
	}
	assert.Equal(t, 4, c)

	i, done := MultiGlob{"]]] == bad pattern == [[["}.Iterator()
	assert.True(t, done)
	assert.Equal(t, filepath.ErrBadPattern, i.Err())

}

func setupForTestMultiGlobIter() func() {
	restoreGlob := Glob
	restoreReadFile := ReadFile
	restoreStat := Stat

	ReadFile = mockReadFileAsName
	Stat = mockStat
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

	return func() { Glob, ReadFile, Stat = restoreGlob, restoreReadFile, restoreStat }
}
