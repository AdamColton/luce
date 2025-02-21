package cli_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/stretchr/testify/assert"
)

// service makes a Commander, and is closed if it is an io.Closer.
type service struct {
	ec *cli.ExitClose
}

// ecCommander is a Commander that has the ExitClose it was made with.
type ecCommander struct {
	stubCommander
	ec *cli.ExitClose
}

func (c *ecCommander) EC() *cli.ExitClose { return c.ec }

func (s *service) NewCommander(ec *cli.ExitClose) cli.Commander {
	s.ec = ec
	return &ecCommander{
		stubCommander: stubCommander{cmds: []*handler.Command{{Name: "ping", Action: func() string { return "pong" }}}},
		ec:            ec,
	}
}

type closableService struct {
	service
	closed int
}

func (cs *closableService) Close() error {
	cs.closed++
	return nil
}

// wrappedService hides a closable service behind Wrapped.
type wrappedService struct {
	cli.CommanderFactory
	wrapped any
}

func (ws wrappedService) Wrapped() any { return ws.wrapped }

func TestSetupRunner(t *testing.T) {
	ctx, _ := newTestContext()
	exits := 0

	// Without a Close, the close command is not offered.
	svc := &service{}
	r := cli.SetupRunner("hello\n", svc, ctx, func() { exits++ })
	assert.Equal(t, "hello\n", r.StartMessage)
	assert.Same(t, svc.ec, r.ExitClose)
	assert.True(t, svc.ec.CanExit)
	assert.False(t, svc.ec.CanClose)
	svc.ec.OnExit()
	assert.Equal(t, 1, exits)

	// A service that can be closed is closed when the Runner is.
	closable := &closableService{}
	r = cli.SetupRunner("", closable, ctx, nil)
	assert.True(t, closable.ec.CanClose)
	assert.False(t, closable.ec.CanExit)
	r.OnClose()
	assert.Equal(t, 1, closable.closed)

	// Closers are found through wrappers.
	r = cli.SetupRunner("", wrappedService{CommanderFactory: closable, wrapped: closable}, ctx, nil)
	r.OnClose()
	assert.Equal(t, 2, closable.closed)
	_ = lerr.Str("")
}
