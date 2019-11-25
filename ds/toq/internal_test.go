package toq

import (
	"fmt"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestLinkedLists(t *testing.T) {
	tq := New(time.Second*5, 10)

	assert.EqualValues(t, empty, tq.head)
	assert.EqualValues(t, empty, tq.tail)
	assert.EqualValues(t, empty, tq.free)

	action := func() {}

	cs := make([]Token, 3)

	cs[0] = tq.Add(action)
	assert.EqualValues(t, 0, tq.head)
	assert.EqualValues(t, 0, tq.tail)
	assert.EqualValues(t, empty, tq.free)

	assert.True(t, cs[0].Cancel())
	assert.EqualValues(t, empty, tq.head)
	assert.EqualValues(t, empty, tq.tail)
	assert.EqualValues(t, 0, tq.free)

	cs[1] = tq.Add(action)
	assert.EqualValues(t, 0, tq.head)
	assert.EqualValues(t, 0, tq.tail)
	assert.EqualValues(t, empty, tq.free)

	// calling previous cancel again does nothing
	assert.False(t, cs[0].Cancel())
	assert.EqualValues(t, 0, tq.head)
	assert.EqualValues(t, 0, tq.tail)
	assert.EqualValues(t, empty, tq.free)

	cs[0] = tq.Add(action)
	assert.EqualValues(t, 0, tq.head)
	assert.EqualValues(t, 1, tq.tail)
	assert.EqualValues(t, empty, tq.free)

	cs[2] = tq.Add(action)
	assert.EqualValues(t, 0, tq.head)
	assert.EqualValues(t, 2, tq.tail)
	assert.EqualValues(t, empty, tq.free)

	assert.True(t, cs[2].Cancel())
	assert.EqualValues(t, 0, tq.head)
	assert.EqualValues(t, 1, tq.tail)
	assert.EqualValues(t, 2, tq.free)

	assert.True(t, cs[1].Cancel())
	assert.EqualValues(t, 1, tq.head)
	assert.EqualValues(t, 1, tq.tail)
	assert.EqualValues(t, 0, tq.free)

	assert.True(t, cs[0].Cancel())
	assert.EqualValues(t, empty, tq.head)
	assert.EqualValues(t, empty, tq.tail)
	assert.EqualValues(t, 1, tq.free)

	// Just to get to 100% test coverage
	Token(token{}).private()
}

func getAction(ch chan<- int, i int) func() {
	return func() {
		ch <- i
		fmt.Println(i)
	}
}

func TestReset(t *testing.T) {
	defer func() {
		now = time.Now
		sleep = time.Sleep
	}()
	mt := timeout.NewMockTime()
	now = mt.Now
	sleep = mt.Sleep

	var multiTick = func(ticks int) {
		for ; ticks > 0; ticks-- {
			time.Sleep(time.Millisecond)
			mt.Tick(0)
		}
	}

	tq := New(time.Millisecond*20, 2)
	ch := make(chan int)

	tokens := make([]Token, 3)
	tokens[0] = tq.Add(getAction(ch, 1))
	mt.Tick(1)
	tokens[1] = tq.Add(getAction(ch, 2))
	mt.Tick(1)
	tokens[2] = tq.Add(getAction(ch, 3))
	mt.Tick(1)

	assert.True(t, tokens[1].Reset())
	assert.NoError(t, timeout.After(100, func() {
		mt.Tick(30)
		fmt.Println("A")
		go multiTick(50)
		fmt.Println("B")
		got := map[int]bool{
			<-ch: true,
			<-ch: true,
			<-ch: true,
		}
		fmt.Println("C")
		assert.True(t, got[1])
		assert.True(t, got[2])
		assert.True(t, got[3])
		fmt.Println("D")
	}))

	// mt.Tick(1)
	// tokens[0] = tq.Add(getAction(ch, 4))
	// mt.Tick(1)
	// tokens[1] = tq.Add(getAction(ch, 5))
	// mt.Tick(1)
	// tokens[2] = tq.Add(getAction(ch, 6))
	// mt.Tick(1)
	// assert.True(t, tokens[1].Reset())
	// mt.Tick(1)
	// assert.True(t, tokens[0].Reset())
	// mt.Tick(1)
	// assert.True(t, tokens[1].Cancel())
	// assert.False(t, tokens[1].Reset())
	// mt.Tick(1)
	// tq.Add(getAction(ch, 7))
	// assert.NoError(t, timeout.After(100, func() {
	// 	mt.Tick(21)
	// 	go multiTick(5)
	// 	got := map[int]bool{
	// 		<-ch: true,
	// 		<-ch: true,
	// 		<-ch: true,
	// 	}
	// 	assert.True(t, got[6])
	// 	assert.True(t, got[4])
	// 	assert.True(t, got[7])
	// }))
}
