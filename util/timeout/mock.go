package timeout

import (
	"sync"
	"sync/atomic"
	"time"
)

// Staging:
// sleeps all wait for wg
// sleeps all check

type MockTime struct {
	N            time.Time
	TickDuration time.Duration

	mux    sync.Mutex
	sleeps atomic.Int32
	tick   sync.WaitGroup
	tock   sync.WaitGroup
	sleep  sync.WaitGroup
}

func NewMockTime() *MockTime {
	mt := &MockTime{
		N:            time.Now(),
		TickDuration: time.Millisecond,
	}
	mt.tick.Add(1)
	return mt
}

func (mt *MockTime) Now() time.Time {
	return mt.N
}

func (mt *MockTime) Tick(ticks int) {
	mt.N = mt.N.Add(time.Duration(ticks) * mt.TickDuration)

	mt.mux.Lock()
	s := int(mt.sleeps.Load())
	mt.tock.Add(1)
	mt.sleep.Add(s)
	mt.tick.Done() // unlock all sleep threads
	mt.sleep.Wait()
	mt.tick.Add(1)
	mt.sleep.Add(s)
	mt.tock.Done()
	mt.sleep.Wait()
	mt.mux.Unlock()

}

func (mt *MockTime) Sleep(d time.Duration) {
	end := mt.N.Add(d)
	done := false
	mt.mux.Lock()
	mt.sleeps.Add(1)
	mt.mux.Unlock()
	for !done {
		mt.tick.Wait()
		done = end.After(mt.N)
		if done {
			mt.sleeps.Add(-1)
		}
		mt.sleep.Done()

		mt.tock.Wait()
		mt.sleep.Done()
	}
}
