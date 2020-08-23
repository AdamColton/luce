package gothicgo

import (
	"fmt"
	"io"
	"strings"

	"github.com/adamcolton/luce/ds/bufpool"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// FuncSig is a function signature and fulfills Type.
type FuncSig interface {
	Type
	Name() string
	Args() []NameType
	Variadic() bool
	Returns() []NameType
	// AsType returns a FuncSig where the args and returns are unnamed. The
	// clearName options controls if the name is cleared.
	AsType(clearName bool) FuncSig
}

type funcSigT struct {
	name     string
	args     []NameType
	rets     []NameType
	variadic bool
}

// ErrMixedParameters is returned if a
const ErrMixedParameters = lerr.Str("Mixed named and unnamed function parameters")

// NewFuncSig returns a new function signature.
func NewFuncSig(name string, args, rets []NameType, variadic bool) FuncSig {
	return &funcSigHT{
		typeWrapper{
			&funcSigT{
				name:     name,
				args:     args,
				rets:     rets,
				variadic: variadic,
			},
		},
	}
}

func (f *funcSigT) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("func")
	if f.name != "" {
		sw.WriteRune(' ')
		sw.WriteString(f.name)
	}
	sw.WriteRune('(')
	var str string
	str, sw.Err = nameTypeSliceToString(pre, f.args, f.variadic)
	sw.WriteString(str)
	end := ""
	if len(f.rets) > 1 {
		sw.WriteString(") (")
		end = ")"
	} else {
		sw.WriteString(") ")
	}
	str, sw.Err = nameTypeSliceToString(pre, f.rets, false)
	sw.WriteString(str)
	sw.WriteString(end)
	sw.Err = lerr.Wrap(sw.Err, "While writing function signature %s", f.name)

	return sw.Rets()
}

func (f *funcSigT) Kind() Kind             { return FuncKind }
func (f *funcSigT) PackageRef() PackageRef { return pkgBuiltin }
func (f *funcSigT) RegisterImports(i *Imports) {
	for _, arg := range f.args {
		arg.RegisterImports(i)
	}
	for _, ret := range f.rets {
		ret.RegisterImports(i)
	}
}

func expectNamed(name string) error {
	if name == "" {
		return ErrMixedParameters
	}
	return nil
}

func expectUnnamed(name string) error {
	if name != "" {
		return ErrMixedParameters
	}
	return nil
}
func nameTypeSliceToString(pre Prefixer, nts []NameType, variadic bool) (string, error) {
	l := len(nts)
	if l == 0 {
		return "", nil
	}
	var validateName = expectNamed
	if nts[0].N == "" {
		validateName = expectUnnamed
	}

	var s = make([]string, l)
	l--
	var prevType string
	typeBuf := bufpool.Get()
	defer bufpool.Put(typeBuf)
	// working backwards makes it easier to place types
	for i := l; i >= 0; i-- {
		if err := validateName(nts[i].N); err != nil {
			return "", err
		}
		nts[i].Type.PrefixWriteTo(typeBuf, pre)
		if i == l && variadic {
			if nts[i].N == "" {
				s[i] = fmt.Sprintf("...%s", typeBuf.String())
			} else {
				s[i] = fmt.Sprintf("%s ...%s", nts[i].N, typeBuf.String())
			}
		} else if typeBuf.String() != prevType {
			if nts[i].N == "" {
				s[i] = fmt.Sprintf("%s", typeBuf.String())
			} else {
				s[i] = fmt.Sprintf("%s %s", nts[i].N, typeBuf.String())
				bs := make([]byte, typeBuf.Len())
				copy(bs, typeBuf.Bytes())
				prevType = string(bs)
			}
		} else {
			s[i] = nts[i].N
		}
		typeBuf.Reset()
	}

	return strings.Join(s, ", "), nil
}

type funcSigHT struct {
	typeWrapper
}

func (f *funcSigHT) Returns() []NameType {
	return f.coreType.(*funcSigT).rets
}

func (f *funcSigHT) Args() []NameType {
	return f.coreType.(*funcSigT).args
}

func (f *funcSigHT) Name() string {
	return f.coreType.(*funcSigT).name
}

func (f *funcSigHT) Variadic() bool {
	return f.coreType.(*funcSigT).variadic
}

func (f *funcSigHT) AsType(clearName bool) FuncSig {
	fs := f.coreType.(*funcSigT)
	args := make([]NameType, len(fs.args))
	rets := make([]NameType, len(fs.rets))
	for i, a := range fs.args {
		args[i].Type = a.Type
	}
	for i, r := range fs.rets {
		rets[i].Type = r.Type
	}
	var name string
	if !clearName {
		name = fs.name
	}
	return NewFuncSig(name, args, rets, fs.variadic)
}
