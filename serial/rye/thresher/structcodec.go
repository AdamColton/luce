package thresher

import (
	"crypto/sha256"
	"hash"
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye/compact"
)

type fieldCodec struct {
	idx      int
	name     string
	*encoder //TODO replace this w/ encodingID
}

func newFieldDecoder(d *compact.Deserializer) fieldCodec {
	return fieldCodec{
		name: d.CompactString(),
		idx:  -1,
		encoder: &encoder{
			encodingID: d.CompactSlice(),
		},
	}
}

func (fc *fieldCodec) hash(t reflect.Type, h hash.Hash) []byte {
	f := t.Field(fc.idx)
	d := getEncoder(f.Type)
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
	encodingID        []byte
	includeEncodingID bool
}

func (sc *structCodec) enc(i any, s compact.Serializer, base bool) {
	if base {
		s.CompactSlice(sc.encodingID)
	}
	v := reflect.ValueOf(i)
	for _, fc := range sc.fieldCodecs {
		f := v.Field(fc.idx).Interface()
		fc.encode(f, s, false)
	}
}

func (sc *structCodec) dec(d compact.Deserializer) any {
	srct := reflect.New(sc.Type).Elem()
	for _, fc := range sc.fieldCodecs {
		idx := fc.idx
		dec := getDecoder(sc.Type.Field(idx).Type, fc.encodingID)
		i := dec(d)
		if i != nil && idx >= 0 {
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

var structCodecs = lmap.Map[reflect.Type, *structCodec]{}

func getStructEncoder(t reflect.Type) *encoder {
	c := &encoder{}
	encoders[t] = c
	sc := &structCodec{
		Type:              t,
		fieldCodecs:       fieldsRecur(0, t.NumField(), t, 0),
		includeEncodingID: true,
	}
	sc.fieldCodecs.Sort(func(i, j fieldCodec) bool {
		return i.name < j.name
	})

	baseHash := sha256.New()
	fieldHash := sha256.New()

	s := fieldsHashRecur(0, 0, sc.Type, sc.fieldCodecs, baseHash, fieldHash, make([][]byte, len(sc.fieldCodecs)))
	sc.encodingID = baseHash.Sum(nil)
	store[string(sc.encodingID)] = s.Data

	addDecoder(t, sc.encodingID, sc.dec)
	structCodecs[t] = sc

	c.encodingID = sc.encodingID
	c.encode = sc.enc
	c.size = sc.size
	c.roots = sc.roots
	return c
}

func fieldsHashRecur(i int, size uint64, t reflect.Type, fcs slice.Slice[fieldCodec], baseHash, fieldHash hash.Hash, hashes [][]byte) compact.Serializer {
	if i == len(fcs) {
		i64 := uint64(i)
		size += compact.SizeUint64(i64)
		s := compact.MakeSerializer(int(size))
		s.CompactUint64(i64)
		return s
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
			if c := getEncoder(f.Type); c != nil {
				fcs := fieldsRecur(i+1, ln, t, fields+1)
				fcs[fields].idx = i
				fcs[fields].encoder = c
				fcs[fields].name = f.Name
				return fcs
			}
		}
	}
	return make(slice.Slice[fieldCodec], fields)
}

func makeStructDecoder(t reflect.Type, id []byte) decoder {
	d := compact.NewDeserializer(store[string(id)])
	fcs := make([]fieldCodec, 0, d.CompactInt64())
	for !d.Done() {
		// TODO: compact.(Serializer/Deserializer) should NOT be pointers
		// because they are wrappers
		sub := compact.NewDeserializer(store[string(d.CompactSlice())])
		fc := newFieldDecoder(&sub)
		f, _ := t.FieldByName(fc.name)
		fc.idx = f.Index[0]
		fcs = append(fcs, fc)
	}

	sc := &structCodec{
		Type:        t,
		fieldCodecs: fcs,
	}
	return sc.dec
}
