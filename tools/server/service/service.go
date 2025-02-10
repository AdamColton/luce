package service

import (
	"fmt"
	"path"

	"github.com/adamcolton/luce/lerr"
)

// Link is a named URL a Service exposes, relative to the server's host and
// port.
type Link struct {
	Name       string
	Host, Path string
}

// Get builds the Link's full URL against host and port.
func (l Link) Get(host, port string) string {
	if l.Host == "" {
		return l.Path
	}
	return fmt.Sprintf("https://%s.%s%s%s", l.Host, host, port, l.Path)
}

// Service describes a registered service: its name, host and base path,
// routes, and the links it exposes.
type Service struct {
	Name   string
	Host   string
	Base   string
	Routes []Route
	Links  []Link
}

// TypeID32 fulfills TypeIDer32. The ID was choosen at random.
func (*Service) TypeID32() uint32 {
	return 2516527266
}

// Validate validates every one of the Service's Routes.
func (s *Service) Validate() error {
	return lerr.NewSliceErrs(len(s.Routes), -1, func(i int) error {
		r := &(s.Routes[i])
		return r.Validate()
	})
}

// AddLink adds a Link named name, joining the Service's Base with pth.
func (s *Service) AddLink(name, host string, pth ...string) {
	wBase := make([]string, len(pth)+1)
	wBase[0] = s.Base
	copy(wBase[1:], pth)
	s.Links = append(s.Links, Link{
		Name: name,
		Host: host,
		Path: path.Join(wBase...),
	})
}
