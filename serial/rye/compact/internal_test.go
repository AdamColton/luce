package compact

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Each tier starts one past the end of the one before, and the first starts
// one past the longest short length.
func TestTierStart(t *testing.T) {
	assert.Equal(t, uint64(sliceMeta-1+shortMin+1), tierStart[0])
	for i := 1; i < len(tierStart); i++ {
		assert.Equal(t, tierStart[i-1]+1<<(8*i), tierStart[i], "tier %d", i)
	}
	assert.Equal(t, sliceNil-sliceMeta, len(tierStart))
}

func TestSliceLen(t *testing.T) {
	lns := []uint64{shortMin, tierStart[0] - 1}
	for _, start := range tierStart {
		lns = append(lns, start, start+1)
	}
	last := tierStart[len(tierStart)-1]
	lns = append(lns, last+1<<56-1)

	for _, ln := range lns {
		t.Run(fmt.Sprint(ln), func(t *testing.T) {
			s := NewSerializer(9)
			s.Make()
			s.sliceLen(ln)
			if ln < tierStart[0] {
				assert.Equal(t, 1, s.Idx)
			} else {
				assert.Equal(t, 2+int(tier(ln)), s.Idx)
			}

			d := NewDeserializer(s.Data)
			assert.Equal(t, ln, d.sliceLen(d.Byte()))
			assert.Equal(t, s.Idx, d.Idx)
		})
	}
}
