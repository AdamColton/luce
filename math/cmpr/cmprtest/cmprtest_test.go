package cmprtest_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/math/cmpr/cmprtest"
	"github.com/stretchr/testify/assert"
)

func TestCmprTest(t *testing.T) {
	d := float64(cmprtest.Small / 10)
	a := 3.1415
	b := a + d
	cmprtest.Equal(t, a, b)
}

type mock struct {
	msgs []string
}

func (m *mock) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	m.msgs = append(m.msgs, msg)
}

func (m *mock) clear() {
	m.msgs = m.msgs[:0]
}

func TestWithMock(t *testing.T) {
	m := &mock{}

	cmprtest.Equal(m, 3, 4)
	assert.Len(t, m.msgs, 1)
	assert.Equal(t, "Expected 3 got 4", m.msgs[0])

	m.clear()
	cmprtest.Equal(m, 3, 4, "MESSAGE")
	assert.Len(t, m.msgs, 1)
	assert.Equal(t, "Expected 3 got 4: MESSAGE", m.msgs[0])
}

// TODO: Add helpers
// AssertEqualValue() float64
// cmpr.AssertEqual[T any](to T, t cmpr.Tolerance) error

type equalizer struct {
	PartName string
	Length   float64
}

func (e *equalizer) AssertEqual(to interface{}, t cmpr.Tolerance) error {
	e2, ok := to.(*equalizer)
	if !ok {
		return lerr.Str("equalizer.AssertEqual requires equalizer for argument 'to'")
	}
	if !t.Equal(e.Length, e2.Length) {
		return lerr.Str("values not equal")
	}
	return nil
}

func TestAssertEqual(t *testing.T) {
	a := &equalizer{
		PartName: "widget",
		Length:   1.23,
	}
	b := &equalizer{
		PartName: "gadget",
		Length:   1.23,
	}

	cmprtest.Equal(t, a, b)

	m := &mock{}
	b.Length = 3.14
	cmprtest.Equal(m, a, b)
	assert.Len(t, m.msgs, 1)
	assert.Equal(t, "values not equal", m.msgs[0])
}

type sliceAlias []float64

func TestSlice(t *testing.T) {
	a := sliceAlias{1, 2, 3, 4, 5}
	b := sliceAlias{1, 2, 3, 4, 5}
	cmprtest.Equal(t, a, b)

	m := &mock{}
	c := []string{"1", "2", "3", "4", "5"}
	cmprtest.Equal(m, a, c)
	assert.Len(t, m.msgs, 1)
	expected := make([]string, 5)
	for i := range expected {
		expected[i] = fmt.Sprintf(`%d: Types do not match: expected "float64", got "string"`, i)
	}
	assert.Equal(t, "\t"+strings.Join(expected, "\n\t"), m.msgs[0])

	m.clear()
	cmprtest.Equal(m, a, "test")
	assert.Len(t, m.msgs, 1)
	assert.Equal(t, `Types do not match: expected "cmprtest_test.sliceAlias", got "string"`, m.msgs[0])

	m.clear()
	cmprtest.Equal(m, equalizer{}, "test")
	assert.Len(t, m.msgs, 1)
	assert.Equal(t, `unsupported_type: cmprtest_test.equalizer (*cmprtest_test.equalizer fulfills AssertEqualizer)`, m.msgs[0])
}

func TestTypes(t *testing.T) {
	cmprtest.Equal(t, float64(10), int(10))
	cmprtest.Equal(t, float64(10), int8(10))
	cmprtest.Equal(t, float64(10), int16(10))
	cmprtest.Equal(t, float64(10), int32(10))
	cmprtest.Equal(t, float64(10), int64(10))
	cmprtest.Equal(t, float64(10), uint(10))
	cmprtest.Equal(t, float64(10), uint8(10))
	cmprtest.Equal(t, float64(10), uint16(10))
	cmprtest.Equal(t, float64(10), uint32(10))
	cmprtest.Equal(t, float64(10), uint64(10))
	cmprtest.Equal(t, float64(10), float32(10))
}

func TestMessage(t *testing.T) {
	assert.Equal(t, "", cmprtest.Message())
	assert.Equal(t, "123", cmprtest.Message(123))
	str := "%d"
	assert.Equal(t, "123", cmprtest.Message(str, 123))
	assert.Equal(t, "123 456", cmprtest.Message(123, 456))
}
