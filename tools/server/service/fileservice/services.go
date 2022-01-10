package fileservice

import (
	"net/http"
	"strings"

	"github.com/adamcolton/luce/tools/server/service"
)

// Services maps URLs to individual Clients and provides a set of shared
// Handlers.
type Services struct {
	DirHandler Handler
	Handlers   map[Extension]Handler
	byURL      map[string]*Service
}

// New creates an empty Services collection.
func New() *Services {
	return &Services{
		Handlers: make(map[Extension]Handler),
		byURL:    make(map[string]*Service),
	}
}

// New creates a Service.
func (s *Services) New(baseURL, rootDir string, handlers map[Extension]Handler, handleDir Handler) *Service {
	baseURL = "/" + strings.Trim(baseURL, "/")
	rootDir = strings.TrimRight(rootDir, "/")
	srv := &Service{
		BaseURL:   baseURL,
		RootDir:   rootDir,
		parent:    s,
		Handlers:  handlers,
		HandleDir: handleDir,
	}
	s.byURL[baseURL] = srv
	return srv
}

// Run the services.
func (s *Services) Run(addr string) {
	conn := service.MustClient(addr)

	s.Register(conn)

	conn.Run()
}

func (s *Services) Register(conn *service.Client) {
	for _, srv := range s.byURL {
		srv.Register(conn)
	}
}

// Handler finds the correct service to handle the request and returns
// it's response.
func (s *Services) Handler(req *service.Request) *service.Response {
	var srv *Service
	sIdx := strings.IndexRune(req.Path[1:], '/')
	if sIdx == -1 {
		srv = s.byURL[req.Path]
	} else {
		srv = s.byURL[req.Path[:sIdx+1]]
	}
	if srv == nil {
		return req.ResponseErr(nil, http.StatusNotFound)
	}
	return srv.Handler(req)
}
