package rbtree_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/ds/graph/rbtree"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/lrand"
	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	ptr := rbtree.MakePtrType[int, string]()
	rb := rbtree.New(ptr, filter.Comparer[int]())
	r := lrand.New()
	rng := 100

	randKV := func() graph.KV[int, string] {
		k := r.Intn(rng * 100)
		return graph.NewKV(k, strconv.Itoa(k))
	}

	vals := lset.New[graph.KV[int, string]]()
	var max int
	for i := 0; i < rng; i++ {
		kv := randKV()
		if vals.Contains(kv) {
			i--
			continue
		}
		fmt.Println(kv)
		rb.Add(kv.K, kv.V)
		vals.Add(kv)
		if !assert.True(t, rb.Validate(), i) {
			rb.Print()
		}
		max = cmpr.Max(max, kv.K)
	}

	vals.Each(func(kv graph.KV[int, string]) (done bool) {
		n, found := rb.Seek(kv.K)
		assert.True(t, found)
		assert.Equal(t, kv.V, n.Val())
		return false
	})

	for i := 0; i < 100; i++ {
		kv := randKV()
		if i == 0 {
			kv.K = max + 1
			kv.V = strconv.Itoa(kv.K)
		}
		if vals.Contains(kv) {
			i--
			continue
		}
		n, found := rb.Seek(kv.K)
		assert.False(t, found)
		if kv.K > max {
			assert.Nil(t, n)
		} else {
			assert.Greater(t, n.Key(), kv.K)
		}
	}

	fmt.Println("---")
	s := vals.Slice()
	for _, skv := range s {
		vals.Each(func(kv graph.KV[int, string]) (done bool) {
			n, found := rb.Seek(kv.K)
			assert.True(t, found, kv.K)
			assert.Equal(t, kv.K, n.Key())
			assert.Equal(t, kv.V, n.Val())
			return false
		})

		fmt.Println(skv.K)
		rb.Remove(skv.K)
		vals.Remove(skv)
		assert.True(t, rb.Validate())
	}
}

// func TestFoo(t *testing.T) {
// 	ptr := rbtree.MakePtrType[int]()
// 	rb := rbtree.New(ptr, filter.Comparer[int]())
// 	add := []int{737, 244, 263, 302, 559, 581, 508, 246, 780, 403}
// 	rm := []int{559, 581, 508, 246, 780, 403, 244, 263, 737, 302}
// 	for _, a := range add {
// 		rb.Add(a)
// 	}
// 	for _, r := range rm {
// 		rb.Remove(r)
// 	}

// 	rb.Remove(41)
// }
