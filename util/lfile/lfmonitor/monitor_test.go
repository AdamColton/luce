package lfmonitor_test

import (
	"fmt"
	"testing"
	"testing/fstest"
	"time"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfmonitor"
	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

var monday = time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)

func notes() fstest.MapFS {
	return fstest.MapFS{
		"notes/groceries.tasks": {Data: []byte("[ ] milk"), ModTime: monday},
		"notes/garden.tasks":    {Data: []byte("[ ] weed"), ModTime: monday},
		"notes/readme.md":       {Data: []byte("not a task"), ModTime: monday},
	}
}

// paths walks src and returns what it found, marking deleted paths.
func paths(src lfile.IteratorSource) []string {
	var out []string
	i, _ := src.Iterator()
	liter.For(i, func(path string) {
		if lfile.Deleted(i) {
			path = "deleted " + path
		}
		out = append(out, path)
	})
	return out
}

func ExampleMonitor() {
	fsys := notes()
	m := &lfmonitor.Monitor{
		Source: lfile.MatchExt(true, "tasks").Root("notes").SetCoreFS(fsys),
	}

	fmt.Println("first walk", paths(m))
	fmt.Println("nothing changed", paths(m))

	fsys["notes/garden.tasks"].ModTime = monday.Add(time.Hour)
	fsys["notes/work.tasks"] = &fstest.MapFile{Data: []byte("[ ] report"), ModTime: monday}
	fmt.Println("edited and added", paths(m))

	delete(fsys, "notes/groceries.tasks")
	fmt.Println("removed", paths(m))
	fmt.Println("nothing changed", paths(m))

	// Output:
	// first walk [notes/groceries.tasks notes/garden.tasks]
	// nothing changed []
	// edited and added [notes/work.tasks notes/garden.tasks]
	// removed [deleted notes/groceries.tasks]
	// nothing changed []
}

func TestIteratorPassesThrough(t *testing.T) {
	fsys := notes()
	m := &lfmonitor.Monitor{
		Source: lfile.MatchExt(true, "tasks").Root("notes").SetCoreFS(fsys),
	}
	i, done := m.Iterator()
	assert.False(t, done)
	assert.Equal(t, 0, i.Idx())
	assert.Equal(t, "notes/groceries.tasks", i.Path())
	assert.Equal(t, "[ ] milk", string(i.Data()))
	assert.Equal(t, monday, i.Stat().ModTime())

	path, done := i.Next()
	assert.False(t, done)
	assert.Equal(t, "notes/garden.tasks", path)
	assert.Equal(t, 1, i.Idx())

	_, done = i.Next()
	assert.True(t, done)
	assert.NoError(t, i.Err())
	assert.Equal(t, "", i.Path())

	assert.True(t, i.Reset(), "Reset reports only what changed since")

	fsys["notes/groceries.tasks"].ModTime = monday.Add(-time.Hour)
	assert.False(t, i.Reset(), "an older copy is a change too")
	assert.Equal(t, "notes/groceries.tasks", i.Path())
	assert.Equal(t, 0, i.Idx())
}

func TestStoppedEarly(t *testing.T) {
	m := &lfmonitor.Monitor{
		Source: lfile.MatchExt(true, "tasks").Root("notes").SetCoreFS(notes()),
	}
	i, _ := m.Iterator()
	assert.Equal(t, "notes/groceries.tasks", i.Path())

	assert.Equal(t, []string{"notes/garden.tasks"}, paths(m), "a path never reached is reported by the next walk")
}

func TestWrapped(t *testing.T) {
	src := lfile.MatchExt(true, "tasks").Root("notes").SetCoreFS(notes())
	m := &lfmonitor.Monitor{Source: src}
	i, _ := m.Iterator()
	w, ok := i.(interface{ Wrapped() any })
	assert.True(t, ok)
	_, ok = w.Wrapped().(lfile.Iterator)
	assert.True(t, ok)
}

func TestDeleted(t *testing.T) {
	fsys := notes()
	m := &lfmonitor.Monitor{
		Source: lfile.MatchExt(true, "tasks").Root("notes").SetCoreFS(fsys),
	}
	paths(m)

	delete(fsys, "notes/groceries.tasks")
	delete(fsys, "notes/garden.tasks")
	fsys["notes/work.tasks"] = &fstest.MapFile{Data: []byte("[ ] report"), ModTime: monday}
	i, done := m.Iterator()
	assert.False(t, done)
	assert.Equal(t, "notes/work.tasks", i.Path())
	assert.False(t, lfile.Deleted(i))

	path, done := i.Next()
	assert.False(t, done)
	assert.Equal(t, "notes/garden.tasks", path)
	assert.True(t, lfile.Deleted(i))
	assert.Nil(t, i.Stat())
	assert.Nil(t, i.Data())
	assert.Equal(t, 1, i.Idx())

	path, _ = i.Next()
	assert.Equal(t, "notes/groceries.tasks", path)
	assert.Equal(t, 2, i.Idx())

	_, done = i.Next()
	assert.True(t, done)
	assert.False(t, lfile.Deleted(i))
	assert.Equal(t, "", i.Path())
	_, done = i.Next()
	assert.True(t, done, "Next past the end stays done")
	assert.Equal(t, 2, i.Idx())
	assert.NoError(t, i.Err())
}

func TestFailedWalkDeletesNothing(t *testing.T) {
	fsys := notes()
	m := &lfmonitor.Monitor{
		Source: lfile.MatchExt(true, "tasks").Root("notes").SetCoreFS(fsys),
	}
	paths(m)

	// The root is gone, as if a drive were unmounted. That is an error, not a
	// reason to report every file as deleted.
	for name := range fsys {
		delete(fsys, name)
	}
	i, done := m.Iterator()
	assert.True(t, done)
	assert.Error(t, i.Err())
	assert.False(t, lfile.Deleted(i))

	restored := notes()
	for name, f := range restored {
		fsys[name] = f
	}
	assert.Empty(t, paths(m), "the files are still recorded, so nothing is new")
}
