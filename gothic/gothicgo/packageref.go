package gothicgo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/adamcolton/luce/lerr"
)

// PackageRef represents a reference to a package.
type PackageRef interface {
	ImportPath() string
	Name() string

	// ImportSpec is the specification and may include a package name or a
	// modifier like _.
	ImportSpec() string

	// PackageRef is not meant to be implemented, it's meant as an accessor to the
	// underlying packageRef. All instances should be created with NewPackageRef
	// to guarentee that the reference is well formed.
	privatePkgRef()
}

// ExternalPackageRef represents an external package - one that will not be
// generated.
type ExternalPackageRef interface {
	PackageRef

	// ExternalType represents a type in an external package. The name must
	// be exported (begin with an uppercase character).
	ExternalType(name string) (ExternalType, error)
	MustExternalType(name string) ExternalType
}

var pkgBuiltin = &packageRef{}

// PkgBuiltin returns a PackageRef with an empty name.
func PkgBuiltin() PackageRef {
	return pkgBuiltin
}

type packageRef struct {
	path, name string
}

func (p packageRef) ImportPath() string {
	return p.path
}

func (p packageRef) Name() string {
	return p.name
}

func (p packageRef) ImportSpec() string {
	return strings.Join([]string{"\"", p.path, "\""}, "")
}

func (packageRef) privatePkgRef() {}

type externalPackageRef struct {
	packageRef
}

func (p externalPackageRef) ExternalType(name string) (ExternalType, error) {
	if !IsExported(name) {
		return nil, fmt.Errorf(`ExternalType "%s" in package "%s" is not exported`, name, p.Name())
	}
	return &externalTypeWrapper{
		HelpfulTypeWrapper{
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

// TODO: this regex is only mostly right
var packageRefRegex = regexp.MustCompile(`^(?:[\w\-\.]+\/)*([\w\-]+)$`)

// ErrBadPackageRef indicates a poorly formatted package ref string.
const ErrBadPackageRef = lerr.Str("Bad Package Ref")

// NewExternalPackageRef takes the string used to import a pacakge and returns
// an ExternalPackageRef.
func NewExternalPackageRef(ref string) (ExternalPackageRef, error) {
	m := packageRefRegex.FindStringSubmatch(ref)
	if len(m) == 0 {
		return nil, ErrBadPackageRef
	}
	return externalPackageRef{packageRef{
		path: m[0],
		name: m[1],
	}}, nil
}

// MustExternalPackageRef returns a new PackageRef and panics if there is an error
func MustExternalPackageRef(ref string) ExternalPackageRef {
	p, err := NewExternalPackageRef(ref)
	lerr.Panic(err)
	return p
}
