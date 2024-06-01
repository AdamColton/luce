package thresher

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"reflect"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye/compact"
)

type fieldCodec struct {
	idx  int
	name string
	*codec
}

func (fc *fieldCodec) hash(t reflect.Type, h hash.Hash) []byte {
	f := t.Field(fc.idx)
	d := getCodec(f.Type)
	n := f.Name

	size := compact.SizeString(n) + compact.Size(d.encodingID)
	s := compact.MakeSerializer(int(size))
	s.CompactString(n)
	s.CompactSlice(d.encodingID)

	h.Write([]byte(f.Name))
	h.Write(d.encodingID)
	id := h.Sum(nil)

	store[string(id)] = s.Data

	return id
}

type structCodec struct {
	fieldCodecs slice.Slice[fieldCodec]
	reflect.Type
	encodingID []byte
}

func (sc *structCodec) enc(i any, s compact.Serializer) {
	v := reflect.ValueOf(i)
	s.CompactSlice(sc.encodingID)
	for _, fc := range sc.fieldCodecs {
		f := v.Field(fc.idx).Interface()
		fc.enc(f, s)
	}
}

func (sc *structCodec) dec(d compact.Deserializer) any {
	srct := reflect.New(sc.Type).Elem()
	id := d.CompactSlice()
	if !bytes.Equal(id, sc.encodingID) {
		panic("encodingIDs not equal")
	}
	for _, fc := range sc.fieldCodecs {
		idx := fc.idx
		str := sc.Type.Field(idx).Type.String()
		_ = str
		dec := decoders[typeEncoding{
			encID: string(fc.encodingID),
			t:     sc.Type.Field(idx).Type,
		}]
		i := dec(d)
		if i != nil {
			fv := reflect.ValueOf(i)
			srct.Field(idx).Set(fv)
		}
	}
	return srct.Interface()
}

func (sc *structCodec) size(i any) uint64 {
	v := reflect.ValueOf(i)
	sum := compact.Size(sc.encodingID)
	for _, fc := range sc.fieldCodecs {
		f := v.Field(fc.idx).Interface()
		sum += fc.size(f)
	}
	return sum
}

func (sc *structCodec) roots(v reflect.Value) (out []*rootObj) {
	for _, fc := range sc.fieldCodecs {
		if fc.roots != nil {
			f := v.Field(fc.idx)
			out = append(out, fc.roots(f)...)
		}
	}
	return
}

func makeStructCodec(t reflect.Type) *codec {
	c := &codec{}
	codecs[t] = c
	sc := &structCodec{
		Type:        t,
		fieldCodecs: fieldsRecur(0, t.NumField(), t, 0),
	}
	sc.fieldCodecs.Sort(func(i, j fieldCodec) bool {
		return i.name < j.name
	})

	c.enc = sc.enc
	c.size = sc.size
	c.roots = sc.roots

	baseHash := sha256.New()
	fieldHash := sha256.New()

	s := fieldsHashRecur(0, 0, sc.Type, sc.fieldCodecs, baseHash, fieldHash, make([][]byte, len(sc.fieldCodecs)))
	sc.encodingID = baseHash.Sum(nil)
	c.encodingID = sc.encodingID
	store[string(sc.encodingID)] = s.Data

	decoders[typeEncoding{
		encID: string(sc.encodingID),
		t:     t,
	}] = sc.dec

	return c
}

func fieldsHashRecur(i int, size uint64, t reflect.Type, fcs slice.Slice[fieldCodec], baseHash, fieldHash hash.Hash, hashes [][]byte) compact.Serializer {
	if i == len(fcs) {
		return compact.MakeSerializer(int(size))
	}
	fc := fcs[i]
	fieldHash.Reset()
	hashes[i] = fc.hash(t, fieldHash)
	baseHash.Write(hashes[i])
	size += compact.Size(hashes[i])

	s := fieldsHashRecur(i+1, size, t, fcs, baseHash, fieldHash, hashes)
	if i == 0 {
		for _, h := range hashes {
			s.CompactSlice(h)
		}
	}

	return s
}

func fieldsRecur(i int, ln int, t reflect.Type, fields int) slice.Slice[fieldCodec] {
	for ; i < ln; i++ {
		f := t.Field(i)
		if f.IsExported() {
			if c := getCodec(f.Type); c != nil {
				fcs := fieldsRecur(i+1, ln, t, fields+1)
				fcs[fields].idx = i
				fcs[fields].codec = c
				fcs[fields].name = f.Name
				return fcs
			}
		}
	}
	return make(slice.Slice[fieldCodec], fields)
}
