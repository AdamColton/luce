package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// PackageVar is a Package level variable
type PackageVar struct {
	NT      NameType
	file    *File
	Value   PrefixWriterTo
	Comment string
}

// NewPackageVar generates a variable declaration.
func (f *File) NewPackageVar(name string, t Type) (*PackageVar, error) {
	pv := &PackageVar{
		NT:   NameType{name, t},
		file: f,
	}
	return pv, lerr.Wrap(f.AddGenerator(pv), "NewPackageVar")
}

// File returns the file that the variable declaration will be written to.
func (pv *PackageVar) File() *File {
	return pv.file
}

// PrefixWriteTo fulfills PrefixWriterTo. Writes the variable declaration.
func (pv *PackageVar) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	WriteComment(sw, pre, pv.NT.N, pv.Comment)
	sw.WriteString("var ")
	sumPrefixWriteTo(sw, pre, pv.NT)

	if pv.Value != nil {
		sw.WriteString(" = ")
		sumPrefixWriteTo(sw, pre, pv.Value)
	}

	sw.Err = lerr.Wrap(sw.Err, "While writing PackageVar %s", pv.NT.N)
	return sw.Rets()
}

// PackageRef returns the PackageRef where the variable will be defined.
func (pv *PackageVar) PackageRef() PackageRef {
	return pv.file.Package()
}

// ScopeName fulfills Namer and helps prevent name collision in a package.
func (pv *PackageVar) ScopeName() string {
	return pv.NT.N
}

// RegisterImports fulfills ImportsRegistrar. It calls RegisterImports on the
// Type and the Value if it implements ImportsRegistrar.
func (pv *PackageVar) RegisterImports(i *Imports) {
	pv.NT.RegisterImports(i)
	if r, ok := pv.Value.(ImportsRegistrar); ok {
		r.RegisterImports(i)
	}
}

// Ref returns a PackageVarRef referencing this PackageVar.
func (pv *PackageVar) Ref() *PackageVarRef {
	return &PackageVarRef{
		NT:  pv.NT,
		Pkg: pv.file.pkg,
	}
}
