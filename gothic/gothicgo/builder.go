package gothicgo

import (
	"fmt"
	"io"

	"github.com/adamcolton/luce/util/luceio"
)

// Builder is a helper for constructing generators.
type Builder struct {
	WriterTos []PrefixWriterTo
}

// NewBuilder will cast the writers to PrefixWriterTo and create a Builder. The
// writers must be PrefixWriterTo, io.WriterTo or string.
func NewBuilder(writers ...interface{}) (*Builder, error) {
	b := &Builder{
		WriterTos: make([]PrefixWriterTo, 0, len(writers)),
	}
	err := b.Append(writers...)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Append writers to the builder. The writers must be PrefixWriterTo,
// io.WriterTo or string.
func (b *Builder) Append(writers ...interface{}) error {
	for _, w := range writers {
		p, err := CastPrefixWriterTo(w)
		if err != nil {
			return err
		}
		b.WriterTos = append(b.WriterTos, p)
	}
	return nil
}

// CastPrefixWriterTo takes a PrefixWriterTo, io.WriterTo or string and converts
// it to a PrefixWriterTo.
func CastPrefixWriterTo(i interface{}) (PrefixWriterTo, error) {
	switch w := i.(type) {
	case PrefixWriterTo:
		return w, nil
	case io.WriterTo:
		return IgnorePrefixer{w}, nil
	case string:
		return IgnorePrefixer{luceio.StringWriterTo(w)}, nil
	}
	return nil, fmt.Errorf("Could not convert to PrefixWriterTo")
}

// PrefixWriteTo fulfills PrefixWriterTo. It call PrefixWriteTo on all
// PrefixWriterTos in the builder.
func (b *Builder) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	for _, p := range b.WriterTos {
		sumPrefixWriteTo(sw, pre, p)
	}
	return sw.Rets()
}

// RegisterImports fulfills ImportsRegistrar. It calls RegisterImports on the
// Type and the Value if it implements ImportsRegistrar.
func (b *Builder) RegisterImports(i *Imports) {
	for _, p := range b.WriterTos {
		if r, ok := p.(ImportsRegistrar); ok {
			r.RegisterImports(i)
		}
	}
}

// Layout is a helper that composites Builders into sections. This allows each
// section to be extended.
type Layout struct {
	Order    []string
	Builders map[string]*Builder
}

// NewLayout initilizes a Layout.
func NewLayout() *Layout {
	return &Layout{
		Builders: make(map[string]*Builder),
	}
}

// Section creates the section if it does not exist or appends the writers if it
// does.
func (l *Layout) Section(name string, writers ...interface{}) (*Builder, error) {
	b := l.Builders[name]
	if b != nil {
		return b, b.Append(writers...)
	}
	b, err := NewBuilder(writers...)
	if err != nil {
		return nil, err
	}
	l.Builders[name] = b
	l.Order = append(l.Order, name)
	return b, nil
}

// PrefixWriteTo fulfills PrefixWriterTo and calls PrefixWriteTo on each section
// in order.
func (l *Layout) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	for _, section := range l.Order {
		sumPrefixWriteTo(sw, pre, l.Builders[section])
	}
	return sw.Rets()
}

// RegisterImports calls RegisterImports on each section.
func (l *Layout) RegisterImports(i *Imports) {
	for _, b := range l.Builders {
		b.RegisterImports(i)
	}
}
