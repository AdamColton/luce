package service

import (
	"fmt"
	"path"

	"github.com/adamcolton/luce/lerr"
)

type Link struct {
	Name       string
	Host, Path string
}

func (l Link) Get(host, port string) string {
	if l.Host == "" {
		return l.Path
	}
	return fmt.Sprintf("https://%s.%s%s%s", l.Host, host, port, l.Path)
}

type Service struct {
	Name   string
	Host   string
	Base   string
	Routes []Route
	Links  []Link
}

func (*Service) TypeID32() uint32 {
	return 2516527266
}

func (s *Service) Validate() error {
	return lerr.NewSliceErrs(len(s.Routes), -1, func(i int) error {
		r := &(s.Routes[i])
		return r.Validate()
	})
}

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
