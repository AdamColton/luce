package entity

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"golang.org/x/exp/maps"
)

// this just provides a small amount of type safety
type entRefer interface {
	isEntRef()
	Clear(bool)
	setIdx(int)
}

type cacheRecord struct {
	t   reflect.Type
	ref entRefer
}

// allRefs tracks all refs whose pointer is not nil
// first entry is nil so 0 can be used as "unset"
// kind of hacky, but works (for now)
// this is just here to make ClearCache work
var allRefs = []entRefer{nil}

var cache = lmap.EmptySafe[uint64, *cacheRecord](0)

func ClearCache() {
	for len(allRefs) > 1 {
		allRefs[1].Clear(false)
	}
	maps.Clear(cache.Map())
}

// T does not need to fulfill entity because it could be an interface where
// the underlying type fulfills both Entity and ...
// found indicates if it was found in cache, not if it exists in the entity
// store.
func Get[T any, E EntPtr[T]](key Key) (out *Ref[T, E], found bool) {
	var cr *cacheRecord
	cr, found = cache.Get(key.Hash64())
	//TODO: add typecheck here
	if found && cr.ref != nil {
		out, found = cr.ref.(*Ref[T, E])
		DeferStrategy.DeferCacheClear(out)
	}
	return
}

func NewRef[T any, E EntPtr[T]](ent E) *Ref[T, E] {
	r := &Ref[T, E]{
		key: ent.EntKey(),
		ent: entPtr[T, E]{p: ent},
	}
	r.addToAllRefs()
	DeferStrategy.DeferCacheClear(r)
	return r
}

func KeyRef[T any, E EntPtr[T]](k Key) *Ref[T, E] {
	return &Ref[T, E]{
		key: k,
	}
}

// ErrBadKey is thrown when attempting to Put an entity with no key or a key
// comprised of all zeros. This detects when keys are not being set correctly
// or when incrementing values are being used which will result in collision.
const ErrBadKey = lerr.Str("bad entity key")

// Put the entity into the entity cache. If there is already an a record for
// this entity, the pointer will be updated.
func Put[T any, E EntPtr[T]](ent E) *Ref[T, E] {
	k := ent.EntKey()

	// Temporary code to confirm keys are getting set correctly
	badKey := true
	for _, b := range k {
		badKey = b == 0
		if !badKey {
			break
		}
	}
	if badKey {
		panic(ErrBadKey)
	}

	h := k.Hash64()
	cr := cache.GetVal(h)
	var er *Ref[T, E]
	if cr == nil {
		cr = &cacheRecord{
			t: reflect.TypeOf(ent),
		}
		cache.Set(h, cr)
	}
	if cr.ref == nil {
		er = &Ref[T, E]{
			key: k,
			ent: entPtr[T, E]{p: ent},
		}
		er.addToAllRefs()
		cr.ref = er
	} else {
		er = cr.ref.(*Ref[T, E])
		er.ent = entPtr[T, E]{p: ent}
		er.addToAllRefs()
	}
	DeferStrategy.DeferCacheClear(er)

	return er
}

type weakGetter interface {
	GetEnt() (e Entity, ok bool)
	Referer
}

func GetEnt(k Key) (ent Entity, found bool) {
	var cr *cacheRecord
	cr, found = cache.Get(k.Hash64())
	//TODO: add typecheck here
	if found && cr.ref != nil {
		var getter weakGetter
		getter, found = cr.ref.(weakGetter)
		if found {
			ent, found = getter.GetEnt()
		}
		if found {
			DeferStrategy.DeferCacheClear(getter)
		}
	}
	return
}
