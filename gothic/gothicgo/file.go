package gothicgo

import (
	"bytes"
	"go/format"
	"io"

	"github.com/adamcolton/luce/gothic"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// File represents a Go file. Writer is intended as a hook for testing. If it is
// nil, the code will be written to the file normally, if it's set to an
// io.WriteCloser, it will write to that instead
type File struct {
	*Imports
	name    string
	code    []io.WriterTo
	pkg     *Package
	Comment *Comment
	// CW is CommentWidth, that name is not used because it collides with the
	// method name.
	CW int
}

// File creates a file within the package. The name should not include ".go"
// which will be automatically appended.
func (p *Package) File(name string) *File {
	if file, exists := p.files[name]; exists {
		return file
	}
	f := &File{
		Imports: NewImports(p),
		name:    name,
		pkg:     p,
		CW:      p.CommentWidth,
	}
	if p.DefaultComment != "" {
		f.Comment = &Comment{
			Comment: p.DefaultComment,
			Width:   p.CommentWidth,
		}
	}
	p.files[name] = f
	return f
}

// Prepare runs prepare on all the generators in the file
func (f *File) Prepare() error {
	for _, w := range f.code {
		if p, ok := w.(gothic.Preparer); ok {
			err := p.Prepare()
			if err != nil {
				return lerr.Wrap(err, "While preparing file %s", f.name)
			}
		}
	}
	return nil
}

// AddWriterTo adds a WriterTo that will be invoked when the file is written. If
// the WriterTo fulfils gothic.Prepper then it's Prepare method will be called
// while File.Prepare is called. If the WriterTo fulfills Namer, it's ScopeName
// will be added to the package.
func (f *File) AddWriterTo(writerTo io.WriterTo) error {
	f.code = append(f.code, writerTo)
	n, isNamer := writerTo.(Namer)
	if !isNamer {
		return nil
	}
	return lerr.Wrap(f.pkg.AddNamer(n), "File.AddWriterTo")
}

// Generate the file
func (f *File) Generate() (err error) {
	f.Imports.RemoveRef(f.pkg)

	buf := bytes.NewBuffer(nil)
	sw := luceio.NewSumWriter(buf)

	if f.Comment != nil {
		f.Comment.WriteTo(sw)
	}
	sw.WriteRune('\n')
	sw.WriteString("package ")
	sw.WriteString(f.pkg.name)
	sw.WriteString("\n\n")
	_, err = f.Imports.WriteTo(sw)
	if err != nil {
		return lerr.Wrap(err, "WriteTo Error in Generate file %s/%s", f.pkg.name, f.name)
	}
	_, err = luceio.MultiWrite(sw, f.code, "\n\n")
	if sw.Err != nil {
		return lerr.Wrap(sw.Err, "Writer Error in Generate file %s/%s", f.pkg.name, f.name)
	}
	if err != nil {
		return lerr.Wrap(err, "WriteTo Error in Generate file %s/%s", f.pkg.name, f.name)
	}

	code := buf.Bytes()
	fmtCode, fmtErr := format.Source(code)

	w, err := f.pkg.context.FileWriter(f)
	if err != nil {
		return lerr.Wrap(err, "Getting writer for file %s/%s", f.pkg.name, f.name)
	}
	if wc, ok := w.(io.Closer); ok {
		defer func() {
			closeErr := wc.Close()
			if err == nil && closeErr != nil {
				err = lerr.Wrap(closeErr, "Closing writer for file %s/%s", f.pkg.name, f.name)
			}
		}()
	}

	if fmtErr == nil {
		_, err = w.Write(fmtCode)
		err = lerr.Wrap(err, "Writing file %s/%s", f.pkg.name, f.name)
	} else {
		_, err = w.Write(code)
		if err == nil {
			err = lerr.Wrap(fmtErr, "Failed to format %s/%s:", f.pkg.name, f.name)
		}
	}

	return
}

// Package the file is in
func (f *File) Package() *Package { return f.pkg }

// Name returns the name of the file.
func (f *File) Name() string { return f.name }

// UpdateNamer allows a Namer to change it's name within a package.
func (f *File) UpdateNamer(n Namer) error {
	return lerr.Wrap(f.Package().UpdateNamer(n), "UpdateNamer in file %s", f.name)
}

// CommentWidth fulfills CommentWidther. It preserves the comment width from
// the package at the time the File is created.
func (f *File) CommentWidth() int {
	return f.CW
}
