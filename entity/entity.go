package entity

type Refser interface {
	EntRefs() []Key
}

type Entity interface {
	EntKey() Key
	EntVal(buf []byte) ([]byte, error)
	EntLoad(k Key, data []byte) error
}

// EntIniter allows a type to perform initilization when a new instance is
// created, before EntLoad is called
type EntIniter interface {
	EntInit()
}
