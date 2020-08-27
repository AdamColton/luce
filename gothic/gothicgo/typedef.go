package gothicgo

import (
	"fmt"
	"io"
	"strings"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

const (
	// ErrTypeDefOnIface if an InterfaceType is passed into a TypeDef constructor.
	ErrTypeDefOnIface = lerr.Str("Cannot create TypeDef on interface")
)

// TypeDef within a file. Ptr defaults to true if the Type is a Struct.
type TypeDef struct {
	baseType     Type
	name         string
	file         *File
	ReceiverName string
	methods      map[string]*Method

	// Ptr determines if the TypeDef should be primarily treated as a pointer or
	// a literal. This changes the receiver type when defining methods and how
	// it will be referenced.
	Ptr     bool
	Comment string
}

// NewTypeDef adds a type to a Context, the package and file for the struct is
// automatically generated.
func (c *BaseContext) NewTypeDef(name string, t Type) (*TypeDef, error) {
	p, err := c.Package(strings.ToLower(name))
	if err != nil {
		return nil, err
	}
	return p.NewTypeDef(name, t)
}

// MustTypeDef calls NewTypeDef and panics if there is an error
func (c *BaseContext) MustTypeDef(name string, t Type) *TypeDef {
	td, err := c.NewTypeDef(name, t)
	lerr.Panic(err)
	return td
}

// NewTypeDef adds a type to a Package, the file for the struct is automatically
// generated.
func (p *Package) NewTypeDef(name string, t Type) (*TypeDef, error) {
	if _, ok := p.names[name]; ok {
		return nil, fmt.Errorf("Cannot define type %s in package %s; name already exists in scope", name, p.Name())
	}
	if t.Kind() == InterfaceKind {
		return nil, ErrTypeDefOnIface
	}
	return p.File(name+".gothic").NewTypeDef(name, t)
}

// MustTypeDef calls NewTypeDef and panics if there is an error.
func (p *Package) MustTypeDef(name string, t Type) *TypeDef {
	td, err := p.NewTypeDef(name, t)
	lerr.Panic(err)
	return td
}

// NewTypeDef adds a type to a file
func (f *File) NewTypeDef(name string, t Type) (*TypeDef, error) {
	if _, ok := f.pkg.names[name]; ok {
		return nil, fmt.Errorf("Cannot define type %s in package %s; name already exists in scope", name, f.pkg.Name())
	}
	if t.Kind() == InterfaceKind {
		return nil, ErrTypeDefOnIface
	}
	td := &TypeDef{
		baseType:     t,
		name:         name,
		file:         f,
		methods:      make(map[string]*Method),
		ReceiverName: strings.ToLower(string([]rune(name)[0])),
		Ptr:          t.Kind() == StructKind,
	}
	return td, f.AddWriterTo(td)
}

// MustTypeDef calls NewTypeDef and panics if there is an error.
func (f *File) MustTypeDef(name string, t Type) *TypeDef {
	td, err := f.NewTypeDef(name, t)
	lerr.Panic(err)
	return td
}

// RegisterImports from the base Type.
func (td *TypeDef) RegisterImports(i *Imports) {
	td.baseType.RegisterImports(td.File().Imports)
}

// WriteTo writes the TypeDef.
func (td *TypeDef) WriteTo(w io.Writer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	if td.Comment != "" {
		sw.WriterTo(&Comment{
			Comment: luceio.Join(td.name, td.Comment, " "),
			Width:   td.file.CommentWidth(),
		})
	}
	sw.WriteString("type ")
	sw.WriteString(td.name)
	sw.WriteRune(' ')
	td.baseType.PrefixWriteTo(sw, td.file)
	sw.Err = lerr.Wrap(sw.Err, "While writing type %s", td.name)
	return sw.Rets()
}

// PrefixWriteTo fulfills type.
func (td *TypeDef) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	if td.Ptr {
		sw.WriteRune('*')
	}
	sw.WriteString(p.Prefix(td.file.Package()))
	sw.WriteString(td.name)
	sw.Err = lerr.Wrap(sw.Err, "While writing type %s", td.name)
	return sw.Rets()
}

// PackageRef fulfills Type.
func (td *TypeDef) PackageRef() PackageRef {
	return td.file.Package()
}

// File TypeDef will be written to.
func (td *TypeDef) File() *File {
	return td.file
}

// ScopeName fulfills Namer.
func (td *TypeDef) ScopeName() string {
	return td.name
}

// Elem returns the base type.
func (td *TypeDef) Elem() Type {
	return td.baseType
}

// Method on a struct
type Method struct {
	*FuncSig
	Comment string
	Body    PrefixWriterTo
	Ptr     bool
	typeDef *TypeDef
}

// MustMethod calls NewMethod and panics if there is an error.
func (td *TypeDef) MustMethod(name string, args ...NameType) *Method {
	m, err := td.NewMethod(name, args...)
	lerr.Panic(err)
	return m
}

// NewMethod on the struct
func (td *TypeDef) NewMethod(name string, args ...NameType) (*Method, error) {
	if name == "" {
		return nil, lerr.Str("Cannot have unnamed method")
	}
	m := &Method{
		typeDef: td,
		Ptr:     td.Ptr,
		FuncSig: NewFuncSig(name, args...),
	}
	err := td.File().AddWriterTo(m)
	if err != nil {
		return nil, err
	}
	td.methods[name] = m
	return m, nil
}

// Method gets a method by name
func (td *TypeDef) Method(name string) (*Method, bool) {
	m, ok := td.methods[name]
	return m, ok
}

// ErrUnnamedMethod is returned if a Method is created or set to an empty name.
const ErrUnnamedMethod = lerr.Str("Cannot have unnamed method")

// WriteTo writes the Method to the writer
func (m *Method) WriteTo(w io.Writer) (int64, error) {
	if m.Name == "" {
		return 0, ErrUnnamedMethod
	}
	sw := luceio.NewSumWriter(w)
	if m.Comment != "" {
		sw.WriterTo(&Comment{
			Comment: luceio.Join(m.Name, m.Comment, " "),
			Width:   m.typeDef.file.CommentWidth(),
		})
	}
	sw.WriteStrings("func (")
	m.Receiver().PrefixWriteTo(sw, m.typeDef.file)
	sw.WriteStrings(") ", m.Name, "(")
	var str string
	str, sw.Err = nameTypeSliceToString(m.typeDef.file, m.Args, m.Variadic)
	sw.WriteString(str)
	end := " {\n\t"
	if ln := len(m.Rets); ln > 1 || (len(m.Rets) > 0 && m.Rets[0].N != "") {
		sw.WriteString(") (")
		end = ") {\n\t"
	} else {
		sw.WriteString(")")
		if ln == 1 {
			sw.WriteString(" ")
		}
	}
	str, sw.Err = nameTypeSliceToString(m.typeDef.file, m.Rets, false)
	sw.WriteStrings(str, end)

	if m.Body != nil {
		m.Body.PrefixWriteTo(sw, m.typeDef.file)
	}
	sw.WriteString("\n}")
	sw.Err = lerr.Wrap(sw.Err, "While writing method %s", m.Name)

	return sw.Rets()
}

// Receiver returns the NameType of the method receiver.
func (m *Method) Receiver() NameType {
	n := NameType{
		N: m.typeDef.ReceiverName,
		T: m.typeDef,
	}
	if m.Ptr {
		n.T = m.typeDef.Pointer()
	}
	return n
}

// Returns sets the return types on the function
func (m *Method) Returns(rets ...NameType) *Method {
	m.FuncSig.Returns(rets...)
	return m
}

// UnnamedRets sets the return types on the function.
func (m *Method) UnnamedRets(rets ...Type) *Method {
	m.FuncSig.UnnamedRets(rets...)
	return m
}

// BodyWriterTo is a helper allowing the body to be set to a Writer that
// ignores the prefixer.
func (m *Method) BodyWriterTo(w io.WriterTo) *Method {
	m.Body = IgnorePrefixer{w}
	return m
}

// BodyString is a helper allowing the body to be set to a string that
// ignores the prefixer.
func (m *Method) BodyString(str string) *Method {
	m.Body = IgnorePrefixer{luceio.StringWriterTo(str)}
	return m
}

// ErrMethodAlreadyExists is return if a method name is redeclared on a TypeDef.
const ErrMethodAlreadyExists = lerr.Str("Method already exists")

// Rename checks for collision, changes the name and updates the name index.
func (m *Method) Rename(name string) error {
	if name == "" {
		return ErrUnnamedMethod
	}
	if _, found := m.typeDef.methods[name]; found {
		return ErrMethodAlreadyExists
	}
	delete(m.typeDef.methods, m.FuncSig.Name)
	m.FuncSig.Name = name
	m.typeDef.methods[name] = m
	return nil
}

// RegisterImports from the function signature and checks if the body is an
// ImportsRegistrar.
func (m *Method) RegisterImports(i *Imports) {
	m.FuncSig.RegisterImports(i)
	if r, ok := m.Body.(ImportsRegistrar); ok {
		r.RegisterImports(i)
	}
}
