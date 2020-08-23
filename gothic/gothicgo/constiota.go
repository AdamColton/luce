package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// ConstIotaBlock creates a block of const values of a given type using an iota
// expression.
type ConstIotaBlock struct {
	t       Type
	rows    []string
	Comment string
	Iota    string
	file    interface {
		AddWriterTo(io.WriterTo) error
		Package() *Package
		Prefixer
		CommentWidth() int
	}
}

const (
	// ErrEmptyConstIotaBlock is returned if ConstIotaBlock.WriteTo is called on
	// an empty ConstIotaBlock.
	ErrEmptyConstIotaBlock = lerr.Str("ConstIotaBlock requires at least one row")

	// ErrConstIotaBlockType is returned if ConstIotaBlock.WriteTo is called on
	// a ConstIotaBlock with no type defined.
	ErrConstIotaBlockType = lerr.Str("ConstIotaBlock requires a type")
)

type constRow string

func (c constRow) ScopeName() string {
	return string(c)
}

// WriteTo renders the ConstIotaBlock to a writer
func (cib *ConstIotaBlock) WriteTo(w io.Writer) (int64, error) {
	if len(cib.rows) == 0 {
		return 0, ErrEmptyConstIotaBlock
	}
	if cib.t == nil {
		return 0, ErrConstIotaBlockType
	}
	sw := luceio.NewSumWriter(w)
	if cib.Comment != "" {
		(&Comment{
			Comment: cib.Comment,
			Width:   cib.file.CommentWidth(),
		}).WriteTo(sw)
	}
	sw.WriteString("const (\n\t")
	sw.WriteString(cib.rows[0])
	sw.WriteRune(' ')
	cib.t.PrefixWriteTo(sw, cib.file)
	sw.WriteString(" = ")
	if cib.Iota == "" {
		sw.WriteString("iota")
	} else {
		sw.WriteString(cib.Iota)
	}
	for _, r := range cib.rows[1:] {
		sw.WriteString("\n\t")
		sw.WriteString(r)
	}
	sw.WriteString("\n)\n")
	if sw.Err != nil {
		sw.Err = lerr.Wrap(sw.Err, "While writing ConstIotaBlock")
	}
	return sw.Rets()
}

// Append rows to a ConstIotaBlock.
func (cib *ConstIotaBlock) Append(rows ...string) error {
	pkg := cib.file.Package()
	for _, r := range rows {
		err := pkg.AddNamer(constRow(r))
		if err != nil {
			return lerr.Wrap(err, "ConstIotaBlock")
		}
	}
	cib.rows = append(cib.rows, rows...)
	return nil
}

// NewConstIotaBlock creates a ConstIotaBlock on the given File.
func (f *File) NewConstIotaBlock(t Type, rows ...string) (*ConstIotaBlock, error) {
	if t == nil {
		return nil, ErrConstIotaBlockType
	}
	cib := &ConstIotaBlock{
		t:    t,
		file: f,
		rows: rows,
	}
	err := cib.Append(rows...)
	if err != nil {
		return nil, lerr.Wrap(err, "New")
	}

	return cib, lerr.Wrap(f.AddWriterTo(cib), "NewConstIotaBlock")
}

// MustConstIotaBlock creates a new ConstIotaBlock and panics if there is an
// error.
func (f *File) MustConstIotaBlock(t Type, rows ...string) *ConstIotaBlock {
	c, err := f.NewConstIotaBlock(t, rows...)
	lerr.Panic(err)
	return c
}
