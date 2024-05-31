package thresher

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

func clearMemory() {
	byPtr = lmap.Map[uintptr, *rootObj]{}
	byID = lmap.Map[string, *rootObj]{}
}

type Node struct {
	Next  *Node
	Value string
}

func TestGraph(t *testing.T) {
	words := []string{
		"young", "unity", "defend", "storage", "law", "pack", "strike",
		"triangle", "agenda", "knee", "model", "resist", "hike", "aspect",
		"wander", "photography", "strain", "school", "definite", "advocate",
		"map", "projection", "warm", "research", "instinct", "parking",
		"contain", "danger", "deadly", "premature", "day", "brilliance",
		"diplomatic", "colony", "effort", "faith", "harbor", "weigh",
		"impound", "bond", "acquit", "apparatus", "tile", "heart", "wait",
	}

	first := &Node{
		Value: words[0],
	}
	cur := first

	for _, w := range words[1:] {
		cur.Next = &Node{
			Value: w,
		}
		cur = cur.Next
	}
	cur.Next = first

	firstID := Save(first)
	clearMemory()

	cur = lerr.OK(Get[*Node](firstID))(errBadDecode)

	for _, w := range words {
		assert.Equal(t, w, cur.Value)
		cur = cur.Next
	}
}

func TestProofGetPointer(t *testing.T) {
	i := 123
	v := reflect.ValueOf(i)
	v2 := reflect.New(v.Type())
	v2.Elem().Set(v)
	fmt.Println(v2.Type())

	pi := v2.Interface().(*int)
	assert.Equal(t, &i, pi)
}

func TestStructCodec(t *testing.T) {
	c := makeStructCodec(reflector.Type[Node]())

	n := Node{
		Value: "this is a test",
	}

	size := c.size(n)
	s := compact.MakeSerializer(int(size))
	c.enc(n, s)
	d := compact.NewDeserializer(s.Data)
	n2 := c.dec(d).(Node)
	assert.Equal(t, n, n2)
}

var errBadDecode = lerr.Str("bad decode")

func TestStructPtrGraph(t *testing.T) {
	expected := "this is a test"
	n := &Node{
		Value: expected,
	}

	nID := Save(n)
	n = nil
	clearMemory()

	n = lerr.OK(Get[*Node](nID))(errBadDecode)
	assert.Equal(t, expected, n.Value)
}

func TestRing(t *testing.T) {
	n1 := &Node{
		Value: "node 1",
	}
	n2 := &Node{
		Value: "node 2",
	}
	n1.Next = n2
	n2.Next = n1

	n1id := Save(n1)
	n1 = nil
	n2 = nil
	clearMemory()

	n1 = lerr.OK(Get[*Node](n1id))(errBadDecode)

	assert.Equal(t, "node 1", n1.Value)
	assert.Equal(t, "node 2", n1.Next.Value)
}

func TestRootObject(t *testing.T) {
	n := &Node{
		Value: "testing",
	}
	rootObjByV(reflect.ValueOf(n))
}

func TestIntSlice(t *testing.T) {
	s := []int{3, 1, 4, 1, 5}

	sid := Save(s)
	clearMemory()

	got := lerr.OK(Get[[]int](sid))(errBadDecode)
	assert.Equal(t, s, got)
}

type Person struct {
	Name string
	Age  int
}

func TestStructSlice(t *testing.T) {
	s := []Person{
		{
			Name: "Adam",
			Age:  39,
		},
		{
			Name: "Lauren",
			Age:  38,
		},
		{
			Name: "Fletcher",
			Age:  5,
		},
	}
	sid := Save(s)

	clearMemory()

	got := lerr.OK(Get[[]Person](sid))(errBadDecode)
	assert.Equal(t, s, got)
}

func TestPointerSlice(t *testing.T) {
	s := []*Person{
		{
			Name: "Adam",
			Age:  39,
		},
		{
			Name: "Lauren",
			Age:  38,
		},
		{
			Name: "Fletcher",
			Age:  5,
		},
	}
	sid := Save(s)
	clearMemory()

	got := lerr.OK(Get[[]*Person](sid))(errBadDecode)
	assert.Equal(t, s, got)
}

func TestProofReflectCast(t *testing.T) {
	var i64 int64 = 31415
	var i32 int32
	reflect.ValueOf(&i32).Elem().SetInt(i64)
	assert.Equal(t, int32(i64), i32)
}

func TestStructEncoding(t *testing.T) {
	c := getCodec(reflector.Type[Person]())

	e := compact.NewDeserializer(encodings[string(c.encodingID)])

	age := compact.NewDeserializer(encodings[string(e.CompactSlice())])
	assert.Equal(t, "Age", age.CompactString())
	assert.Equal(t, intEncID, age.CompactSlice())

	name := compact.NewDeserializer(encodings[string(e.CompactSlice())])
	assert.Equal(t, "Name", name.CompactString())
	assert.Equal(t, compactSliceEncID, name.CompactSlice())

}
