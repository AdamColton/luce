package lfile

import (
	"fmt"
	"io"
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

// RunHandler will create an Iter from Iterator and call the HandleIter method
// on the IterHandler for each value in the iterator.
func RunHandlerSource(ii IteratorSource, ih IterHandler) error {
	i, done := ii.Iterator()
	for ; !done; _, done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}

// RunHandler will create an Iter from Iterator and call the HandleIter method
// on the IterHandler for each value in the iterator.
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

func GetFiles(cutPrefix string, buf []string) *GetType {
	return &GetType{
		GetDirs:   false,
		CutPrefix: cutPrefix,
		Matches:   buf[:0],
	}
}

func GetDirs(cutPrefix string, buf []string) *GetType {
	return &GetType{
		GetDirs:   true,
		CutPrefix: cutPrefix,
		Matches:   buf[:0],
	}
}

type GetType struct {
	GetDirs   bool
	Matches   slice.Slice[string]
	CutPrefix string
}

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

type FileTreeNode interface {
	Children() lmap.Map[string, FileTreeNode]
	Next(key string, create bool, ctx *FileTreeCtx) (FileTreeNode, bool)
	Write(w io.Writer, cur, pad string)
	IsDir() bool
	Name() string
}

type fileTreeNode struct {
	children map[string]FileTreeNode
	name     string
}

type FilesTree struct {
	CutPrefix string
	tree      *fileTreeNode
}

func NewFilesTree(cutPrefix string) *FilesTree {
	return &FilesTree{
		CutPrefix: cutPrefix,
		tree: &fileTreeNode{
			children: map[string]FileTreeNode{},
		},
	}
}

type FileTreeCtx struct {
	depth int
	isDir bool
}

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
		if ctx.depth != 0 || ctx.isDir {
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
	for k, c := range ftn.children {
		fmt.Fprintln(w, cur, k)
		c.Write(w, cur+pad, pad)
	}
}

var notEmptyStr = filter.NEQ("")

func (ft *FilesTree) HandleIter(i Iterator) {
	p := i.Path()
	if ft.CutPrefix != "" {
		p, _ = strings.CutPrefix(p, ft.CutPrefix)
	}
	keys := slice.New(strings.Split(p, "/"))
	keys = notEmptyStr.SliceBuf(keys, keys)
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
