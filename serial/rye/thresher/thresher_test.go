package thresher

import (
	"fmt"
	//	"github.com/adamcolton/luce/serial/rye"
	"bytes"
	"encoding/gob"
	"strconv"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	First        string
	Last         string
	Age          int
	Role         int
	StreetNumber int
	StreetName   string
	City         string
}

func (p *Person) TypeID() uint64 {
	return 2
}

func TestRegister(t *testing.T) {
	p := &Person{
		First:        "Adam",
		Last:         "Colton",
		Age:          34,
		Role:         2,
		StreetNumber: 31415,
		StreetName:   "Pi Dr",
		City:         "Williamston",
	}

	//p = nil

	th := &Thresher{}
	th.Register((*Person)(nil))

	b, err := th.Marshal(p)
	assert.NoError(t, err)

	i, err := th.Unmarshal(b)
	assert.NoError(t, err)

	p2 := i.(*Person)
	assert.Equal(t, p, p2)
}

type Foo []string

func (*Foo) TypeID() uint64   { return 3 }
func (f *Foo) String() string { return strings.Join(*f, "|") }

func TestStringSlice(t *testing.T) {
	th := &Thresher{}
	th.Register((*Foo)(nil))

	f := &Foo{"this", "is", "a", "test"}
	b, err := th.Marshal(f)
	assert.NoError(t, err)

	i, err := th.Unmarshal(b)
	f2 := (i.(*Foo))
	assert.Equal(t, f, f2)
}

type AllTypes struct {
	Int       int
	Int8      int8
	Int16     int16
	Int32     int32
	Int64     int64
	Uint      uint
	Byte      byte
	Uint8     uint8
	Uint16    uint16
	Uint32    uint32
	Uint64    uint64
	Float32   float32
	Float64   float64
	PtrInt    *int
	Interface fmt.Stringer
}

func (*AllTypes) TypeID() uint64 {
	return 4
}

func TestAllTypes(t *testing.T) {
	th := &Thresher{}
	assert.NoError(t, th.Register((*AllTypes)(nil), (*Foo)(nil)))

	iPtr := 123
	ai := &AllTypes{
		Int:       1,
		Int8:      2,
		Int16:     3,
		Int32:     4,
		Int64:     5,
		Uint:      6,
		Byte:      7,
		Uint8:     8,
		Uint16:    9,
		Uint32:    10,
		Uint64:    11,
		Float32:   3.1415,
		Float64:   3.141592653,
		PtrInt:    &iPtr,
		Interface: &Foo{"a", "b", "c", "d"},
	}
	b, err := th.Marshal(ai)
	assert.NoError(t, err)

	i, err := th.Unmarshal(b)
	assert.NoError(t, err)

	ai2 := i.(*AllTypes)

	assert.Equal(t, ai, ai2)
}

func TestAllTypesZero(t *testing.T) {
	th := &Thresher{}
	th.Register((*AllTypes)(nil))

	ai3 := &AllTypes{}
	b, err := th.Marshal(ai3)
	assert.NoError(t, err)
	// 0: TypeID
	// 1: Ptr not null
	// 2: End
	assert.Len(t, b, 3)
}

type Bar struct {
	Foo string `RyeField:"1"`
	Bar int    `RyeField:"2"`
}

type BarSlice []Bar

func (*BarSlice) TypeID() uint64 { return 5 }

func TestSliceOfStruct(t *testing.T) {
	th := &Thresher{}
	th.Register((*BarSlice)(nil))

	bs := &BarSlice{
		Bar{"A", 1},
		Bar{"B", 2},
		Bar{"C", 3},
	}
	b, err := th.Marshal(bs)
	assert.NoError(t, err)

	i, err := th.Unmarshal(b)
	bs2 := i.(*BarSlice)

	assert.Equal(t, bs, bs2)
}

type A struct {
	A int
	B *B
	C HasType
}

func (*A) TypeID() uint64 { return 6 }

type B struct {
	B int
	A *A
	C HasType
}

func (*B) TypeID() uint64 { return 7 }

func TestCyclic(t *testing.T) {
	th := &Thresher{}
	th.Register((*A)(nil))
	th.Register((*B)(nil))

	a := &A{
		A: 5,
		B: &B{
			B: 10,
		},
		C: &B{
			B: 15,
		},
	}
	// a.B.A = a // Can this be made to work?
	b, err := th.Marshal(a)
	assert.NoError(t, err)

	i, err := th.Unmarshal(b)
	a2 := i.(*A)

	assert.Equal(t, a, a2)
	a.B.B = 20
	assert.NotEqual(t, a, a2)
}

const (
	sflag uint64 = (1 << 63) - 1
)

func TestFloat(t *testing.T) {
	t.Skip() // looking into better way to compress a float
	f := 1.5
	u := *(*uint64)(unsafe.Pointer(&f))
	fmt.Println(b(u))
	s := u >> 63
	u &= sflag
	e := u >> 52
	m := u ^ (e << 52)
	fmt.Println(s, b(e), b(m))
}

func b(u uint64) string {
	s := strconv.FormatUint(u, 2)
	for len(s) < 64 {
		s = "0" + s
	}
	return s
}

func BenchmarkRye(b *testing.B) {
	th := &Thresher{}
	th.Register((*AllTypes)(nil), (*Foo)(nil))

	iPtr := 123
	ai := &AllTypes{
		Int:       1,
		Int8:      2,
		Int16:     3,
		Int32:     4,
		Int64:     5,
		Uint:      6,
		Byte:      7,
		Uint8:     8,
		Uint16:    9,
		Uint32:    10,
		Uint64:    11,
		Float32:   3.1415,
		Float64:   3.141592653,
		PtrInt:    &iPtr,
		Interface: &Foo{"a", "b", "c", "d"},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b, _ := th.Marshal(ai)

		th.Unmarshal(b)
	}
}

func BenchmarkGob(b *testing.B) {
	iPtr := 123
	ai := &AllTypes{
		Int:       1,
		Int8:      2,
		Int16:     3,
		Int32:     4,
		Int64:     5,
		Uint:      6,
		Byte:      7,
		Uint8:     8,
		Uint16:    9,
		Uint32:    10,
		Uint64:    11,
		Float32:   3.1415,
		Float64:   3.141592653,
		PtrInt:    &iPtr,
		Interface: &Foo{"a", "b", "c", "d"},
	}

	var out *AllTypes
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(nil)
		gob.NewEncoder(buf).Encode(ai)
		gob.NewDecoder(buf).Decode(out)
		buf.Reset()
	}
}
