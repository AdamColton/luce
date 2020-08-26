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
type ExternalType interface {
	Type
	ExternalPackageRef() ExternalPackageRef
}

func (p externalPackageRef) ExternalType(name string) (ExternalType, error) {
	if !IsExported(name) {
		return nil, fmt.Errorf(`ExternalType "%s" in package "%s" is not exported`, name, p.Name())
	}
	return &externalTypeWrapper{
		typeWrapper{
			&externalType{
				ref:  p,
				name: name,
			},
		},
	}, nil
}

func (p externalPackageRef) MustExternalType(name string) ExternalType {
	et, err := p.ExternalType(name)
	lerr.Panic(err)
	return et
}

type externalTypeWrapper struct {
	typeWrapper
}

func (e *externalTypeWrapper) ExternalPackageRef() ExternalPackageRef {
	return e.coreType.(*externalType).ref
}

type externalType struct {
	ref  ExternalPackageRef
	name string
}

// PrefixWriteTo writes the ExternalType handling prefixing
func (e *externalType) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString(p.Prefix(e.ref))
	sw.WriteString(e.name)
	sw.Err = lerr.Wrap(sw.Err, "While writing external type %s", e.name)
	return sw.Rets()
}

// PackageRef gets the PackageRef ExternalType belongs to
func (e *externalType) PackageRef() PackageRef { return e.ref }

// Kind of ExternalType is TypeDefKind
func (e *externalType) Kind() Kind {
	return TypeDefKind
}

func (e *externalType) RegisterImports(i *Imports) {
	i.Add(e.ref)
}
