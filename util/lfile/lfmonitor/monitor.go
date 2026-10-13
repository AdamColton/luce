// Package lfmonitor reports the files that are new, have changed or have been
// deleted since the last time they were looked at.
package lfmonitor

import (
	"io/fs"
	"os"
	"slices"
	"time"

	"github.com/adamcolton/luce/util/lfile"
)

// Monitor walks the files its Source finds and reports only the ones that are
// new, or whose modification time has changed, since the Monitor last reported
// them, followed by the ones it reported before that the Source no longer
// finds. A lfile.MatchRoot is a typical Source, and it decides which file system
// is walked. The zero value is ready to use once Source is set. A Monitor is not
// safe for concurrent use, and its walks must not overlap: finish one, or stop
// using it, before starting the next.
type Monitor struct {
	Source lfile.IteratorSource
	seen   map[string]record
	walk   uint
}

type record struct {
	mod  time.Time
	walk uint
}

// Iterator fulfills lfile.IteratorSource. It walks the Source and stops only at
// paths that are new or changed. Each path is recorded when the Iterator
// reaches it, so a path that is never reached, because the walk was stopped
// early, is reported again by the next walk. Once the Source is done without an
// error, the Iterator goes on to the paths it has reported before that the walk
// did not find, in order of name. For those lfile.Deleted is true, and Stat and
// Data are nil, and they are forgotten, so they are reported only once. Reset
// walks the Source again and reports what has changed since. Otherwise Data,
// Stat and Err are those of the Source's Iterator.
func (m *Monitor) Iterator() (i lfile.Iterator, done bool) {
	if m.seen == nil {
		m.seen = make(map[string]record)
	}
	src, _ := m.Source.Iterator()
	mi := &monitorIter{
		Iterator: src,
		m:        m,
	}
	m.walk++
	mi.seek()
	return mi, mi.Done()
}

// changed reports whether path is new or modified, and records its time and
// that this walk saw it. It compares with Equal rather than Before so a file
// replaced by an older copy is still reported.
func (m *Monitor) changed(path string, info fs.FileInfo) bool {
	prev, ok := m.seen[path]
	m.seen[path] = record{
		mod:  info.ModTime(),
		walk: m.walk,
	}
	return !ok || !prev.mod.Equal(info.ModTime())
}

// forgetUnseen removes the paths this walk did not see and appends them to buf.
func (m *Monitor) forgetUnseen(buf []string) []string {
	for path, r := range m.seen {
		if r.walk != m.walk {
			buf = append(buf, path)
			delete(m.seen, path)
		}
	}
	slices.Sort(buf)
	return buf
}

// monitorIter has two phases. While walking, it sits the Source's Iterator on
// a changed path, and Cur, Done, Path, Data, Stat and Err pass through to it.
// Once the walk is done, it steps through deleted, and the methods that
// describe the current path are answered from there.
type monitorIter struct {
	lfile.Iterator
	m        *Monitor
	idx      int
	deleting bool
	deleted  []string
	delIdx   int
}

func (mi *monitorIter) Wrapped() any {
	return mi.Iterator
}

// seek moves the Source's Iterator forward from its current path to the first
// one that changed. If there are none, and the walk didn't fail, it moves on to
// the deleted paths.
func (mi *monitorIter) seek() {
	for path, done := mi.Iterator.Cur(); !done; path, done = mi.Iterator.Next() {
		if mi.m.changed(path, mi.Iterator.Stat()) {
			return
		}
	}
	if mi.Iterator.Err() == nil {
		mi.deleting = true
		mi.deleted = mi.m.forgetUnseen(mi.deleted[:0])
		mi.delIdx = 0
	}
}

func (mi *monitorIter) Next() (string, bool) {
	if mi.deleting {
		if !mi.Done() {
			mi.delIdx++
		}
	} else {
		mi.Iterator.Next()
		mi.seek()
	}
	if !mi.Done() {
		mi.idx++
	}
	return mi.Cur()
}

func (mi *monitorIter) Cur() (string, bool) {
	return mi.Path(), mi.Done()
}

func (mi *monitorIter) Done() bool {
	if mi.deleting {
		return mi.delIdx >= len(mi.deleted)
	}
	return mi.Iterator.Done()
}

func (mi *monitorIter) Deleted() bool {
	return mi.deleting && !mi.Done()
}

func (mi *monitorIter) Path() string {
	if !mi.deleting {
		return mi.Iterator.Path()
	}
	if mi.Done() {
		return ""
	}
	return mi.deleted[mi.delIdx]
}

func (mi *monitorIter) Data() []byte {
	if mi.deleting {
		return nil
	}
	return mi.Iterator.Data()
}

func (mi *monitorIter) Stat() os.FileInfo {
	if mi.deleting {
		return nil
	}
	return mi.Iterator.Stat()
}

func (mi *monitorIter) Idx() int {
	return mi.idx
}

func (mi *monitorIter) Reset() bool {
	mi.deleting = false
	mi.m.walk++
	mi.Iterator.Reset()
	mi.seek()
	mi.idx = 0
	return mi.Done()
}
