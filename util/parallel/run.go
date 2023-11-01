package parallel

import (
	"runtime"
	"sync"
	"sync/atomic"
)

var cpus = runtime.NumCPU()

// Run takes a worker that will operate on some workload until complete and
// runs it on each core. The worker is called in its own Go routine for every
// CPU and is passed the index of its core. Run returns at once with a WaitGroup
// that is done when every worker has returned.
func Run(worker func(coreIdx int)) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(cpus)
	wrap := func(coreIdx int) {
		worker(coreIdx)
		wg.Done()
	}
	for i := range cpus {
		go wrap(i)
	}
	return wg
}

// RunRange will call worker in parallel with every int from 0 to max (0
// inclusive, max exclusive). The cores share the range, so worker is called
// once for each value, with the index of the core that is calling it. It
// returns a WaitGroup that is done when every value has been processed.
func RunRange(max int, worker func(rangeIdx, coreIdx int)) *sync.WaitGroup {
	idx32 := atomic.Int32{}
	idx32.Store(-1)
	nxt := func() int { return int(idx32.Add(1)) }
	return Run(func(coreIdx int) {
		for idx := nxt(); idx < max; idx = nxt() {
			worker(idx, coreIdx)
		}
	})
}
