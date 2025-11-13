package entity

import (
	"reflect"
	"sync/atomic"
	"time"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/lerr"
)

type Refs []Key

func (r Refs) EntRefs() []Key {
	return r
}

type gcStruct struct {
	add         chan addReq
	q           chan addReq
	buf         []addReq
	runningLock uint32
	state       map[string]byte
	pass        byte
	roots       *lset.Set[string]
}

type action byte

const (
	undefined action = iota
	fromScan
	fromSave
	fromRoot
	fromSweep
	populateQ
)

func (a action) String() string {
	switch a {
	case undefined:
		return "undefined"
	case fromScan:
		return "fromScan"
	case fromSave:
		return "fromSave"
	case fromRoot:
		return "fromRoot"
	case fromSweep:
		return "fromSweep"
	case populateQ:
		return "populateQ"
	}
	return "unknown"
}

type addReq struct {
	Key
	action
}

var gc = &gcStruct{
	add:   make(chan addReq, 10),
	q:     make(chan addReq, 1),
	pass:  1,
	roots: lset.Safe[string](),
	state: make(map[string]byte),
}

func (gc *gcStruct) run() {
	if !atomic.CompareAndSwapUint32(&(gc.runningLock), 0, 1) {
		return
	}
	go func() {
		gc.add <- addReq{action: fromScan}
	}()
	for req := range gc.add {
		switch req.action {
		case populateQ:
			ln := len(gc.buf)
			if len(gc.q) == 0 && ln > 0 {
				ln--
				k := gc.buf[ln]
				gc.buf = gc.buf[:ln]
				gc.q <- k
			}
		default:
			if len(gc.q) == 0 {
				gc.q <- req
			} else {
				gc.buf = append(gc.buf, req)
			}
		}
	}
}

var stepTimeout = time.Microsecond * 100

func (gc *gcStruct) step() {
	var (
		k  Key
		ks string
	)
	for {
		// loop until we exhaust the q
		// or perform a store operation
		req, _ := channel.Timeout(stepTimeout, gc.q)
		if req.action == undefined {
			if len(gc.add) == 0 && len(gc.buf) == 0 && len(gc.q) == 0 {
				gc.finishPass()
				return
			}
			continue
		}
		gc.add <- addReq{action: populateQ} // steps run to move a value from the buffer
		k = req.Key
		ks = string(k)

		if req.action == fromScan {
			next := entstore.Next(k)
			if next != nil {
				ns := string(next)
				if _, found := gc.state[ns]; !found {
					gc.state[ns] = 0
				}
				gc.add <- addReq{Key: next, action: fromScan}
			}
			return
		}

		if gc.state[ks] != gc.pass {
			break
		}
	}

	gc.state[ks] = gc.pass
	rec := entstore.Get(k)
	if !rec.Found {
		return
	}

	t, data, err := typer.GetType(rec.Value)
	lerr.Panic(lerr.Wrap(err, "get type failed during ent gc"))
	i := reflect.New(t).Elem().Interface().(Entity)

	refs, err := i.EntRefs(data)
	lerr.Panic(err) // TODO: I don't know, something
	for _, k := range refs {
		gc.add <- addReq{Key: k, action: fromSweep}
	}
}

var garbage []Key

func (gc *gcStruct) finishPass() {
	var updateState []string
	var newGarbage []Key
	for ks, s := range gc.state {
		if s == 0 {
			updateState = append(updateState, ks)
		} else if s != gc.pass {
			newGarbage = append(newGarbage, Key(ks))
		}
	}
	for _, ks := range updateState {
		gc.state[ks] = gc.pass
	}
	gc.roots.Each(func(k string, done *bool) {
		gc.add <- addReq{Key: Key(k), action: fromRoot}
	})
	garbage = newGarbage
	gc.pass = 3 - gc.pass
}

func Garbage() []Key {
	if gc.runningLock == 0 {
		go gc.run()
	}
	p := gc.pass
	for gc.pass == p {
		gc.step()
	}
	return garbage
}

func AddGCRoots(keys ...Key) {
	if gc.runningLock == 0 {
		go gc.run()
	}
	for _, k := range keys {
		gc.roots.Add(string(k))
		gc.add <- addReq{Key: k, action: fromRoot}
	}
}
