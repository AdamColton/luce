package gothicgo

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// Package represents a directory containing Go code. Package also fulfills the
// PackageRef interface.
type Package struct {
	name       string
	importPath string
	OutputPath string
	context    Context
	files      map[string]*File
	names      map[string]Namer
	namers     map[Namer]string
	// CommentWidth is taken from the context when the Package is created.
	CommentWidth int
	// DefaultComment will be set on every File created. It is intended to be a
	// "Generated Code - Do Not Modify" style comment.
	DefaultComment string
}

// Namer represents a code generator that introduce a name to a package scope.
// The is used to prevent name collisions.
type Namer interface {
	ScopeName() string
}

const (
	// ErrBadPackageName is returned when a package name is not allowed
	ErrBadPackageName = lerr.Str("Bad package name")

	// ErrNilContext is returned if NewPackage is called with a nil Context
	ErrNilContext = lerr.Str("Context cannot be nil")
)

var nameRe = regexp.MustCompile(`^[\w\-]+$`)

// NewPackage creates a new Package. The import path will use the ImportPath
// set on the project. This should only be called by the Context. To create a
// Package call Context.Package.
func NewPackage(ctx Context, name string) (*Package, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	if !nameRe.MatchString(name) {
		return nil, ErrBadPackageName
	}
	pkg := &Package{
		name:           name,
		context:        ctx,
		files:          make(map[string]*File),
		importPath:     ctx.ImportPath(),
		OutputPath:     ctx.OutputPath(name),
		names:          make(map[string]Namer),
		namers:         make(map[Namer]string),
		DefaultComment: ctx.DefaultComment(),
		CommentWidth:   ctx.CommentWidth(),
	}
	ctx.AddPackage(pkg)
	return pkg, nil
}

// Prepare is currently a placeholder
func (p *Package) Prepare() error {
	for _, f := range p.files {
		err := f.Prepare()
		if err != nil {
			return lerr.Wrap(err, "Prepare package %s", p.name)
		}
	}
	return nil
}

// Generate is currently a placeholder
func (p *Package) Generate() error {
	path, _ := filepath.Abs(p.OutputPath)
	err := p.context.MakeDir(path)
	if err != nil {
		return lerr.Wrap(err, "Generate package %s", p.name)
	}
	for _, f := range p.files {
		err := f.Generate()
		if err != nil {
			return lerr.Wrap(err, "Generate package %s", p.name)
		}
	}
	return nil
}

// SetImportPath sets the import path for the package not including the name.
func (p *Package) SetImportPath(path string) error {
	if !importPathRe.MatchString(path) {
		return ErrBadImportPath
	}
	p.importPath = path
	return nil
}

// ImportPath returns the full import path including the package name.
func (p *Package) ImportPath() string {
	return path.Join(p.importPath, p.name)
}

// Name returns the package name and fulfills the PackageRef and Type
// interfaces.
func (p *Package) Name() string {
	return p.name
}

// ImportSpec returns import specification as it would be used in an import
// statement.
func (p *Package) ImportSpec() string {
	return luceio.Join("\"", p.ImportPath(), "\"", "")
}

func (*Package) privatePkgRef() {}

// AddNamer to Package checking for name scope collision
func (p *Package) AddNamer(n Namer) error {
	name := n.ScopeName()
	if name == "" {
		return nil
	}
	if _, ok := p.names[name]; ok {
		return fmt.Errorf("Name '%s' already exists in package '%s'", name, p.name)
	}
	p.names[name] = n
	p.namers[n] = name
	return nil
}

// UpdateNamer allows a Namer to change it's name within a package.
func (p *Package) UpdateNamer(n Namer) error {
	old, found := p.namers[n]
	if found {
		delete(p.names, old)
		delete(p.namers, n)
	}
	return p.AddNamer(n)
}

// NewTypeRef fulfills PackageRef.
func (p *Package) NewTypeRef(name string, t Type) *TypeRef {
	return NewTypeRef(p, name, t)
}

// NewFuncRef creates a FuncRef in this Package.
func (p *Package) NewFuncRef(name string, args ...NameType) *FuncRef {
	return NewFuncRef(p, name, args...)
}

func (p *Package) NewInterfaceRef(name string) *InterfaceRef {
	return NewInterfaceRef(p, name)
}
