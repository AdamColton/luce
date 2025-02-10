package service_test

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	r := &service.Request{}
	assert.Equal(t, service.RequestTypeID32, r.TypeID32())
}
