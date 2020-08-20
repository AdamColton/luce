package gothicgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKindString(t *testing.T) {
	assert.Equal(t, "StringKind", StringKind.String())
	assert.Equal(t, "UndefinedKind", (Kind(255)).String())
}
