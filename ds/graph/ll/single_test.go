package ll_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/ds/graph/ll"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestSinglePtr(t *testing.T) {
	expected := []graph.KV[string, int]{
		{"3", 3},
		{"1", 1},
		{"4", 4},
		{"1", 1},
		{"5", 5},
		{"9", 9},
	}
	p := graph.RawPointer[ll.Single[string, int]]{}
	n := ll.NewSingleLoop(p, expected[0].K, expected[0].V)
	s := n
	assert.Equal(t, lerr.OK(n.Next.Get())(ll.ErrPtrMiss), n)
	for _, kv := range expected[1:] {
		nxt := n.InsertAfter(kv.K, kv.V)
		assert.Equal(t, lerr.OK(n.Next.Get())(ll.ErrPtrMiss), nxt)
		n = nxt
	}

	got := slice.IterSlice(s.Iter(), nil)
	assert.Equal(t, expected, got)
}

func TestFoo(t *testing.T) {
	b := []byte{0, 0}
	rye.Serialize.Uint16(b, 59495)
	fmt.Println(b)

	var prev uint64 = 249
	for i := 1; i < 8; i++ {
		nxt := prev + exp(i)
		fmt.Println(prev, nxt)
		prev = nxt + 1
	}
}

func exp(i int) (u uint64) {
	if i == 0 {
		return 0
	}
	u = 256
	for j := 1; j < i; j++ {
		u *= 256
	}
	u--
	return u
}
