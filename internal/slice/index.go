package slice

import (
	"errors"
	"fmt"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/model"
)

// FromIndex returns a Reader, and a LinkReader, over view v of x. It reads
// the same answers FromModel gives over the model the view stores: the
// index's queries keep the model's order (children by key then ID, links
// by ID, ancestors root first), and a short ID is found from the view's
// two neighbouring node IDs.
func FromIndex(x *index.Index, v index.View) Reader {
	return indexReader{x: x, v: v}
}

type indexReader struct {
	x *index.Index
	v index.View
}

func (r indexReader) Node(i id.ID) (model.Node, bool, error) {
	return r.x.Node(r.v, i)
}

func (r indexReader) Link(i id.ID) (model.Link, bool, error) {
	return r.x.Link(r.v, i)
}

func (r indexReader) Ancestors(i id.ID) ([]model.Node, error) {
	ns, err := r.x.Ancestors(r.v, i)
	if errors.Is(err, index.ErrNotFound) {
		return nil, fmt.Errorf("%w: no node %s", ErrNotFound, i)
	}
	return ns, err
}

func (r indexReader) Children(parent id.ID) ([]model.Node, error) {
	return r.x.Children(r.v, parent)
}

func (r indexReader) Links(i id.ID, dir Direction) ([]model.Link, error) {
	var d index.Direction
	switch dir {
	case From:
		d = index.From
	case To:
		d = index.To
	case Both:
		d = index.Both
	default:
		return nil, fmt.Errorf("slice: unknown direction %d", dir)
	}
	return r.x.Links(r.v, i, d)
}

func (r indexReader) ShortID(i id.ID) (string, error) {
	prev, next, err := r.x.Neighbours(r.v, i)
	if err != nil {
		return "", err
	}
	return Short(i, prev, next), nil
}
