package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/util/luceio"
)

// PackageVarRef is a reference to a variable defined in a package.
type PackageVarRef struct {
	NT  NameType
	Pkg PackageRef
}

// PrefixWriteTo fulfills PrefixWriterTo. Writes the PackageVarRef.
func (pv *PackageVarRef) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteStrings(
		pre.Prefix(pv.Pkg),
		pv.NT.N,
	)
	return sw.Rets()
}

// RegisterImports fulfills ImportRegistrar. Adds the VarRef Package to the
// imports.
func (pv *PackageVarRef) RegisterImports(i *Imports) {
	i.Add(pv.Pkg)
}

// PackageRef returns the PackageRef where the variable was be defined.
func (pv *PackageVarRef) PackageRef() PackageRef {
	return pv.Pkg
}
