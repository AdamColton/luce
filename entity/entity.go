package entity

// TODO: I believe Refser is used to for garbage collection.
// An entity returns which other entities it has a reference to.
// but I don't have tests around this.
// Not quite true, enttest.Foo includes Refs[]
// which implements Refser
//
// So I think this should just be added to the Entity interface.
// Often, it will be implemented with a blank function,
// but it will save me (and possibly others) from forgetting.

type Entity interface {
	EntKey() Key
	EntVal(buf []byte) ([]byte, error)
	EntLoad(k Key, data []byte) error
	// EntRefs is unusual in that it is going to be invoked on a nil instance
	// Given the serialized data, return all entities this record refers to
	EntRefs(data []byte) ([]Key, error)
}

// EntIniter allows a type to perform initilization when a new instance is
// created, before EntLoad is called
type EntIniter interface {
	EntInit()
}
