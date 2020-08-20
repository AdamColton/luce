package gothicgo

import (
	"path"
	"regexp"

	"github.com/adamcolton/luce/gothic"
	"github.com/adamcolton/luce/lerr"
)

// ErrBadImportPath indicates a poorly formatted import path. Path must end with
// /
const ErrBadImportPath = lerr.Str("Bad Import Path")

var importPathRe = regexp.MustCompile(`^[\w\-\.]+(\/[\w\-\.]+)*\/?$`)

// Context manages high level generation configuration.
type Context interface {
	gothic.Generator
	Export() error
	AddGenerators(g ...gothic.Generator) error

	Package(name string) (*Package, error)
	MustPackage(name string) *Package

	AddPackage(*Package) error
	OutputPath(name string) string
	SetOutputPath(path string)
	ImportPath() string
	SetImportPath(path string) error
}

// CtxFactory is used to configure the construction of a BaseContext.
type CtxFactory struct {
	OutputPath string
	ImportPath string
}

// New BaseContext access the os with the Factory settings.
func (c CtxFactory) New() *BaseContext {
	return &BaseContext{
		outputPath: c.OutputPath,
		importPath: c.ImportPath,
		project:    gothic.New(),
		packages:   make(map[string][]*Package),
	}
}

// BaseContext provides a Context for rendering Go code to the local file
// system. Swapping out the CreateFile, Abs and MkdirAll funcs allows it to be
// used for unit testing.
type BaseContext struct {
	outputPath string
	importPath string
	project    *gothic.Project
	packages   map[string][]*Package
}

// Prepare calls Prepare on all registered Generators
func (bc *BaseContext) Prepare() error {
	return bc.project.Prepare()
}

// Generate calls Generate on all registered Generators
func (bc *BaseContext) Generate() error {
	return bc.project.Generate()
}

// Export calls Prepare then Generate on all registered Generators
func (bc *BaseContext) Export() error {
	return bc.project.Export()
}

// AddGenerators to the Context.
func (bc *BaseContext) AddGenerators(g ...gothic.Generator) error {
	return bc.project.AddGenerators(g...)
}

// Package gets the package by name (matching the ImportPath) or creates a new
// Package.
func (bc *BaseContext) Package(name string) (*Package, error) {
	for _, p := range bc.packages[name] {
		if p.importPath == bc.importPath {
			return p, nil
		}
	}

	return NewPackage(bc, name)
}

// MustPackage gets the package by name (matching the ImportPath) or creates a
// new Package. If there is an error creating the package, MustPackage will
// panic.
func (bc *BaseContext) MustPackage(name string) *Package {
	p, err := bc.Package(name)
	lerr.Panic(err)
	return p
}

// AddPackage is called by NewPackage when a package is created. It should not
// be called directly.
func (bc *BaseContext) AddPackage(pkg *Package) error {
	n := pkg.Name()
	bc.packages[n] = append(bc.packages[n], pkg)

	return bc.AddGenerators(pkg)
}

// SetOutputPath changes the output path for a context. This is used to set the
// output path for a package. Packages cache the OutputPath at the time they
// are created, so the output path can be changed when creating packages.
func (bc *BaseContext) SetOutputPath(path string) {
	bc.outputPath = path
}

// OutputPath appends name to the end of the current OutputPath using path.Join.
func (bc *BaseContext) OutputPath(name string) string {
	return path.Join(bc.outputPath, name)
}

// SetImportPath for the project. It is safe to change import path during
// generation, anything that uses the default import path will get a copy at the
// time of it's instantiation.
func (bc *BaseContext) SetImportPath(path string) error {
	if !importPathRe.MatchString(path) {
		return ErrBadImportPath
	}
	bc.importPath = path
	return nil
}

// ImportPath returns the current import path
func (bc *BaseContext) ImportPath() string {
	return bc.importPath
}
