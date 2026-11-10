package api

// This file holds ls: the nodes of a view, sorted by ID, filtered by kind
// and parent, a page at a time. The cursor is opaque to callers: it is the
// base64url (unpadded) form of the last ID of the page before.

import (
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "ls",
		ReadOnly: true,
		Summary:  "List nodes, filtered by kind and parent, a page at a time",
		Request:  reflect.TypeFor[LsRequest](),
		Response: reflect.TypeFor[LsResponse](),
		Run:      runLs,
	})
}

// The page sizes of ls.
const (
	lsDefaultLimit = 100
	lsMaxLimit     = 1000
)

// lsDescribeNext is the next command for a bad ls request.
const lsDescribeNext = "rhiza describe ls"

// LsRequest filters and pages the nodes.
type LsRequest struct {
	Kinds   []string `json:"kinds,omitempty" flag:"kind" enum:"intent,capability,constraint,decision,question,evidence,source,actor,anchor,note" desc:"Keep the nodes of this kind; repeat for several; default every kind"`
	Parent  string   `json:"parent,omitempty" desc:"Keep the nodes directly under this node: an ID, a unique ID prefix of 4+ characters, or a live node's handle"`
	Retired bool     `json:"retired,omitempty" desc:"Keep retired nodes too"`
	Limit   int      `json:"limit,omitempty" desc:"The most nodes to return, 1 to 1000; default 100"`
	Cursor  string   `json:"cursor,omitempty" desc:"Where to go on: the next_cursor of the page before"`
	View    string   `json:"view,omitempty" desc:"One of as-is, to-be or what-if:<change set>; default follows the effective change set"`
}

// LsResponse is one page of nodes.
type LsResponse struct {
	Nodes      []LsNode `json:"nodes" desc:"The nodes, sorted by ID"`
	NextCursor string   `json:"next_cursor,omitempty" desc:"The cursor of the next page; absent when no more remain"`
	View       string   `json:"view" desc:"The view it was read in: as-is, to-be or what-if:<change set ID>"`

	// command is the ls command line of this page, without --cursor, for
	// the next line.
	command string
}

// LsNode is one node of a page.
type LsNode struct {
	ID      id.ID   `json:"id" desc:"The node's ID"`
	Kind    op.Kind `json:"kind" desc:"The node's kind"`
	Handle  string  `json:"handle" desc:"The node's handle"`
	Parent  id.ID   `json:"parent,omitzero" desc:"The node's parent, if it has one"`
	Retired bool    `json:"retired" desc:"Whether it is retired"`
}

// runLs checks the request, then reads one row more than the page to know
// whether another page follows.
func runLs(ctx *Ctx, req any) (any, error) {
	q := req.(LsRequest)
	var kinds []op.Kind
	for _, k := range q.Kinds {
		if !slices.Contains(op.Kinds(), op.Kind(k)) {
			return nil, Usagef("unknown kind %q: want one of %s", k, joinKinds()).WithNext(lsDescribeNext)
		}
		if !slices.Contains(kinds, op.Kind(k)) {
			kinds = append(kinds, op.Kind(k))
		}
	}
	limit := q.Limit
	if limit == 0 {
		limit = lsDefaultLimit
	}
	if limit < 0 || limit > lsMaxLimit {
		return nil, Usagef("limit %d is out of range: give 1 to %d", q.Limit, lsMaxLimit).WithNext(lsDescribeNext)
	}
	var after id.ID
	if q.Cursor != "" {
		c, err := lsDecodeCursor(q.Cursor)
		if err != nil {
			return nil, err
		}
		after = c
	}
	r, err := openReader(ctx, "ls", q.View, defaultView)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	lq := index.ListQuery{Kinds: kinds, Retired: q.Retired, After: after, Limit: limit + 1}
	if q.Parent != "" {
		p, err := r.Resolve(strings.TrimSpace(q.Parent))
		if err != nil {
			return nil, err
		}
		lq.Parent, lq.HasParent = p.ID, true
	}
	ns, err := r.List(lq)
	if err != nil {
		return nil, err
	}
	resp := LsResponse{Nodes: []LsNode{}, View: r.View(), command: lsCommand(q)}
	if len(ns) > limit {
		ns = ns[:limit]
		resp.NextCursor = lsEncodeCursor(ns[limit-1].ID)
	}
	for _, n := range ns {
		resp.Nodes = append(resp.Nodes, LsNode{ID: n.ID, Kind: n.Kind, Handle: n.Handle, Parent: n.Parent, Retired: n.Retired})
	}
	return resp, nil
}

// lsEncodeCursor is the cursor that goes on after ID i.
func lsEncodeCursor(i id.ID) string {
	return base64.RawURLEncoding.EncodeToString([]byte(i.String()))
}

// lsDecodeCursor is the ID cursor c goes on after; anything that is not a
// cursor is a usage error.
func lsDecodeCursor(c string) (id.ID, error) {
	b, err := base64.RawURLEncoding.DecodeString(c)
	if err == nil {
		var i id.ID
		if i, err = id.Parse(string(b)); err == nil {
			return i, nil
		}
	}
	return id.ID{}, Usagef("bad cursor %q: give a next_cursor from ls", c).WithNext(lsDescribeNext)
}

// lsCommand is the command line that asks for q again, without its cursor.
func lsCommand(q LsRequest) string {
	parts := []string{"rhiza", "ls"}
	for _, k := range q.Kinds {
		parts = append(parts, "--kind", k)
	}
	if q.Parent != "" {
		parts = append(parts, "--parent", q.Parent)
	}
	if q.Retired {
		parts = append(parts, "--retired")
	}
	if q.Limit != 0 {
		parts = append(parts, "--limit", strconv.Itoa(q.Limit))
	}
	if q.View != "" {
		parts = append(parts, "--view", q.View)
	}
	return commandLine(parts...)
}

// Human lists the page one node a line, then the command for the next
// page when more remain.
func (r LsResponse) Human(w io.Writer) {
	if len(r.Nodes) == 0 {
		fmt.Fprintln(w, "No nodes match.")
		fmt.Fprintln(w, "\nnext: rhiza add")
		return
	}
	for _, n := range r.Nodes {
		line := fmt.Sprintf("- %s %s %s", n.Kind, n.Handle, n.ID)
		if n.Retired {
			line += " [retired]"
		}
		fmt.Fprintln(w, line)
	}
	if r.NextCursor != "" {
		cmd := r.command
		if cmd == "" {
			cmd = "rhiza ls"
		}
		fmt.Fprintf(w, "\nnext: %s --cursor %s\n", cmd, r.NextCursor)
	}
}
