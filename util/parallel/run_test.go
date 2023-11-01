package parallel_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/util/parallel"
	"github.com/stretchr/testify/assert"
)

type testRecord struct {
	core, idx int
}

func TestRunRange(t *testing.T) {
	out := make([]*testRecord, 1000)
	wg := parallel.RunRange(len(out), func(rangeIdx, coreIdx int) {
		assert.Nil(t, out[rangeIdx])
		out[rangeIdx] = &testRecord{
			core: coreIdx,
			idx:  rangeIdx,
		}
		time.Sleep(time.Millisecond)
	})
	wg.Wait()

	for i, r := range out {
		assert.Equal(t, r.idx, i)
	}
}
