package gothicgo

import (
	"fmt"
	"io"
	"strings"

	"github.com/adamcolton/luce/ds/bufpool"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

const (
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

type NewTypeDefiner interface {
	NewTypeDef(name string, t Type) (*TypeDef, error)
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

func (td *TypeDef) Prepare() error {
	td.baseType.RegisterImports(td.File().Imports)
	return nil
}

func (td *TypeDef) WriteTo(w io.Writer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	if td.Comment != "" {
		sw.WriterTo(&Comment{
			Comment: luceio.Join(td.Name(), td.Comment, " "),
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
func (td *TypeDef) PackageRef() PackageRef {
	return td.file.Package()
}
func (td *TypeDef) File() *File {
	return td.file
}

func (td *TypeDef) Name() string {
	return td.name
}
func (td *TypeDef) ScopeName() string {
	return td.name
}

func (td *TypeDef) RegisterImports(i *Imports) {
	i.Add(td.file.Package())
}

func (td *TypeDef) StructEmbedName() string {
	return td.name
}

func (td *TypeDef) Base() Type {
	return td.Elem()
}

func (td *TypeDef) Elem() Type {
	if td.Ptr {
		return PointerTo(td.baseType)
	}
	return td.baseType
}

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
	fn := &Func{
		FuncSig: NewFuncSig(name, args...),
		file:    td.file,
	}
	m := &Method{
		typeDef: td,
		Ptr:     td.Ptr,
		Func:    fn,
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

// Method on a struct
type Method struct {
	*Func
	Ptr     bool
	typeDef *TypeDef
}

// SetName of the method, also updates the method map in the struct.
func (m *Method) SetName(name string) {
	delete(m.typeDef.methods, m.Func.Name)
	m.Rename(name)
	m.typeDef.methods[name] = m
}

// ScopeName overrides ScopeName on Func, not the best solution.
func (m *Method) ScopeName() string { return "" }

// String outputs the entire function as a string
func (m *Method) String() string {
	buf := bufpool.Get()
	m.WriteTo(buf)
	return bufpool.PutStr(buf)
}

// WriteTo writes the Method to the writer
func (m *Method) WriteTo(w io.Writer) (int64, error) {
	if m.Name == "" {
		return 0, lerr.Str("Cannot have unnamed method")
	}
	sw := luceio.NewSumWriter(w)
	if m.Comment != "" {
		sw.WriterTo(&Comment{
			Comment: luceio.Join(m.Name, m.Comment, " "),
			Width:   m.file.CommentWidth(),
		})
	}

	sw.WriteStrings("func (", m.typeDef.ReceiverName)
	if m.Ptr {
		sw.WriteString(" *")
	} else {
		sw.WriteRune(' ')
	}
	sw.WriteStrings(m.typeDef.Name(), ") ", m.Name, "(")
	var str string
	str, sw.Err = nameTypeSliceToString(m.typeDef.file, m.Args, m.Variadic)
	sw.WriteString(str)
	end := " {\n\t"
	if ln := len(m.Rets); ln > 1 {
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

	if m.Func.Body != nil {
		m.Func.Body.PrefixWriteTo(sw, m.typeDef.file)
	}
	sw.WriteString("\n}")
	sw.Err = lerr.Wrap(sw.Err, "While writing method %s", m.Name)

	return sw.Rets()
}

func (m *Method) Receiver() NameType {
	n := NameType{
		N: m.typeDef.ReceiverName,
	}
	if m.Ptr {
		n.T = m.typeDef.Pointer()
	} else {
		n.T = m.typeDef
	}
	return n
}

// Returns sets the return types on the function
func (m *Method) Returns(rets ...NameType) *Method {
	m.Func.Returns(rets...)
	return m
}

func (m *Method) UnnamedRets(rets ...Type) *Method {
	m.Func.UnnamedRets(rets...)
	return m
}
