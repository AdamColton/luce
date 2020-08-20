package gothic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testGen struct {
	prepared  bool
	generated bool
	err       error
}

func (t *testGen) Prepare() error {
	t.prepared = true
	return t.err
}

func (t *testGen) Generate() error {
	t.generated = true
	return t.err
}

func TestHappyPath(t *testing.T) {
	tg := &testGen{}
	p := New()

	err := p.AddGenerators(tg)
	assert.NoError(t, err)
	assert.False(t, tg.prepared, "tg.prepared should be false")
	assert.False(t, tg.generated, "tg.prepared should be false")

	err = p.Generate()
	assert.Error(t, err)

	err = p.Prepare()
	assert.NoError(t, err)
	assert.True(t, tg.prepared, "tg.prepared should be true")
	assert.False(t, tg.generated, "tg.prepared should be false")

	err = p.Prepare()
	assert.Equal(t, ErrPrepareCalled, err)

	err = p.AddGenerators(tg)
	assert.Error(t, err)

	err = p.Generate()
	assert.NoError(t, err)
	assert.True(t, tg.prepared, "tg.prepared should be true")
	assert.True(t, tg.generated, "tg.prepared should be false")

	err = p.Prepare()
	assert.Error(t, err)

	err = p.Generate()
	assert.Error(t, err)
}

func TestHappyPathExport(t *testing.T) {
	tg := &testGen{}
	p := New()

	err := p.AddGenerators(tg)
	assert.NoError(t, err)
	assert.False(t, tg.prepared, "tg.prepared should be false")
	assert.False(t, tg.generated, "tg.prepared should be false")

	err = p.Export()
	assert.NoError(t, err)

	assert.True(t, tg.prepared, "tg.prepared should be true")
	assert.True(t, tg.generated, "tg.prepared should be false")
}

func TestPrepareError(t *testing.T) {
	tg := &testGen{
		err: fmt.Errorf("test error"),
	}
	p := New()

	err := p.AddGenerators(tg)
	assert.NoError(t, err)

	err = p.Prepare()
	assert.Equal(t, "test error", err.Error())

	err = p.Prepare()
	assert.Equal(t, "Project is in an error state.", err.Error())

	err = p.Export()
	assert.Equal(t, "Project is in an error state.", err.Error())
}

func TestGenerateError(t *testing.T) {
	tg := &testGen{}
	p := New()

	err := p.AddGenerators(tg)
	assert.NoError(t, err)

	err = p.Prepare()
	assert.NoError(t, err)

	tg.err = fmt.Errorf("test error")

	err = p.Generate()
	assert.Equal(t, "test error", err.Error())

	err = p.Generate()
	assert.Equal(t, "Project is in an error state.", err.Error())
}
