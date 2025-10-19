package server_test

import (
	"bytes"
	"html/template"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/store/ephemeral/quicknested"
	server "github.com/adamcolton/luce/tools/server"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/lexec"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/gorilla/sessions"
	"github.com/quasoft/memstore"
	"github.com/stretchr/testify/assert"
)

// syncBuffer is a bytes.Buffer safe to read from one goroutine while the CLI
// loop writes from another.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func (s *syncBuffer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buf.Reset()
}

// testTemplates builds the minimal set of templates Server needs to render
// its built-in user pages.
func testTemplates() *template.Template {
	tmpl := template.Must(template.New("SignIn").Parse("signin: {{.Message}}"))
	template.Must(tmpl.New("Home").Parse("home"))
	template.Must(tmpl.New("HomeSignedIn").Parse("home, signed in"))
	return tmpl
}

// testConfig builds a Config backed by in-memory stores, ready for New.
func testConfig(addr string) server.Config {
	keyPairs := [][]byte{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	}
	cfg := server.Config{
		SessionStore: sessions.Store(memstore.NewMemStore(keyPairs...)),
		Templates:    testTemplates(),
		UserStore:    quicknested.New(10),
		TemplateNames: server.TemplateNames{
			SignIn:       "SignIn",
			Home:         "Home",
			HomeSignedIn: "HomeSignedIn",
		},
	}
	cfg.Addr = addr
	return cfg
}

// runCLI drives srv's admin CLI over an io.Pipe, sending each line in turn
// and giving the CLI time to respond before the next one. It returns every
// line written to the CLI's output after the test runs.
func runCLI(t *testing.T, srv *server.Server, lines []string) string {
	t.Helper()
	oldIn, oldOut := cli.StdIn, cli.StdOut
	defer func() { cli.StdIn, cli.StdOut = oldIn, oldOut }()
	pr, pw := io.Pipe()
	cli.StdIn = pr
	buf := &syncBuffer{}
	cli.StdOut = buf

	go srv.RunStdIO()
	err := timeout.After(3000, func() {
		for !strings.Contains(buf.String(), "> ") {
			time.Sleep(time.Millisecond)
		}
	})
	assert.NoError(t, err)

	for _, line := range lines {
		before := strings.Count(buf.String(), "> ")
		pw.Write([]byte(line + "\n"))
		err := timeout.After(3000, func() {
			for strings.Count(buf.String(), "> ") <= before {
				time.Sleep(time.Millisecond)
			}
		})
		assert.NoError(t, err, "command %q", line)
	}
	return buf.String()
}

func TestServerCLIUsersAndGroups(t *testing.T) {
	cfg := testConfig(":54401")
	srv, err := (&cfg).New()
	assert.NoError(t, err)

	out := runCLI(t, srv, []string{
		"u Name:bob Password:secret",
		"lu",
		"g Name:admins",
		"lg",
		"ug User:bob Group:admins",
		"s",
		// the CLI's bool parser (util/reflector/parsers.Bool) only treats
		// "y"/"Y" as true; anything else, including "true", is silently
		// false (a pre-existing cross-package quirk, not fixed here).
		"su AdminLockUserCreation:y",
		"s",
	})

	assert.Contains(t, out, "Created User")
	assert.Contains(t, out, "bob")
	assert.Contains(t, out, "Created Group")
	assert.Contains(t, out, "admins")
	assert.Contains(t, out, "Added User to Group")
	assert.Contains(t, out, "Admin lock setting updated")
	assert.Contains(t, out, "AdminLockUserCreation false")
	assert.Contains(t, out, "AdminLockUserCreation true")
}

func TestServerCLIUserGroupErrors(t *testing.T) {
	cfg := testConfig(":54402")
	srv, err := (&cfg).New()
	assert.NoError(t, err)

	out := runCLI(t, srv, []string{
		"ug User:nobody Group:nogroup",
		"u Name:amy Password:secret",
		"ug User:amy Group:nogroup",
	})
	assert.Contains(t, out, "group not found")
	assert.Contains(t, out, "Created User")
	// still "group not found": the group genuinely doesn't exist either time
	assert.Equal(t, 2, strings.Count(out, "group not found"))
}

func TestServerCLIRoutesAndServices(t *testing.T) {
	cfg := testConfig(":54403")
	srv, err := (&cfg).New()
	assert.NoError(t, err)

	out := runCLI(t, srv, []string{
		"sr",
		"ls",
	})
	// With no service connected over the service socket, both lists are
	// legitimately empty; this just confirms the commands run cleanly.
	assert.NotContains(t, out, "unknown command")
	assert.NotContains(t, out, "no handler found")
}

func TestServerCLIBashCommands(t *testing.T) {
	cfg := testConfig(":54404")
	cfg.BashCommands = []server.BashCmd{
		{Format: "echo hello"},
	}
	cfg.Commander = lexec.NewMock()
	srv, err := (&cfg).New()
	assert.NoError(t, err)

	out := runCLI(t, srv, []string{
		"lbc",
		"rbc ID:0",
		"sbc ID:0",
		"rbc ID:5",
	})
	assert.Contains(t, out, "echo hello")
	assert.Contains(t, out, "OK")
	assert.Contains(t, out, "ID out of range")
}

func TestBashCmdRunning(t *testing.T) {
	bc := server.BashCmd{}
	assert.False(t, bc.Running())
}

func TestServerCLISetPort(t *testing.T) {
	cfg := testConfig(":54405")
	srv, err := (&cfg).New()
	assert.NoError(t, err)

	closed := make(chan bool)
	go func() {
		srv.Run()
		closed <- true
	}()

	out := runCLI(t, srv, []string{
		"sp Port::54406",
	})
	assert.Contains(t, out, "Port changed")

	assert.NoError(t, srv.Close())
	err = timeout.After(300, closed)
	assert.NoError(t, err)
}
