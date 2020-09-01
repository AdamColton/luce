package gothicgo

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// StructEmbeddable is used to embed a named type in a struct. The returned
// string is what the field name will be. So when embedding *foo.Bar, the
// StructEmbedName will be Bar.
type StructEmbeddable interface {
	Type
	StructEmbedName() string
}

// StructType represents a Go struct.
type StructType struct {
	fields     map[string]*Field
	fieldOrder []string
}

// NewStructType defines a struct
func NewStructType(fields ...PrefixWriterTo) (*StructType, error) {
	s := &StructType{
		fields: make(map[string]*Field),
	}
	_, err := s.AddFields(fields...)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// MustStructType defines a struct and panics if it fails
func MustStructType(fields ...PrefixWriterTo) *StructType {
	s, err := NewStructType(fields...)
	lerr.Panic(err)
	return s
}

// PackageRef gets the name of the package.
func (s *StructType) PackageRef() PackageRef { return pkgBuiltin }

// Field returns a field by name
func (s *StructType) Field(name string) (*Field, bool) {
	f, ok := s.fields[name]
	return f, ok
}

// RegisterImports fulfills ImportsRegistrar. Register imports on all fields.
func (s *StructType) RegisterImports(i *Imports) {
	for _, f := range s.fields {
		f.NameType.T.RegisterImports(i)
	}
}

// Fields returns the fields in order.
func (s *StructType) Fields() []string {
	fs := make([]string, len(s.fieldOrder))
	copy(fs, s.fieldOrder)
	return fs
}

// FieldCount returns how many fields the struct has
func (s *StructType) FieldCount() int {
	return len(s.fieldOrder)
}

const (
	// ErrBadField is returned if a PrefixWriterTo cannot be converted to a field.
	ErrBadField = lerr.Str("Given type cannot be converted to struct field")
	// ErrBadFieldName is returned if the Field name is malformed.
	ErrBadFieldName = lerr.Str("Field must either be named or type must be StructEmbeddable")
)

// AddFields to the struct. Fields must be *Field, NameType, StructEmbeddable.
func (s *StructType) AddFields(fields ...PrefixWriterTo) ([]*Field, error) {
	out := make([]*Field, 0, len(fields))
	for _, p := range fields {
		f, err := s.AddField(p)
		if err != nil {
			return out, err
		}
		out = append(out, f)
	}
	return out, nil
}

// AddField to struct. Field must be *Field, NameType, StructEmbeddable.
func (s *StructType) AddField(field PrefixWriterTo, tags ...string) (*Field, error) {
	f, err := NewField(field, tags...)
	if err != nil {
		return nil, err
	}

	key := f.Name()
	if key == "" {
		return nil, ErrBadFieldName
	}
	if _, exists := s.fields[key]; exists {
		return nil, fmt.Errorf(`Field "%s" already exists in struct`, key)
	}
	s.fields[key] = f
	s.fieldOrder = append(s.fieldOrder, key)
	return f, nil
}

// PrefixWriteTo fulfills PrefixWriterTo. Writes the Struct to the writer.
func (s *StructType) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sum := luceio.NewSumWriter(w)
	if len(s.fieldOrder) == 0 {
		sum.WriteString("struct{}")
		return sum.Rets()
	}
	sum.WriteString("struct {")
	for _, f := range s.fieldOrder {
		sum.WriteString("\n\t")
		s.fields[f].PrefixWriteTo(sum, p)
	}
	sum.WriteString("\n}")
	sum.Err = lerr.Wrap(sum.Err, "While writing struct")
	return sum.Rets()
}

// Field is a struct field. Tags follows the convention of `key1:"value1"
// key2:"value2"`. If no value is defined only the key is printed.
type Field struct {
	NameType
	Tags map[string]string
}

// NewField creates a field from the Name and Type
func NewField(field PrefixWriterTo, tags ...string) (*Field, error) {
	var f *Field
	switch t := field.(type) {
	case *Field:
		f = t
	case NameType:
		f = &Field{NameType: t}
	case StructEmbeddable:
		f = &Field{NameType: NameType{"", t}}
	default:
		return nil, ErrBadField
	}

	f.AddTags(tags...)
	return f, nil
}

// Name of the field. If the field is explicit, it will be the defined name, if
// the field is embedded it will the embedded name.
func (f *Field) Name() string {
	if f.N != "" {
		return f.N
	}
	if emb, ok := f.T.(StructEmbeddable); ok {
		return emb.StructEmbedName()
	}
	return ""
}

// AddTags to the field
func (f *Field) AddTags(tags ...string) {
	if len(tags) == 0 {
		return
	}
	if len(tags)%2 == 1 {
		tags = append(tags, "")
	}
	if f.Tags == nil {
		f.Tags = make(map[string]string, len(tags)/2)
	}
	for i := 0; i < len(tags); i += 2 {
		f.AddTag(tags[i], tags[i+1])
	}
}

// AddTag to the field
func (f *Field) AddTag(key, value string) {
	if f.Tags == nil {
		f.Tags = map[string]string{
			key: value,
		}
		return
	}
	if s, ok := f.Tags[key]; ok {
		f.Tags[key] = s + ";" + value
		return
	}
	f.Tags[key] = value
}

// PrefixWriteTo fulfills PrefixWriterTo. Writes the Field to the writer.
func (f *Field) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sum := luceio.NewSumWriter(w)
	if f.N != "" {
		sum.WriteString(f.N)
		sum.WriteString(" ")
	}
	indentWriter := luceio.ReplacerWriter{
		Writer:   sum,
		Replacer: strings.NewReplacer("\n", "\n\t"),
	}
	f.T.PrefixWriteTo(indentWriter, p)

	if len(f.Tags) > 0 {
		sum.WriteString(" `")
		tags := make([]string, 0, len(f.Tags))
		for k := range f.Tags {
			tags = append(tags, k)
		}
		sort.Strings(tags)
		for i, tag := range tags {
			if i > 0 {
				sum.WriteString(" ")
			}
			sum.WriteString(tag)
			if v := f.Tags[tag]; v != "" {
				sum.WriteString(":\"")
				sum.WriteString(v)
				sum.WriteString("\"")
			}
		}
		sum.WriteString("`")
	}

	sum.Err = lerr.Wrap(sum.Err, "While writing field %s", f.Name())

	return sum.Rets()
}
