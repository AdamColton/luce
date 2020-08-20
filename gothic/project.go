package gothic

import "github.com/adamcolton/luce/lerr"

// State consts for generators
const (
	StateError = int8(iota - 1)
	StateReady
	StatePrepared
	StateGenerated
)

const (
	ErrPrepareCalled = lerr.Str("Prepare has already been called")
)

// BadState is a sentinal error
type BadState string

// Error fulfills error on BadState
func (bs BadState) Error() string { return string(bs) }

// Project is a collection of generators. It also implements the Generator
// interface, so projects can contain sub-projects.
type Project struct {
	state      int8
	generators []Generator
}

// New Project
func New() *Project {
	return &Project{}
}

// Prepare calls Prepare() on all Generators
func (p *Project) Prepare() error {
	if p.state != StateReady {
		if p.state == StatePrepared || p.state == StateGenerated {
			return ErrPrepareCalled
		}
		return BadState("Project is in an error state.")
	}
	for _, g := range p.generators {
		err := g.Prepare()
		if err != nil {
			p.state = StateError
			return err
		}
	}
	p.state = StatePrepared
	return nil
}

// Generate calls Generate() on all Generators
func (p *Project) Generate() error {
	if p.state != StatePrepared {
		if p.state == StateReady {
			return BadState("Cannot Generate, Prepare has not been called")
		}
		if p.state == StateGenerated {
			return BadState("Generate has already been called")
		}
		return BadState("Project is in an error state.")

	}
	for _, g := range p.generators {
		err := g.Generate()
		if err != nil {
			p.state = StateError
			return err
		}
	}
	p.state = StateGenerated
	return nil
}

// Export wraps a call to Prepare() and Generate()
func (p *Project) Export() error {
	err := p.Prepare()
	if err != nil {
		return err
	}
	return p.Generate()
}

// AddGenerators adds a generator to the project
func (p *Project) AddGenerators(g ...Generator) error {
	if p.state != StateReady {
		return BadState("Cannot Add Generator after Prepare has been called")
	}
	p.generators = append(p.generators, g...)
	return nil
}
