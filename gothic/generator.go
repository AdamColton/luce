package gothic

// Prepper is a type that requires preperation before generation
type Preparer interface {
	Prepare() error
}

// Generator is expected to write any relevant data to a persistant medium when
// Generate is called
type Generator interface {
	Preparer
	Generate() error
}
