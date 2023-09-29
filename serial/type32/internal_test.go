package type32

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceToUint32(t *testing.T) {
	assert.Equal(t, uint32(0), sliceToUint32(nil))
}
