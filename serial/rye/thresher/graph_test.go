package thresher

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
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

	firstID := rootObjByV(reflect.ValueOf(first)).id

	c := getCodec(reflect.TypeOf(first))
	size := c.size(first)
	s := compact.MakeSerializer(int(size))
	c.enc(first, s)
	d := compact.NewDeserializer(s.Data)
	n2 := (<-c.dec(d)).(*Node)
	assert.Equal(t, first, n2)

	g := Graph(first)
	assert.Equal(t, 1, g.types.Len())
	assert.Equal(t, len(words), g.ptrs.Len())

	data := g.enc()
	assert.NotNil(t, data)

	clearMemory()

	g2 := &grapher{}
	g2.dec(data)
	assert.Equal(t, len(words), g.ptrs.Len())

	cur = rootObjByID(firstID).v.Interface().(*Node)

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
	n2 := (<-c.dec(d)).(Node)
	assert.Equal(t, n, n2)
}

func TestStructPtrGraph(t *testing.T) {
	expected := "this is a test"
	n := &Node{
		Value: expected,
	}

	g := Graph(n)
	data := g.enc()
	n = nil
	clearMemory()
	if len(data) > 100 {
		panic("wtf")
	}

	g = newGrapher()
	g.dec(data)
	for _, ro := range g.ptrs {
		n = ro.v.Interface().(*Node)
	}

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

	g := Graph(n1)
	data := g.enc()
	n1id := rootObjByV(reflect.ValueOf(n1)).id
	n1 = nil
	n2 = nil
	clearMemory()

	g = newGrapher()
	g.dec(data)
	n1 = rootObjByID(n1id).v.Interface().(*Node)
	assert.Equal(t, "node 1", n1.Value)
}

func TestRootObject(t *testing.T) {
	n := &Node{
		Value: "testing",
	}
	rootObjByV(reflect.ValueOf(n))
}
