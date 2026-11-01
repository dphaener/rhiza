package api

// change.diff reports a change set's field-level diff: per node it writes,
// the action, each changed field's old and new value taken from the node at
// the first op's base, and the links it adds and retires.

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
)

func init() {
	register(Operation{
		Name:     "change.diff",
		Summary:  "Show a change set's changes: per node, each field before and after",
		Request:  reflect.TypeFor[ChangeDiffRequest](),
		Response: reflect.TypeFor[ChangeDiffResponse](),
		Run:      runChangeDiff,
	})
}

// ChangeDiffRequest names the change set to diff.
type ChangeDiffRequest struct {
	Ref string `json:"ref,omitempty" arg:"0" desc:"The change set: a full ID or a unique prefix of 4+ characters; default --change, else the current one"`
}

// ChangeDiffResponse is the change set's diff.
type ChangeDiffResponse struct {
	ID    id.ID        `json:"id" desc:"The change set's ID"`
	Nodes []NodeChange `json:"nodes" desc:"One entry per node the change set writes, in the order it first writes them"`
	Links []LinkChange `json:"links" desc:"The links it adds and retires, in change set order"`
}

// NodeChange is the change set's writes to one node.
type NodeChange struct {
	Node    id.ID         `json:"node" desc:"The node's ID"`
	Action  string        `json:"action" enum:"created,updated,moved,retired" desc:"What the change set does to the node"`
	Kind    string        `json:"kind" desc:"The node's kind"`
	Handle  string        `json:"handle" desc:"The node's handle after the change set"`
	Reason  string        `json:"reason,omitempty" desc:"Why it is retired, for a retire"`
	Changes []FieldChange `json:"changes" desc:"Each aspect that differs; handle, parent and key are named @handle, @parent and @key"`
}

// FieldChange is one aspect before and after. A string is shown unquoted,
// any other value as compact JSON; an absent Old or New is omitted.
type FieldChange struct {
	Name string  `json:"name" desc:"The field's name"`
	Old  *string `json:"old,omitempty" desc:"The value at the first op's base; absent for a created node or a field not set there"`
	New  *string `json:"new,omitempty" desc:"The value after the change set; absent when the field ends unset"`
}

// LinkChange is a link the change set adds or retires.
type LinkChange struct {
	Link       id.ID  `json:"link" desc:"The link's ID"`
	Action     string `json:"action" enum:"added,retired" desc:"Whether the change set adds or retires the link"`
	Type       string `json:"type" desc:"The link type"`
	From       id.ID  `json:"from" desc:"The link's source"`
	To         id.ID  `json:"to" desc:"The link's target"`
	FromHandle string `json:"from_handle,omitempty" desc:"The source's handle, when it is a node"`
	ToHandle   string `json:"to_handle,omitempty" desc:"The target's handle, when it is a node"`
	Note       string `json:"note,omitempty" desc:"The link's note"`
}

func runChangeDiff(ctx *Ctx, req any) (any, error) {
	r := req.(ChangeDiffRequest)
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	var cs id.ID
	if ref := strings.TrimSpace(r.Ref); ref != "" {
		cs, err = w.resolveChangeSet(ref)
	} else {
		cs, err = w.effectiveChangeSet(ctx)
	}
	if err != nil {
		return nil, err
	}
	if cs.IsZero() {
		return nil, Usagef("no change set: name one or make one current").WithNext(changeLsNext)
	}
	if len(entriesOf(w.log(), cs)) == 0 {
		return nil, NotFoundf("no change set %s", cs).WithNext(changeLsNext)
	}
	return buildChangeDiff(cs, changeset.BuildDiff(w.changeSets(), cs)), nil
}

// buildChangeDiff turns a domain diff into the response.
func buildChangeDiff(cs id.ID, d changeset.Diff) ChangeDiffResponse {
	resp := ChangeDiffResponse{ID: cs, Nodes: []NodeChange{}, Links: []LinkChange{}}
	for _, n := range d.Nodes {
		nc := NodeChange{Node: n.Node, Action: string(n.Action), Kind: string(n.Kind),
			Handle: n.Handle, Reason: n.Reason, Changes: []FieldChange{}}
		for _, c := range n.Changes {
			nc.Changes = append(nc.Changes, FieldChange{Name: c.Name, Old: diffValue(c.Old), New: diffValue(c.New)})
		}
		resp.Nodes = append(resp.Nodes, nc)
	}
	for _, l := range d.Links {
		action := "added"
		if l.Retired {
			action = "retired"
		}
		resp.Links = append(resp.Links, LinkChange{Link: l.Link, Action: action, Type: string(l.Type),
			From: l.From, To: l.To, FromHandle: l.FromHandle, ToHandle: l.ToHandle, Note: l.Note})
	}
	return resp
}

// diffValue shows a canonical JSON value: a string unquoted, anything else
// as it is. Nil is absent.
func diffValue(raw json.RawMessage) *string {
	if raw == nil {
		return nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return &s
	}
	t := string(raw)
	return &t
}

// Human lists each node's changes, then the links.
func (r ChangeDiffResponse) Human(w io.Writer) {
	if len(r.Nodes) == 0 && len(r.Links) == 0 {
		fmt.Fprintln(w, "No changes.")
		return
	}
	for _, n := range r.Nodes {
		line := fmt.Sprintf("%s %s %s", n.Action, n.Kind, n.Handle)
		if n.Reason != "" {
			line += ": " + n.Reason
		}
		fmt.Fprintln(w, line)
		for _, c := range n.Changes {
			switch {
			case c.Old == nil:
				fmt.Fprintf(w, "  %s: %s\n", c.Name, *c.New)
			case c.New == nil:
				fmt.Fprintf(w, "  %s: %s -> unset\n", c.Name, *c.Old)
			case strings.Contains(*c.Old, "\n") || strings.Contains(*c.New, "\n"):
				fmt.Fprintf(w, "  %s:\n", c.Name)
				for _, l := range changeset.LineDiff(*c.Old, *c.New) {
					fmt.Fprintf(w, "    %c %s\n", l.Op, l.Text)
				}
			default:
				fmt.Fprintf(w, "  %s: %s -> %s\n", c.Name, *c.Old, *c.New)
			}
		}
	}
	for _, l := range r.Links {
		fmt.Fprintf(w, "%s link %s %s %s\n", l.Action, l.Type, nameOr(l.FromHandle, l.From), nameOr(l.ToHandle, l.To))
	}
}

func nameOr(h string, i id.ID) string {
	if h != "" {
		return h
	}
	return i.String()
}
