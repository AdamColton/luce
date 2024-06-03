package entity

import (
	"bytes"
	"encoding/base64"
	"reflect"

	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

type Getter[E Entity] interface {
	Get(key []byte) (found bool, ent E, err error)
}

var getters = map[reflect.Type]any{}

func AddGetter[E Entity, G Getter[E]](getter G) {
	t := reflector.Type[E]()
	getters[t] = getter
}

func GetGetter[E Entity]() Getter[E] {
	t := reflector.Type[E]()
	g, found := getters[t]
	if found {
		return g.(Getter[E])
	}
	return nil
}

type Reference[E Entity] struct {
	ID  []byte
	Ent E    `json:"-"`
	set bool `json:"-"`
}

func NewRef[E Entity](e E) Reference[E] {
	id := e.EntKey()
	return Reference[E]{
		ID:  id,
		Ent: e,
		set: len(id) > 0,
	}
}

func (ref *Reference[E]) Get() (E, bool) {
	if !ref.set {
		g := GetGetter[E]()
		if g != nil {
			found, e, err := g.Get(ref.ID)
			lerr.Panic(err)
			if found {
				ref.set = true
				ref.Ent = e
			}
		}
	}

	return ref.Ent, ref.set
}

func (ref *Reference[E]) GobDecode(id []byte) error {
	ref.ID = id
	return nil
}

func (ref *Reference[E]) GobEncode() ([]byte, error) {
	return ref.ID, nil
}

const dblQuote byte = 34

var b64enc = base64.RawStdEncoding

func (ref *Reference[E]) MarshalJSON() ([]byte, error) {
	ln := b64enc.EncodedLen(len(ref.ID))
	out := make([]byte, ln+2)
	out[0], out[ln+1] = dblQuote, dblQuote
	b64enc.Encode(out[1:], ref.ID)

	return out, nil
}

func (ref *Reference[E]) UnmarshalJSON(str []byte) (err error) {
	str = str[1 : len(str)-1]
	ln := b64enc.DecodedLen(len(str))
	s := string(str)
	_ = s
	ref.ID = make([]byte, ln)
	_, err = b64enc.Decode(ref.ID, str)
	return err
}

func (ref *Reference[E]) Set(e E) graph.Ptr[E] {
	id := e.EntKey()
	if bytes.Equal(ref.ID, id) {
		ref.Ent = e
		ref.set = true
		return ref
	}
	return &Reference[E]{
		ID:  id,
		Ent: e,
		set: true,
	}
}

func (ref *Reference[E]) New() graph.Ptr[E] {
	return &Reference[E]{}
}
