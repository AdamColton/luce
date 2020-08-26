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
type FuncSig struct {
	Name     string
	Args     []NameType
	Rets     []NameType
	Variadic bool
}

// ErrMixedParameters is returned if a
const ErrMixedParameters = lerr.Str("Mixed named and unnamed function parameters")

// NewFuncSig returns a new function signature.
func NewFuncSig(name string, args ...NameType) *FuncSig {
	return &FuncSig{
		Name: name,
		Args: args,
	}
}

func (f *FuncSig) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("func")
	if f.Name != "" {
		sw.WriteRune(' ')
		sw.WriteString(f.Name)
	}
	sw.WriteRune('(')
	var str string
	str, sw.Err = nameTypeSliceToString(pre, f.Args, f.Variadic)
	sw.WriteString(str)
	end := ""
	if len(f.Rets) > 1 {
		sw.WriteString(") (")
		end = ")"
	} else {
		sw.WriteString(") ")
	}
	str, sw.Err = nameTypeSliceToString(pre, f.Rets, false)
	sw.WriteString(str)
	sw.WriteString(end)
	sw.Err = lerr.Wrap(sw.Err, "While writing function signature %s", f.Name)

	return sw.Rets()
}

func (f *FuncSig) PackageRef() PackageRef { return pkgBuiltin }
func (f *FuncSig) RegisterImports(i *Imports) {
	for _, arg := range f.Args {
		arg.T.RegisterImports(i)
	}
	for _, ret := range f.Rets {
		ret.T.RegisterImports(i)
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
	if str := typeBuf.String(); str != "" {
		fmt.Println(str)
	}
	// working backwards makes it easier to place types
	for i := l; i >= 0; i-- {
		if err := validateName(nts[i].N); err != nil {
			return "", err
		}
		nts[i].T.PrefixWriteTo(typeBuf, pre)
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

func (f *FuncSig) AsType(clearName bool) *FuncSig {
	var name string
	if !clearName {
		name = f.Name
	}
	return &FuncSig{
		Name:     name,
		Args:     ClearNames(f.Args...),
		Rets:     ClearNames(f.Rets...),
		Variadic: f.Variadic,
	}
}

// Returns sets the return types on the function
func (f *FuncSig) Returns(rets ...NameType) *FuncSig {
	f.Rets = rets
	return f
}

func (f *FuncSig) UnnamedRets(rets ...Type) *FuncSig {
	f.Rets = make([]NameType, len(rets))
	for i, t := range rets {
		f.Rets[i].T = t
	}
	return f
}
