package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/util/luceio"
)

const defaultCommentWidth = 80

// Comment string that automatically wraps the string
type Comment string

// CommentWidther is anything that can specify a comment width
type CommentWidther interface {
	CommentWidth() int
}

// NewComment takes a string and returns a Comment. The comment width comes from
// the files package context. If File is nil a Comment with no width is
// returned. It does not write the comment to the file.
func (f *File) NewComment(comment string) {
	f.AddGenerator(Comment(comment))
}

var nl = []byte("\n")

// PrefixWriteTo wraps the comment and writes it to the Writer
func (c Comment) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	width := defaultCommentWidth
	if cw, ok := w.(CommentWidther); ok {
		width = cw.CommentWidth()
	} else if cw, ok := pre.(CommentWidther); ok {
		width = cw.CommentWidth()
	}
	lww := luceio.NewLineWrappingWriter(
		luceio.LineWrapperContextWriter{
			Writer: w,
			Width:  width,
			Pad:    "// ",
		},
	)
	s0 := lww.Sum
	lww.WritePadding()
	lww.Write([]byte(c))
	lww.SumWriter.Write(nl)

	return lww.Sum - s0, lww.Err
}

// WriteComment is a helper. If it's given one comment value, it will write that
// if the comment is not blank. Two values is intended for doc comments. If the
// second argument is blank, no comment is written, other wise they are joined
// with a space. More than 3 comments and the comment slice is passed into Join.
func WriteComment(w io.Writer, pre Prefixer, comment ...string) (int64, error) {
	ln := len(comment)
	if ln == 0 || (ln == 1 && comment[0] == "") || (ln == 2 && comment[1] == "") {
		return 0, nil
	}
	if ln == 1 {
		return Comment(comment[0]).PrefixWriteTo(w, pre)
	}
	if ln == 2 {
		return Comment(luceio.Join(comment[0], comment[1], " ")).PrefixWriteTo(w, pre)
	}
	return Comment(luceio.Join(comment...)).PrefixWriteTo(w, pre)
}
