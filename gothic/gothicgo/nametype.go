package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// NameType is used for arguments and returns for function. Type is specifically
// not embedded because doing so causes unexpected results when the NameType
// PrefixWriteTo is invoked when it's being treated as a type.
type NameType struct {
	N string
	T Type
}

// Name value
func (n NameType) Name() string { return n.N }

// Unnamed takes a slice of types and returns them as a slice of NameTypes that are
// unnamed.
func Unnamed(ts ...Type) []NameType {
	nts := make([]NameType, len(ts))
	for i, t := range ts {
		nts[i].T = t
	}
	return nts
}

// ClearNames takes a slice of NameType and returns a matching slice with all
// the names cleared.
func ClearNames(ns ...NameType) []NameType {
	nts := make([]NameType, len(ns))
	for i, nt := range ns {
		nts[i].T = nt.T
	}
	return nts
}

// PrefixWriteTo writes the name followed by a space then the prefixed type.
func (n NameType) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	if n.N != "" {
		sw.WriteString(n.N)
		sw.WriteRune(' ')
	}
	n.T.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing NameType")
	return sw.Rets()
}
