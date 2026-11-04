package slice

import (
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
)

// Shown returns a Reader, and a LinkReader, over r whose nodes carry their
// texts shown: each stored @{<ID>} as @<current handle> in r's view, or as
// @{<ID>} when the node is retired or not there, and a literal @ that would
// start a reference as @@ (see package textref).
func Shown(r Reader) Reader {
	return &shownReader{Reader: r, live: map[id.ID]string{}}
}

type shownReader struct {
	Reader
	live map[id.ID]string
}

func (s *shownReader) handleOf(i id.ID) (string, bool) {
	h, ok := s.live[i]
	if !ok {
		if n, found, err := s.Reader.Node(i); err == nil && found && !n.Retired {
			h = n.Handle
		}
		s.live[i] = h
	}
	return h, h != ""
}

func (s *shownReader) show(n model.Node) model.Node {
	n.Fields = model.ShowFields(n.Kind, n.Fields, s.handleOf)
	return n
}

func (s *shownReader) showAll(ns []model.Node) []model.Node {
	for i := range ns {
		ns[i] = s.show(ns[i])
	}
	return ns
}

func (s *shownReader) Node(i id.ID) (model.Node, bool, error) {
	n, ok, err := s.Reader.Node(i)
	return s.show(n), ok, err
}

func (s *shownReader) Ancestors(i id.ID) ([]model.Node, error) {
	ns, err := s.Reader.Ancestors(i)
	return s.showAll(ns), err
}

func (s *shownReader) Children(parent id.ID) ([]model.Node, error) {
	ns, err := s.Reader.Children(parent)
	return s.showAll(ns), err
}

func (s *shownReader) Link(i id.ID) (model.Link, bool, error) {
	if lr, ok := s.Reader.(LinkReader); ok {
		return lr.Link(i)
	}
	return model.Link{}, false, nil
}
