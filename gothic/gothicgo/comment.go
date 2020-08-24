package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/util/luceio"
)

const defaultCommentWidth = 80

// Comment string that automatically wraps the string
type Comment struct {
	Comment string
	Width   int
}

// CommentWidther is anything that can specify a comment width
type CommentWidther interface {
	CommentWidth() int
}

// NewComment takes a string and returns a Comment. The comment width comes from
// the files package context. If File is nil a Comment with no width is
// returned. It does not write the comment to the file.
func (f *File) NewComment(comment string) *Comment {
	c := &Comment{
		Comment: comment,
		Width:   f.pkg.context.CommentWidth(),
	}
	f.AddWriterTo(c)
	return c
}

var nl = []byte("\n")

// WriteTo wraps the comment and writes it to the Writer
func (c *Comment) WriteTo(w io.Writer) (int64, error) {
	width := c.Width
	if width == 0 {
		width = defaultCommentWidth
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
	lww.Write([]byte(c.Comment))
	lww.SumWriter.Write(nl)

	return lww.Sum - s0, lww.Err
}
