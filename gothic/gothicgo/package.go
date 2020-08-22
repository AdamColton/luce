package gothicgo

import (
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adamcolton/luce/lerr"
)

// Package represents a directory containing Go code. Package also fulfills the
// PackageRef interface.
type Package struct {
	name       string
	importPath string
	OutputPath string
	context    Context
	files      map[string]*File
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
		name:       name,
		context:    ctx,
		files:      make(map[string]*File),
		importPath: ctx.ImportPath(),
		OutputPath: ctx.OutputPath(name),
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
	return strings.Join([]string{"\"", p.ImportPath(), "\""}, "")
}

func (*Package) privatePkgRef() {}
