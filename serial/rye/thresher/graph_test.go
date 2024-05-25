package thresher

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

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

	c := getCodec(reflect.TypeOf(first))
	size := c.size(first)
	s := compact.MakeSerializer(int(size))
	c.enc(first, s)
	d := compact.NewDeserializer(s.Data)
	n2 := c.dec(d).(*Node)
	assert.Equal(t, first, n2)

	g := Graph(first)
	assert.Equal(t, 1, g.types.Len())
	assert.Equal(t, len(words), g.ptrs.Len())

	data := g.enc()
	assert.NotNil(t, data)
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

func TestRootObject(t *testing.T) {
	n := &Node{
		Value: "testing",
	}
	rootObjByV(reflect.ValueOf(n))
}
