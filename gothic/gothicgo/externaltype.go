package gothicgo

import (
	"fmt"
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// ExternalType informs a gothic program of types available in external
// packages. It does not generate these but can use them. For instance
// ExternalType(MustPackageRef("time"), "Time") creates a reference to the Time
// type in the time package.
type ExternalType struct {
	Name string
	ExternalPackageRef
}

func (p externalPackageRef) ExternalType(name string) (*ExternalType, error) {
	if !IsExported(name) {
		return nil, fmt.Errorf(`ExternalType "%s" in package "%s" is not exported`, name, p.Name())
	}
	return &ExternalType{
		Name:               name,
		ExternalPackageRef: p,
	}, nil
}

func (p externalPackageRef) MustExternalType(name string) *ExternalType {
	et, err := p.ExternalType(name)
	lerr.Panic(err)
	return et
}

// PrefixWriteTo fulfills Type. Writes the ExternalType with prefixing.
func (e *ExternalType) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString(p.Prefix(e.ExternalPackageRef))
	sw.WriteString(e.Name)
	sw.Err = lerr.Wrap(sw.Err, "While writing external type %s", e.Name)
	return sw.Rets()
}

// PackageRef fulfills Type. Returns the ExternalPackageRef.
func (e *ExternalType) PackageRef() PackageRef { return e.ExternalPackageRef }

// Kind fulfills Type. Returns TypeDefKind
func (e *ExternalType) Kind() Kind {
	return TypeDefKind
}

// RegisterImports fulfills Type.
func (e *ExternalType) RegisterImports(i *Imports) {
	i.Add(e.ExternalPackageRef)
}
