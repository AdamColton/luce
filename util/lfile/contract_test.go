package lfile_test

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

// These tests pin the ways lfile is used by code that calls it: the services
// that serve templates from an embed.FS, the unixsocket search for sockets, the
// corpus example that reads every file of a directory, code that takes an
// OSRepository as a small interface, and the media/audiobooks repository that
// is tested with lfilemock. Each runs against the real OS, an embed.FS, a
// testing/fstest.MapFS and lfilemock where the caller could use any of them. A
// change to lfile that breaks one of these breaks a caller.

//go:embed testdata
var testdataFS embed.FS

// templates is the content of testdata/templates.
var templates = map[string]string{
	"index.html":       "<h1>index</h1>\n",
	"layout/base.html": "{{define \"base\"}}base{{end}}\n",
	"js/app.js":        "console.log(\"app\")\n",
	"notes.txt":        "not a template\n",
}

func TestTestdataMatchesTemplates(t *testing.T) {
	got := map[string]string{}
	assert.NoError(t, fs.WalkDir(testdataFS, "testdata/templates", func(name string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			data, err := testdataFS.ReadFile(name)
			assert.NoError(t, err)
			got[name[len("testdata/templates/"):]] = string(data)
		}
		return err
	}))
	assert.Equal(t, templates, got)
}

// templateFileSystems are the same files in each kind of file system, and the
// root they are under.
func templateFileSystems(t *testing.T) map[string]struct {
	cfs  lfile.CoreFS
	root string
} {
	type fsRoot = struct {
		cfs  lfile.CoreFS
		root string
	}

	mapFS := fstest.MapFS{}
	tree := map[string]any{}
	dir := t.TempDir()
	for name, data := range templates {
		mapFS["templates/"+name] = &fstest.MapFile{Data: []byte(data)}

		parent := tree
		parts := splitSlash(name)
		for _, p := range parts[:len(parts)-1] {
			next, ok := parent[p].(map[string]any)
			if !ok {
				next = map[string]any{}
				parent[p] = next
			}
			parent = next
		}
		parent[parts[len(parts)-1]] = data

		full := filepath.Join(dir, "templates", filepath.FromSlash(name))
		assert.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
		assert.NoError(t, os.WriteFile(full, []byte(data), 0o644))
	}

	return map[string]fsRoot{
		"embed":    {testdataFS, "testdata/templates"},
		"MapFS":    {mapFS, "templates"},
		"mock":     {lfilemock.Parse(map[string]any{"templates": tree}).Repository(), "templates"},
		"os":       {lfile.OSRepository{}, filepath.ToSlash(dir) + "/templates"},
		"os.DirFS": {os.DirFS(dir).(lfile.CoreFS), "templates"},
	}
}

func splitSlash(name string) (parts []string) {
	start := 0
	for i := 0; i <= len(name); i++ {
		if i == len(name) || name[i] == '/' {
			parts = append(parts, name[start:i])
			start = i + 1
		}
	}
	return
}

// TestContractTemplates is how the services load their templates: a
// MatchRoot over a "templates" directory of an embed.FS, set with SetCoreFS or
// by assigning CoreFS, walked with Iterator or a handler, with the name of each
// template made by PathLength.
func TestContractTemplates(t *testing.T) {
	m := lerr.Must(lfile.RegexMatch(`\.(html|js)$`, "", ""))
	// the name of a template keeps the last two parts of its path
	want := map[string]string{
		"templates/index.html": templates["index.html"],
		"layout/base.html":     templates["layout/base.html"],
		"js/app.js":            templates["js/app.js"],
	}

	for name, fsys := range templateFileSystems(t) {
		t.Run(name, func(t *testing.T) {
			// both ways of setting the file system
			byField := m.Root(fsys.root)
			byField.CoreFS = fsys.cfs
			// and the root may have a trailing slash
			for _, mr := range []lfile.MatchRoot{
				byField,
				m.Root(fsys.root).SetCoreFS(fsys.cfs),
				m.Root(fsys.root + "/").SetCoreFS(fsys.cfs),
			} {
				got := map[string]string{}
				i, done := mr.Iterator()
				for ; !done; _, done = i.Next() {
					assert.False(t, i.Stat().IsDir())
					got[lfile.PathLength(2).Trim(i.Path())] = string(i.Data())
				}
				assert.NoError(t, i.Err())
				assert.Equal(t, want, got)
			}
		})
	}
}

// TestContractEmptyRoot is a file system that is the templates directory, so
// its root is "" (or "."), as when a service is given a directory to serve.
func TestContractEmptyRoot(t *testing.T) {
	m := lerr.Must(lfile.RegexMatch(`\.(html|js)$`, "", ""))
	dir := t.TempDir()
	for name, data := range templates {
		full := filepath.Join(dir, filepath.FromSlash(name))
		assert.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
		assert.NoError(t, os.WriteFile(full, []byte(data), 0o644))
	}
	sub, err := fs.Sub(testdataFS, "testdata/templates")
	assert.NoError(t, err)

	for name, cfs := range map[string]lfile.CoreFS{
		"os.DirFS": os.DirFS(dir).(lfile.CoreFS),
		"embed":    sub.(lfile.CoreFS),
	} {
		t.Run(name, func(t *testing.T) {
			for _, root := range [][]string{nil, {""}, {"."}, {"./"}} {
				files := lfile.GetFiles("", nil)
				assert.NoError(t, lfile.RunHandlerSource(m.Root(root...).SetCoreFS(cfs), files))
				sort.Strings(files.Matches)
				assert.Equal(t, []string{"index.html", "js/app.js", "layout/base.html"}, []string(files.Matches), "%q", root)
			}

			// RootGetFiles ends the root with a slash, which used to panic for ""
			got, err := lfile.RootGetFiles("", cfs, m)
			assert.NoError(t, err)
			sort.Strings(got)
			assert.Equal(t, []string{"index.html", "js/app.js", "layout/base.html"}, []string(got))
		})
	}
}

// TestContractHandler is the corpus example: a handler that reads the data and
// path of every file, with the extensions and whether to enter hidden
// directories chosen at run time.
func TestContractHandler(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"a.txt":         "aaa",
		"sub/b.md":      "bb",
		"sub/c.go":      "not matched",
		".hidden/d.txt": "dddd",
	}
	for name, data := range files {
		full := filepath.Join(dir, filepath.FromSlash(name))
		assert.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
		assert.NoError(t, os.WriteFile(full, []byte(data), 0o644))
	}
	root := filepath.ToSlash(dir)

	for _, hidden := range []bool{false, true} {
		exts := lfile.Exts(hidden, "txt", "md")
		mtch := lerr.Must(lfile.RegexMatch(exts, "", ""))
		mr := mtch.Root(root)

		total := 0
		var paths []string
		hdlr := lfile.IterHandlerFn(func(i lfile.Iterator) {
			total += len(i.Data())
			paths = append(paths, i.Path())
		})
		assert.NoError(t, lfile.RunHandlerSource(mr, hdlr))
		sort.Strings(paths)

		if hidden {
			assert.Equal(t, 9, total)
			assert.Equal(t, []string{root + "/.hidden/d.txt", root + "/a.txt", root + "/sub/b.md"}, paths)
		} else {
			assert.Equal(t, 5, total)
			assert.Equal(t, []string{root + "/a.txt", root + "/sub/b.md"}, paths)
		}
	}
}

// TestContractSockets is the unixsocket search: one Match with a Root that is
// changed and used again, collected with Factory, over a root of "" (here) and
// a root with a trailing slash.
func TestContractSockets(t *testing.T) {
	repo := lfilemock.Parse(map[string]any{
		"a.sock": "",
		"c.txt":  "",
		"nested": map[string]any{"deep.sock": ""},
		"tmp":    map[string]any{"x.sock": "", "y.txt": ""},
	}).Repository()
	m := lerr.Must(lfile.RegexMatch(`.+\.sock`, "", ".*"))

	mr := m.Root("")
	mr.CoreFS = repo
	local := slice.FromIterFactory(mr.Factory, nil)
	assert.Equal(t, slice.Slice[string]{"a.sock"}, local, "the directories are skipped")

	mr.Root = "/tmp/"
	tmp := slice.FromIterFactory(mr.Factory, nil)
	assert.Equal(t, slice.Slice[string]{"/tmp/x.sock"}, tmp)

	// and on the real OS
	dir := t.TempDir()
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "os.sock"), nil, 0o644))
	mr = m.Root("")
	mr.CoreFS = lfile.OSRepository{}
	mr.Root = filepath.ToSlash(dir) + "/"
	assert.Equal(t, slice.Slice[string]{filepath.ToSlash(dir) + "/os.sock"}, slice.FromIterFactory(mr.Factory, nil))
}

// TestContractOSRepository is how OSRepository is used as the smaller
// interfaces: the default of a package variable or a field.
func TestContractOSRepository(t *testing.T) {
	dir := filepath.ToSlash(t.TempDir())
	var (
		reader  lfile.FSFileReader = lfile.OSRepository{}
		remover lfile.FSRemover    = lfile.OSRepository{}
		core    lfile.CoreFS       = lfile.OSRepository{}
		repo    interface {
			lfile.CoreFS
			lfile.FSCreator
		} = lfile.OSRepository{}
	)
	mock := lfilemock.Parse(map[string]any{"dir": map[string]any{}}).Repository()

	// Create and write, then read through each of the small interfaces
	f, err := repo.Create(dir + "/a.txt")
	assert.NoError(t, err)
	_, err = f.(interface{ WriteString(string) (int, error) }).WriteString("hello")
	assert.NoError(t, err)
	assert.NoError(t, f.Close())

	data, err := reader.ReadFile(dir + "/a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(data))
	des, err := core.ReadDir(dir)
	assert.NoError(t, err)
	assert.Len(t, des, 1)

	// an FSRemover fails for a file that is not there, and that is an
	// fs.ErrNotExist whichever file system it is (unixsocket ignores it)
	assert.NoError(t, remover.Remove(dir+"/a.txt"))
	assert.ErrorIs(t, remover.Remove(dir+"/a.txt"), fs.ErrNotExist)
	assert.ErrorIs(t, lfile.FSRemover(mock.(lfile.FSRemover)).Remove("dir/nope"), fs.ErrNotExist)
}

// TestContractAudiobooks is the media/audiobooks repository: a mock tree with
// empty directories and paths like "testroot/Cat1/...", a directory opened for
// both of its readers, RootGetFiles for the pictures and sound files, Size, and
// Create for the file it writes.
func TestContractAudiobooks(t *testing.T) {
	repo := lfilemock.Parse(map[string]any{
		"testroot": map[string]any{
			"Cat1": map[string]any{
				"Author1": map[string]any{
					"book1": map[string]any{},
					"book2": map[string]any{
						"cover.jpg": "12345",
						"part1.mp3": "1234567890",
						"notes.txt": "abc",
						"disc2": map[string]any{
							"part2.mp3": "12345",
						},
					},
				},
				"Author2": map[string]any{},
			},
		},
		"listening": map[string]any{
			"book2": map[string]any{},
		},
	}).Repository()

	// a directory is a DirReader for the categories and a DirNameReader for
	// the listening list, each from its own Open, addressed with or without
	// a trailing slash
	dir, err := repo.Open("testroot/Cat1")
	assert.NoError(t, err)
	rdr, ok := upgrade.To[lfile.DirReader](dir)
	assert.True(t, ok)
	authors, err := rdr.ReadDir(0)
	assert.NoError(t, err)
	if assert.Len(t, authors, 2) {
		assert.Equal(t, "Author1", authors[0].Name())
		assert.True(t, authors[0].IsDir())
	}
	dir, err = repo.Open("listening/")
	assert.NoError(t, err)
	names, ok := upgrade.To[lfile.DirNameReader](dir)
	assert.True(t, ok)
	list, err := names.Readdirnames(-1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"book2"}, list)

	// the pictures and sound files of a book, found recursively, named
	// relative to the book
	m := lerr.Must(lfile.RegexMatch(lfile.Exts(false, "jpg", "jpeg", "mp3"), "", ""))
	files, err := lfile.RootGetFiles("testroot/Cat1/Author1/book2", repo, m)
	assert.NoError(t, err)
	sort.Strings(files)
	assert.Equal(t, []string{"cover.jpg", "disc2/part2.mp3", "part1.mp3"}, []string(files))
	files, err = lfile.RootGetFiles("testroot/Cat1/Author1/book1", repo, m)
	assert.NoError(t, err)
	assert.Empty(t, files)

	// the size of a book and of a category
	size, err := lfile.Size(repo, "testroot/Cat1/Author1/book2")
	assert.NoError(t, err)
	assert.Equal(t, int64(23), size)
	size, err = lfile.Size(repo, "testroot/Cat1")
	assert.NoError(t, err)
	assert.Equal(t, int64(23), size)

	// meta.json is written into a book directory that exists
	meta, err := repo.Create("testroot/Cat1/Author1/book1/meta.json")
	assert.NoError(t, err)
	_, err = meta.(interface{ Write([]byte) (int, error) }).Write([]byte(`{"ToRead":true}`))
	assert.NoError(t, err)
	assert.NoError(t, meta.Close())
	data, err := repo.ReadFile("testroot/Cat1/Author1/book1/meta.json")
	assert.NoError(t, err)
	assert.Equal(t, `{"ToRead":true}`, string(data))
	// and one that is not there is missing, not nil
	_, err = repo.Open("testroot/Cat1/Author1/book9/meta.json")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
