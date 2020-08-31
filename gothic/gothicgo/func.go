package gothicgo

import (
	"io"
	"strings"

	"github.com/adamcolton/luce/ds/bufpool"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// Func function written to a Go file
type Func struct {
	*FuncSig
	Body PrefixWriterTo
	// Comment will automatically be prefixed with the Func name.
	Comment string
	file    *File
}

// FuncCaller returns the string of a function invocation with the provided
// args. Intended to be used when building code generators.
type FuncCaller interface {
	Call(pre Prefixer, args ...string) string
}

// ErrUnnamedFuncArg is returned from NewFunc if an unnamed NameType is used
// as an arg.
const ErrUnnamedFuncArg = lerr.Str("All func args must be nammed")

// NewFunc returns a new Func with File set and add the function to file's
// generators so that when the file is generated, the func will be generated as
// part of the file.
func (f *File) NewFunc(name string, args ...NameType) (*Func, error) {
	for _, arg := range args {
		if arg.N == "" {
			return nil, ErrUnnamedFuncArg
		}
	}

	fn := &Func{
		FuncSig: NewFuncSig(name, args...),
		file:    f,
	}
	return fn, lerr.Wrap(f.AddGenerator(fn), "File.NewFunc")
}

// MustFunc calls NewFunc and panics if there is an error
func (f *File) MustFunc(name string, args ...NameType) *Func {
	fn, err := f.NewFunc(name, args...)
	lerr.Panic(err)
	return fn
}

// ScopeName fulfills Namer registering the function name with the package.
func (f *Func) ScopeName() string { return f.Name }

// PrefixWriteTo fulfilss PrefixWriterTo. It generates the function to the
// writer using the prefixer.
func (f *Func) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	WriteComment(sw, pre, f.Name, f.Comment)
	f.FuncSig.PrefixWriteTo(w, pre)
	sw.WriteString(" {\n")
	if f.Body != nil {
		f.Body.PrefixWriteTo(sw, pre)
	}
	sw.WriteString("\n}")
	sw.Err = lerr.Wrap(sw.Err, "While writing func %s", f.Name)
	return sw.Rets()
}

// BodyWriterTo is a helper allowing the body to be set to a Writer that
// ignores the prefixer.
func (f *Func) BodyWriterTo(w io.WriterTo) *Func {
	f.Body = IgnorePrefixer{w}
	return f
}

// BodyString is a helper allowing the body to be set to a string that
// ignores the prefixer.
func (f *Func) BodyString(str string) *Func {
	f.Body = IgnorePrefixer{luceio.StringWriterTo(str)}
	return f
}

// Call produces a invocation of the function and fulfills the FuncCaller
// interface
func (f *Func) Call(pre Prefixer, args ...string) string {
	return funcCall(pre, f.Name, args, f.file.Package())
}

// Rename the function and update the name in the package.
func (f *Func) Rename(name string) error {
	f.Name = name
	return f.file.pkg.UpdateNamer(f)
}

// File returns the File the function will be written to.
func (f *Func) File() *File {
	return f.file
}

// RegisterImports fulfills ImportsRegistrar. It registers the types from the
// arguments and return values. If the Body implements ImportsRegistrar, it will
// also be invoked.
func (f *Func) RegisterImports(i *Imports) {
	f.FuncSig.RegisterImports(i)
	if ri, ok := f.Body.(ImportsRegistrar); ok {
		ri.RegisterImports(i)
	}
}

func funcCall(pre Prefixer, name string, args []string, pkg PackageRef) string {
	buf := bufpool.Get()
	buf.WriteString(pre.Prefix(pkg))
	buf.WriteString(name)
	buf.WriteRune('(')
	buf.WriteString(strings.Join(args, ", "))
	buf.WriteRune(')')
	return bufpool.PutStr(buf)
}

// Returns sets the return types on the function
func (f *Func) Returns(rets ...NameType) *Func {
	f.FuncSig.Returns(rets...)
	return f
}

// UnnamedRets sets the return types on the function
func (f *Func) UnnamedRets(rets ...Type) *Func {
	f.FuncSig.UnnamedRets(rets...)
	return f
}

// Ref returns a FuncRef, allowing the function to be used as a Type.
func (f *Func) Ref() *FuncRef {
	return &FuncRef{
		FuncSig: f.FuncSig,
		Pkg:     f.file.pkg,
	}
}
