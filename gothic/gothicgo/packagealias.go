package gothicgo

import "github.com/adamcolton/luce/util/luceio"

// Alias a package import. Fulfils PackageRef.
type Alias struct {
	alias string
	ref   PackageRef
}

// NewAlias of the PackageRef to the given alias for import
func NewAlias(ref PackageRef, alias string) Alias {
	return Alias{
		alias: alias,
		ref:   ref,
	}
}

// ImportPath of the Alias, fulfils PackageRef.
func (a Alias) ImportPath() string {
	return a.ref.ImportPath()
}

// Name of the Alias, fulfils PackageRef.
func (a Alias) Name() string {
	return a.alias
}

// ImportSpec returns the import specification including the alias, fulfils
// PackageRef.
func (a Alias) ImportSpec() string {
	return luceio.Join(a.alias, " \"", a.ref.ImportPath(), "\"", "")
}

func (a Alias) privatePkgRef() {}
