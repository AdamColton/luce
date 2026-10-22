package depgraph_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/depgraph"
	"github.com/stretchr/testify/assert"
)

type Top struct {
	Foo string
}

var topDef = depgraph.Define("top", func(in *depgraph.Context) *Top {
	ma := midADef.Resolve(in)
	mb := midBDef.Resolve(in)
	return &Top{
		Foo: ma.A + " : " + mb.B,
	}
})

type MidA struct {
	A string
}

var midADef = depgraph.Define("mid-A", func(in *depgraph.Context) *MidA {
	l := depgraph.ResolveType[*Leaf](in)
	return &MidA{
		A: strings.TrimSpace(l.A),
	}
})

type MidB struct {
	B string
}

var midBDef = depgraph.Define("mid-B", func(in *depgraph.Context) *MidB {
	l := depgraph.ResolveType[*Leaf](in)
	return &MidB{
		B: strings.TrimSpace(l.B),
	}
})

type Leaf struct {
	A, B string
}

func TestDepgraph(t *testing.T) {
	ctx := depgraph.NewContext()
	l := &Leaf{
		A: " this is ",
		B: "  a test  ",
	}
	depgraph.DesignateType("leaf", l, ctx)
	top := topDef.Resolve(ctx)
	assert.Equal(t, "this is : a test", top.Foo)

}

func TestValue(t *testing.T) {
	ctx := depgraph.NewContext()
	thickness := depgraph.Designate("base.thickness", 0.75, ctx)
	assert.Equal(t, 0.75, thickness.Get())
	assert.Equal(t, 0.75, thickness.Value(ctx))
	assert.Equal(t, ctx, thickness.Ctx())
	assert.Equal(t, depgraph.Node(thickness), ctx.Resolve("base.thickness"))

	// A Value only belongs to the Context it was designated in.
	other := depgraph.NewContext()
	assert.Nil(t, thickness.Value(other))
	assert.Equal(t, 0, thickness.Uses(other).Len())
	assert.Equal(t, 0, thickness.UsedBy(other).Len())
}

func TestUpdate(t *testing.T) {
	ctx := depgraph.NewContext()
	thickness := depgraph.Designate("base.thickness", 0.75, ctx)
	doubled := depgraph.Define("base.doubled", func(in *depgraph.Context) float64 {
		return 2 * thickness.Get()
	})
	assert.Equal(t, 1.5, doubled.Resolve(ctx))

	// Updating a Value recomputes what depends on it.
	thickness.Update(0.5)
	assert.Equal(t, 0.5, thickness.Get())
	assert.Equal(t, 1.0, doubled.Resolve(ctx))
}

type Board struct {
	Thickness float64
}

var boardDoubled = depgraph.Define("board.doubled", func(in *depgraph.Context) float64 {
	return 2 * depgraph.ResolveType[*Board](in).Thickness
})

func TestContexts(t *testing.T) {
	// One Definition serves every Context; each Context has its own Values.
	pine := depgraph.NewContext()
	pineBoard := depgraph.DesignateType("board", &Board{Thickness: 0.75}, pine)
	oak := depgraph.NewContext()
	depgraph.DesignateType("board", &Board{Thickness: 1}, oak)

	assert.Equal(t, 1.5, boardDoubled.Resolve(pine))
	assert.Equal(t, 2.0, boardDoubled.Resolve(oak))

	// Mutate, then Update. Only that Context recomputes.
	pineBoard.Get().Thickness = 0.5
	pineBoard.Update(pineBoard.Get())
	assert.Equal(t, 1.0, boardDoubled.Resolve(pine))
	assert.Equal(t, 2.0, boardDoubled.Resolve(oak))

	var b Board
	depgraph.Set(&b, oak)
	assert.Equal(t, 1.0, b.Thickness)
}

func TestDefineType(t *testing.T) {
	ctx := depgraph.NewContext()
	depgraph.DesignateType("board", &Board{Thickness: 0.75}, ctx)
	depgraph.DefineType("board.doubled", func(in *depgraph.Context) float64 {
		return 2 * depgraph.ResolveType[*Board](in).Thickness
	}, ctx)
	assert.Equal(t, 1.5, depgraph.ResolveType[float64](ctx))

	// Adding another node of the same type replaces the first.
	depgraph.DesignateType("board.thick", &Board{Thickness: 1}, ctx)
	assert.Equal(t, 1.0, depgraph.ResolveType[*Board](ctx).Thickness)

	assert.PanicsWithValue(t, depgraph.ErrNotFound, func() {
		depgraph.ResolveType[string](ctx)
	})
}

var chicken, egg *depgraph.Definition[string]

func init() {
	chicken = depgraph.Define("chicken", func(in *depgraph.Context) string {
		return egg.Resolve(in)
	})
	egg = depgraph.Define("egg", func(in *depgraph.Context) string {
		return chicken.Resolve(in)
	})
}

func TestCycle(t *testing.T) {
	assert.PanicsWithValue(t, depgraph.ErrCycle, func() {
		chicken.Resolve(depgraph.NewContext())
	})
}
func TestDynamicDependencies(t *testing.T) {
	ctx := depgraph.NewContext()
	useLid := depgraph.Designate("box.useLid", false, ctx)
	lid := depgraph.Designate("box.lid", 0.25, ctx)
	walls := depgraph.Designate("box.walls", 4.0, ctx)
	// height only depends on lid when useLid is set.
	height := depgraph.Define("box.height", func(in *depgraph.Context) float64 {
		h := walls.Get()
		if useLid.Get() {
			h += lid.Get()
		}
		return h
	})
	assert.Equal(t, 4.0, height.Resolve(ctx))
	assert.Equal(t, 2, height.Uses(ctx).Len())
	assert.Equal(t, 0, lid.UsedBy(ctx).Len())

	useLid.Update(true)
	assert.Equal(t, 4.25, height.Resolve(ctx))
	assert.Equal(t, 3, height.Uses(ctx).Len())
	assert.Equal(t, 1, lid.UsedBy(ctx).Len())

	// Once height stops using lid, the edge is gone and updating lid doesn't
	// reach height.
	useLid.Update(false)
	assert.Equal(t, 4.0, height.Resolve(ctx))
	assert.Equal(t, 2, height.Uses(ctx).Len())
	assert.Equal(t, 0, lid.UsedBy(ctx).Len())
	assert.False(t, lid.AllUsedBy().Contains(height))
}

func TestNoDuplicateEdges(t *testing.T) {
	ctx := depgraph.NewContext()
	side := depgraph.Designate("square.side", 2.0, ctx)
	area := depgraph.Define("square.area", func(in *depgraph.Context) float64 {
		return side.Get() * side.Get()
	})
	assert.Equal(t, 4.0, area.Resolve(ctx))
	side.Update(3)
	assert.Equal(t, 9.0, area.Resolve(ctx))
	assert.Equal(t, 1, area.Uses(ctx).Len())
	assert.Equal(t, 1, side.UsedBy(ctx).Len())
}

func TestResolve(t *testing.T) {
	ctx := depgraph.NewContext()
	depgraph.Designate("base.thickness", 0.75, ctx)
	doubled := depgraph.Define("base.doubled", func(in *depgraph.Context) float64 {
		return 2 * depgraph.Resolve[float64]("base.thickness", in)
	})
	assert.Equal(t, 1.5, doubled.Resolve(ctx))
	assert.Equal(t, 1.5, depgraph.Resolve[float64]("base.doubled", ctx))

	assert.PanicsWithValue(t, depgraph.ErrNotFound, func() {
		depgraph.Resolve[float64]("base.width", ctx)
	})
	assert.PanicsWithValue(t, lerr.ErrTypeMismatch{
		Expected: reflect.TypeFor[string](),
		Actual:   reflect.TypeFor[float64](),
	}, func() {
		depgraph.Resolve[string]("base.thickness", ctx)
	})
}

var wallHeight = depgraph.DesignateDefault("box.wallHeight", 4.0)

var boxHeight = depgraph.Define("box.height", func(in *depgraph.Context) float64 {
	return wallHeight.Resolve(in) + 0.25
})

func TestDesignateDefault(t *testing.T) {
	// A Context that designates the name uses its own value.
	tall := depgraph.NewContext()
	depgraph.Designate("box.wallHeight", 6.0, tall)
	assert.Equal(t, 6.25, boxHeight.Resolve(tall))

	// Otherwise the default is designated, so it can be updated like any
	// other Value.
	ctx := depgraph.NewContext()
	assert.Equal(t, 4.25, boxHeight.Resolve(ctx))
	walls := ctx.Resolve("box.wallHeight").(*depgraph.Value[float64])
	walls.Update(5)
	assert.Equal(t, 5.25, boxHeight.Resolve(ctx))

	mismatch := depgraph.NewContext()
	depgraph.Designate("box.wallHeight", 4, mismatch)
	assert.PanicsWithValue(t, lerr.ErrTypeMismatch{
		Expected: reflect.TypeFor[float64](),
		Actual:   reflect.TypeFor[int](),
	}, func() {
		wallHeight.Resolve(mismatch)
	})
}
