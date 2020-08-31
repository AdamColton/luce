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

// NewTypeRef fulfills PackageRef.
func (a Alias) NewTypeRef(name string, t Type) *TypeRef {
	return NewTypeRef(a, name, t)
}

// NewFuncRef creates a FuncRef in this Package.
func (a Alias) NewFuncRef(name string, args ...NameType) *FuncRef {
	return NewFuncRef(a, name, args...)
}

// NewInterfaceRef creates a Reference to an Interface in the Pacakge.
func (a Alias) NewInterfaceRef(name string) *InterfaceRef {
	return NewInterfaceRef(a, name)
}
