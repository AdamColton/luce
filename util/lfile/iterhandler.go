package lfile

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/navigator"
)

// IterHandler represents something that will handle each value in the iterator.
type IterHandler interface {
	HandleIter(Iterator)
}

// IterHandlerFn is a function that fulfills IterHandler.
type IterHandlerFn func(Iterator)

// HandleIter fulfills IterHandler by calling fn.
func (fn IterHandlerFn) HandleIter(iter Iterator) {
	fn(iter)
}

// RunHandlerSource gets an Iterator from ii and calls HandleIter on ih for each
// value in it. It returns the error that ended the iteration, if there was one.
func RunHandlerSource(ii IteratorSource, ih IterHandler) error {
	i, done := ii.Iterator()
	for ; !done; _, done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}

// RunHandler resets i and calls HandleIter on ih for each value in it. It
// returns the error that ended the iteration, if there was one.
func RunHandler(i Iterator, ih IterHandler) error {
	for done := i.Reset(); !done; _, done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}

// GetByTypeHandler records all the files and directories the Iterator visits
// and seperates them by type.
type GetByTypeHandler struct {
	Files, Dirs []string
}

// HandleIter fulfills IterHandler and records the current location based on
// the type.
func (bt *GetByTypeHandler) HandleIter(i Iterator) {
	if i.Stat().IsDir() {
		bt.Dirs = append(bt.Dirs, i.Path())
	} else {
		bt.Files = append(bt.Files, i.Path())
	}
}

// GetFiles returns a handler that collects the paths of the files an Iterator
// visits. If cutPrefix is not empty it is removed from the start of each path.
// It appends to buf, which can be nil.
func GetFiles(cutPrefix string, buf []string) *GetType {
	return &GetType{
		GetDirs:   false,
		CutPrefix: cutPrefix,
		Matches:   buf[:0],
	}
}

// GetDirs returns a handler that collects the paths of the directories an
// Iterator visits. If cutPrefix is not empty it is removed from the start of
// each path. It appends to buf, which can be nil.
func GetDirs(cutPrefix string, buf []string) *GetType {
	return &GetType{
		GetDirs:   true,
		CutPrefix: cutPrefix,
		Matches:   buf[:0],
	}
}

// GetType collects the paths of either the files or the directories that an
// Iterator visits. Use GetFiles or GetDirs to create one.
type GetType struct {
	// GetDirs is true to collect directories and false to collect files.
	GetDirs bool
	// Matches are the paths that have been collected.
	Matches slice.Slice[string]
	// CutPrefix is removed from the start of each path if it is there.
	CutPrefix string
}

// HandleIter fulfills IterHandler. If the current value of the Iterator is of
// the type that is being collected its path is added to Matches.
func (gt *GetType) HandleIter(i Iterator) {
	if i.Stat().IsDir() == gt.GetDirs {
		p := i.Path()
		if gt.CutPrefix != "" {
			p, _ = strings.CutPrefix(p, gt.CutPrefix)
		}
		gt.Matches = append(gt.Matches, p)
	}
}

// GetContentsHandler reads the contents of all files into a map.
type GetContentsHandler map[string][]byte

// HandleIter fulfills IterHandler. If the current value of the Iterator is a
// file, it's contents are entered into the GetContentsHandler map.
func (c GetContentsHandler) HandleIter(i Iterator) {
	if !i.Stat().IsDir() {
		c[i.Path()] = i.Data()
	}
}

// MultiHandler is a slice of IterHandler. HandleIter will call HandleIter on
// each IterHandler in the slice
type MultiHandler []IterHandler

// HandleIter will call HandleIter on each IterHandler in the slice
func (mh MultiHandler) HandleIter(i Iterator) {
	for _, h := range mh {
		h.HandleIter(i)
	}
}

// FileTreeNode is a file or directory in a FilesTree. It is a
// navigator.Nexter, so a path can be followed through it with a
// navigator.Navigator.
type FileTreeNode interface {
	// Children of a directory by name. It is nil for a file.
	Children() lmap.Map[string, FileTreeNode]
	// Next finds the child called key. If there is none and create is true, one
	// is added. It is a directory if it is not the last of the keys that are being
	// followed, or if ctx says the last one is. A nil ctx creates files.
	Next(key string, create bool, ctx *FileTreeCtx) (FileTreeNode, bool)
	// Write writes the names of the children, and of theirs in turn, to w, each
	// on its own line and in order. Each line starts with cur, and cur is
	// followed by pad for each level down.
	Write(w io.Writer, cur, pad string)
	// IsDir is true for a directory, even one with no children.
	IsDir() bool
	// Name of the file or directory. It is empty for the root.
	Name() string
}

type fileTreeNode struct {
	children map[string]FileTreeNode
	name     string
}

// FilesTree is an IterHandler that builds a tree of the paths an Iterator
// visits. The path is split at every slash, and CutPrefix is removed from the
// start of it first if it is there.
type FilesTree struct {
	CutPrefix string
	tree      *fileTreeNode
}

// NewFilesTree creates an empty FilesTree that removes cutPrefix from the paths
// it is given.
func NewFilesTree(cutPrefix string) *FilesTree {
	return &FilesTree{
		CutPrefix: cutPrefix,
		tree: &fileTreeNode{
			children: map[string]FileTreeNode{},
		},
	}
}

// FileTreeCtx is the context that Next needs to create the nodes of a path. It
// is made by FilesTree, which knows how far along the path it is and whether it
// ends in a directory.
type FileTreeCtx struct {
	depth int
	isDir bool
}

// makesDir is true if a node created now is a directory. Without a context there
// is no way to know, so it is a file.
func (ctx *FileTreeCtx) makesDir() bool {
	return ctx != nil && (ctx.depth != 0 || ctx.isDir)
}

// Root of the tree. It is a directory with no name.
func (ft *FilesTree) Root() FileTreeNode {
	return ft.tree
}

func (ftn *fileTreeNode) Next(key string, create bool, ctx *FileTreeCtx) (FileTreeNode, bool) {
	if ctx != nil {
		ctx.depth--
	}
	child, found := ftn.children[key]
	if !found && create {
		childFtn := &fileTreeNode{
			name: key,
		}
		if ctx.makesDir() {
			childFtn.children = make(map[string]FileTreeNode)
		}
		ftn.children[key] = childFtn
		child = childFtn
		found = true
	}
	return child, found
}

func (ftn *fileTreeNode) Children() lmap.Map[string, FileTreeNode] {
	return ftn.children
}

func (ftn *fileTreeNode) IsDir() bool {
	return ftn.children != nil
}

func (ftn *fileTreeNode) Name() string {
	return ftn.name
}

func (ftn *fileTreeNode) Write(w io.Writer, cur, pad string) {
	keys := make([]string, 0, len(ftn.children))
	for k := range ftn.children {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintln(w, cur, k)
		ftn.children[k].Write(w, cur+pad, pad)
	}
}

var notEmptyStr = filter.NEQ("")

// HandleIter fulfills IterHandler. It adds the path of the current value to the
// tree, as a directory or a file to match the Iterator.
func (ft *FilesTree) HandleIter(i Iterator) {
	p := i.Path()
	if ft.CutPrefix != "" {
		p, _ = strings.CutPrefix(p, ft.CutPrefix)
	}
	keys := slice.New(strings.Split(p, "/"))
	keys = notEmptyStr.Slice(keys, keys)
	n := &navigator.Navigator[string, FileTreeNode, *FileTreeCtx]{
		Cur:  ft.tree,
		Idx:  0,
		Keys: keys,
	}
	ctx := &FileTreeCtx{
		depth: len(n.Keys),
		isDir: i.Stat().IsDir(),
	}
	n.Seek(true, ctx)
}
