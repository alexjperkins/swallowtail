package service

import (
	"path"
	"strings"
)

type Kind int

const (
	KindService = iota
	KindGraphQLIAPIService
)

type Service interface {
	// The name of the service, e.g s.bar.
	Name() string
	// The kind of service.
	Kind() Kind
	// Path is the absolute path to the service with the supplied elements added.
	Path(elems ...string) string
	// Returns the Go import path for the service with the supplied elements added.
	ImportPath(elems ...string) string
}

func FromPath(dir string, kind Kind) Service {
	return &service{
		name: path.Base(dir),
		path: dir,
		kind: kind,
	}
}

type service struct {
	name, path string
	kind       Kind
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Kind() Kind {
	return s.kind
}

func (s *service) Path(elems ...string) string {
	elems = append([]string{s.path}, elems...)
	return path.Join(elems...)
}

func (s *service) ImportPath(elems ...string) string {
	elems = append([]string{"github.com/alexjperkins/swallowtail", s.name}, elems...)
	return strings.Join(elems, "/")
}
