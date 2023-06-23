package corpus

// RootID represents a root word - a lower case alphanumeric string
type RootID uint32

type root struct {
	RootID
	str string
}
