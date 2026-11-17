package api

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// Model is what the write operations need to know about the graph: refs,
// base versions, siblings, handles and actors. It sees the write view: the
// as-is graph (every landed change set), then the in-flight change sets the
// current one depends on, directly or through others, then the current
// change set's ops on top.
//
// The projection's write model (writeModel, write_indexed.go) implements
// it. Nothing outside internal/api may depend on the implementation.
//
// Only Resolve, ResolveLink and LinkBetween give Version and Conflict;
// Node, LinksTo, LinksFrom, Children and ActorByEmail leave them zero, so a
// write that needs a base resolves its target first.
type Model interface {
	// Resolve finds a node by ref: the full canonical ID of a node; a
	// unique prefix, four characters or more, of the node IDs, matched as
	// given (IDs are upper case, so a lower-case ref never matches one);
	// the handle of a live node, in that order. Retired nodes resolve by ID
	// and ID prefix, never by handle. No match is a not_found Error with
	// next "rhiza ls"; several are a usage Error listing the candidates.
	Resolve(ref string) (Node, error)
	// Node returns the node with ID n, retired or not.
	Node(n id.ID) (Node, bool)
	// ResolveLink finds a link by its full ID or a unique ID prefix of at
	// least four characters, retired or not.
	ResolveLink(ref string) (Link, error)
	// LinkBetween finds the live link of type t from from to to. No match
	// is a not_found Error; several are a usage Error listing them.
	LinkBetween(from, to id.ID, t op.LinkType) (Link, error)
	// LinksTo returns the live links of type t whose To is to, sorted by
	// ID. It is nil when there are none.
	LinksTo(to id.ID, t op.LinkType) []Link
	// LinksFrom returns the live links of type t whose From is from, sorted
	// by ID. It is nil when there are none.
	LinksFrom(from id.ID, t op.LinkType) []Link
	// DependsOn returns what n depends on: the To of each live depends-on
	// link from n, sorted by ID, each once. It is nil when there are none.
	DependsOn(n id.ID) []id.ID
	// Children returns the live nodes of kind whose parent is parent, sorted
	// by key, then ID. A zero parent gives the roots of that kind.
	Children(kind op.Kind, parent id.ID) []Node
	// HandleTaken reports which live node has handle h, if any.
	HandleTaken(h string) (id.ID, bool)
	// ActorByEmail finds the live actor node whose email field is the JSON
	// string email.
	ActorByEmail(email string) (Node, bool)
}

// Node is a node in the to-be view.
type Node struct {
	ID     id.ID
	Kind   op.Kind
	Handle string
	// Parent is zero for a root and for kinds without a tree; Key is empty
	// for the latter.
	Parent id.ID
	Key    string
	Fields map[string]json.RawMessage
	// Retired is true once a node.retire applies, landed or current.
	Retired bool
	// Version is the last operation in the node's strict history that is
	// not the current change set's own: the one new writes name as base. It
	// is zero for a node created only in the current change set, and the
	// current change set never changes it.
	Version id.ID
	// Conflict is true when the view holds an unrepaired conflict
	// (model.ErrConflict) on the node from an operation not of the current
	// change set: a stale write that overlaps one it did not see.
	Conflict bool
	// Pending lists the IDs of the current change set's ops on the node, in
	// change set order; nil when it has none. A dependency's ops are never
	// in it.
	Pending []id.ID
}

// Link is a link in the to-be view. Version, Conflict and Pending mean what
// they mean on Node.
type Link struct {
	ID       id.ID
	Type     op.LinkType
	From     id.ID
	To       id.ID
	Note     string
	Decision id.ID
	Retired  bool
	Version  id.ID
	Conflict bool
	Pending  []id.ID
}

// dependencyOrder returns changeset.Closure(sets, current) with every change
// set after the ones it depends on among them, ties by ID. In a loop no
// member is first by that rule, so the smallest ID left goes next.
func dependencyOrder(sets []changeset.ChangeSet, current id.ID) []id.ID {
	left := changeset.Closure(sets, current)
	deps := map[id.ID][]id.ID{}
	for _, cs := range left {
		deps[cs] = changeset.DependsOn(sets, cs)
	}
	var out []id.ID
	for len(left) > 0 {
		i := slices.IndexFunc(left, func(cs id.ID) bool {
			return !slices.ContainsFunc(deps[cs], func(d id.ID) bool {
				return slices.Contains(left, d)
			})
		})
		if i < 0 {
			i = 0
		}
		out = append(out, left[i])
		left = slices.Delete(left, i, i+1)
	}
	return out
}

func ambiguousLinks(what string, ls []Link) error {
	names := make([]string, len(ls))
	for i, k := range ls {
		names[i] = k.ID.String()
	}
	return Usagef("%s matches %d links: %s", what, len(ls), strings.Join(names, ", "))
}
