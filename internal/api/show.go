package api

// This file holds show: one node, with its fields, its links out and in
// and its children, or one link with both ends named, in a view. A ref is
// tried as a node first and then as a link.

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "show",
		Summary:  "Show a node or a link: its fields, links and children",
		Request:  reflect.TypeFor[ShowRequest](),
		Response: reflect.TypeFor[ShowResponse](),
		Run:      runShow,
	})
}

// The values of ShowResponse.Kind.
const (
	showKindNode = "node"
	showKindLink = "link"
)

// ShowRequest names the node or link to show.
type ShowRequest struct {
	Ref  string `json:"ref" arg:"0" desc:"The node or link: an ID, a unique ID prefix of 4+ characters, or a live node's handle"`
	View string `json:"view,omitempty" desc:"as-is, to-be or what-if:<change set>; default follows the effective change set"`
}

// ShowResponse is one node or one link: Kind says which, and only that
// one of Node and Link is present.
type ShowResponse struct {
	Kind string    `json:"kind" enum:"node,link" desc:"What the ref named: node or link"`
	Node *ShowNode `json:"node,omitempty" desc:"The node, when kind is node"`
	Link *ShowLink `json:"link,omitempty" desc:"The link, when kind is link"`
	View string    `json:"view" desc:"The view it was read in: as-is, to-be or what-if:<change set ID>"`
}

// ShowNode is a node as it stands in the view.
type ShowNode struct {
	ID       id.ID                      `json:"id" desc:"The node's ID"`
	Kind     op.Kind                    `json:"kind" desc:"The node's kind"`
	Handle   string                     `json:"handle" desc:"The node's handle"`
	Path     []ShowPathStep             `json:"path" desc:"Its ancestors in the tree, the root first; empty for a root or a kind outside the tree"`
	Retired  bool                       `json:"retired" desc:"Whether it is retired"`
	Reason   string                     `json:"reason,omitempty" desc:"Why it was retired"`
	Version  id.ID                      `json:"version" desc:"The last operation applied to it"`
	Created  id.ID                      `json:"created" desc:"The change set that created it"`
	Actor    string                     `json:"actor,omitempty" desc:"Evidence only: the actor stamped on the entry that created it, its handle (its ID when it has none); no field sets it"`
	Fields   map[string]json.RawMessage `json:"fields" desc:"Every field, name to canonical JSON value"`
	LinksOut []ShowNodeLink             `json:"links_out" desc:"The links from it, live and retired, sorted by ID"`
	LinksIn  []ShowNodeLink             `json:"links_in" desc:"The links to it, live and retired, sorted by ID"`
	Children []ShowChild                `json:"children" desc:"The nodes directly under it, live and retired, in tree order"`
	Notes    []ShowNote                 `json:"notes" desc:"The live notes that discuss it, oldest first"`
}

// ShowPathStep is one ancestor of a node.
type ShowPathStep struct {
	ID     id.ID  `json:"id" desc:"The ancestor's ID"`
	Handle string `json:"handle" desc:"The ancestor's handle"`
}

// ShowNodeLink is one link at a node, read from that node's end.
type ShowNodeLink struct {
	ID       id.ID       `json:"id" desc:"The link's ID"`
	Type     op.LinkType `json:"type" desc:"The link's type"`
	Name     string      `json:"name" desc:"The link read in this direction: its type going out, the type read backwards coming in"`
	Other    ShowEnd     `json:"other" desc:"The node or change set at the other end"`
	Note     string      `json:"note,omitempty" desc:"The link's note"`
	Decision id.ID       `json:"decision,omitzero" desc:"The decision behind an exempts link"`
	Retired  bool        `json:"retired" desc:"Whether the link is retired"`
}

// ShowEnd is one end of a link: a node with its handle, or a change set
// with its intent.
type ShowEnd struct {
	ID        id.ID  `json:"id" desc:"The node's or change set's ID"`
	Handle    string `json:"handle,omitempty" desc:"The node's handle"`
	ChangeSet bool   `json:"change_set,omitempty" desc:"True when the end is a change set"`
	Intent    string `json:"intent,omitempty" desc:"The change set's intent"`
}

// ShowChild is one node directly under the node shown.
type ShowChild struct {
	ID      id.ID   `json:"id" desc:"The child's ID"`
	Kind    op.Kind `json:"kind" desc:"The child's kind"`
	Handle  string  `json:"handle" desc:"The child's handle"`
	Retired bool    `json:"retired" desc:"Whether the child is retired"`
}

// ShowLink is a link as it stands in the view.
type ShowLink struct {
	ID       id.ID       `json:"id" desc:"The link's ID"`
	Type     op.LinkType `json:"type" desc:"The link's type"`
	From     ShowEnd     `json:"from" desc:"The node or change set it goes from"`
	To       ShowEnd     `json:"to" desc:"The node or change set it goes to"`
	Note     string      `json:"note,omitempty" desc:"The link's note"`
	Decision id.ID       `json:"decision,omitzero" desc:"The decision behind an exempts link"`
	Retired  bool        `json:"retired" desc:"Whether the link is retired"`
	Version  id.ID       `json:"version" desc:"The last operation applied to it"`
	Created  id.ID       `json:"created" desc:"The change set that created it"`
}

// runShow resolves the ref as a node, then as a link. Neither is
// not_found with next "rhiza ls"; an ambiguous node ref is usage.
func runShow(ctx *Ctx, req any) (any, error) {
	q := req.(ShowRequest)
	r, err := openReader(ctx, "show", q.View, defaultView)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	ref := strings.TrimSpace(q.Ref)
	if ref == "" {
		return nil, Usagef("show takes a node or link ref").WithNext("rhiza describe show")
	}
	e := &showEnds{r: r}
	n, err := r.Resolve(ref)
	if err == nil {
		sn, err := showNodeOf(r, e, n)
		if err != nil {
			return nil, err
		}
		return ShowResponse{Kind: showKindNode, Node: &sn, View: r.View()}, nil
	}
	if AsError(err).Code != CodeNotFound {
		return nil, err
	}
	l, lerr := r.ResolveLink(ref)
	if lerr != nil {
		if AsError(lerr).Code == CodeNotFound {
			return nil, NotFoundf("no node or link matches %q", ref).WithNext(lsNext)
		}
		return nil, lerr
	}
	from, err := e.end(l.From)
	if err != nil {
		return nil, err
	}
	to, err := e.end(l.To)
	if err != nil {
		return nil, err
	}
	return ShowResponse{Kind: showKindLink, View: r.View(), Link: &ShowLink{
		ID: l.ID, Type: l.Type, From: from, To: to, Note: l.Note, Decision: l.Decision,
		Retired: l.Retired, Version: l.Version, Created: l.Created,
	}}, nil
}

// showNodeOf gathers node n's path, fields, links and children.
func showNodeOf(r *reader, e *showEnds, n gmodel.Node) (ShowNode, error) {
	sn := ShowNode{
		ID: n.ID, Kind: n.Kind, Handle: n.Handle, Retired: n.Retired, Reason: n.Reason,
		Version: n.Version, Created: n.Created,
		Path: []ShowPathStep{}, Fields: map[string]json.RawMessage{},
		LinksOut: []ShowNodeLink{}, LinksIn: []ShowNodeLink{}, Children: []ShowChild{}, Notes: []ShowNote{},
	}
	if n.Kind == op.KindEvidence {
		a, err := showCreator(r, n)
		if err != nil {
			return ShowNode{}, err
		}
		sn.Actor = a
	}
	anc, err := r.Ancestors(n.ID)
	if err != nil {
		return ShowNode{}, err
	}
	for _, a := range anc {
		sn.Path = append(sn.Path, ShowPathStep{ID: a.ID, Handle: a.Handle})
	}
	for name, raw := range n.Fields {
		body, err := canon.Canonicalize(raw)
		if err != nil {
			return ShowNode{}, Errorf("node %s: field %s: %v", n.ID, name, err)
		}
		sn.Fields[name] = json.RawMessage(strings.TrimRight(string(body), "\n"))
	}
	ls, err := r.Links(n.ID, index.Both)
	if err != nil {
		return ShowNode{}, err
	}
	for _, l := range ls {
		if l.From == n.ID {
			nl, err := showLinkAt(e, l, string(l.Type), l.To)
			if err != nil {
				return ShowNode{}, err
			}
			sn.LinksOut = append(sn.LinksOut, nl)
		}
		if l.To == n.ID {
			nl, err := showLinkAt(e, l, backName(l.Type), l.From)
			if err != nil {
				return ShowNode{}, err
			}
			sn.LinksIn = append(sn.LinksIn, nl)
		}
	}
	cs, err := r.Children(n.ID)
	if err != nil {
		return ShowNode{}, err
	}
	for _, c := range cs {
		sn.Children = append(sn.Children, ShowChild{ID: c.ID, Kind: c.Kind, Handle: c.Handle, Retired: c.Retired})
	}
	if sn.Notes, err = showNotes(r, n.ID, sn.LinksIn); err != nil {
		return ShowNode{}, err
	}
	return sn, nil
}

// showCreator names the actor that created node n (gmodel.CreatedBy): it
// folds only the change set of n's create, from the index's entries. The
// name is the actor's handle in the view, else its ID; "" when the fold
// does not say who it was.
func showCreator(r *reader, n gmodel.Node) (string, error) {
	if len(n.History) == 0 {
		return "", nil
	}
	cs := n.History[0].ChangeSet
	es, err := r.Entries(cs)
	if err != nil {
		return "", err
	}
	c, _ := changeset.Fold(cs, es)
	a, ok := gmodel.CreatedBy(n, []changeset.ChangeSet{c})
	if !ok {
		return "", nil
	}
	an, ok, err := r.Node(a)
	if err != nil {
		return "", err
	}
	if ok && an.Handle != "" {
		return an.Handle, nil
	}
	return a.String(), nil
}

// showLinkAt is link l seen from one end: name is how it reads from there
// and other is the far end.
func showLinkAt(e *showEnds, l gmodel.Link, name string, other id.ID) (ShowNodeLink, error) {
	o, err := e.end(other)
	if err != nil {
		return ShowNodeLink{}, err
	}
	return ShowNodeLink{ID: l.ID, Type: l.Type, Name: name, Other: o, Note: l.Note, Decision: l.Decision, Retired: l.Retired}, nil
}

// showEnds names link ends, loading the change sets on first need.
type showEnds struct {
	r   *reader
	css map[id.ID]string
}

// end names the node or change set i; an ID that is neither is just its ID.
func (e *showEnds) end(i id.ID) (ShowEnd, error) {
	n, ok, err := e.r.Node(i)
	if err != nil {
		return ShowEnd{}, err
	}
	if ok {
		return ShowEnd{ID: i, Handle: n.Handle}, nil
	}
	if e.css == nil {
		infos, err := e.r.ChangeSets()
		if err != nil {
			return ShowEnd{}, err
		}
		e.css = map[id.ID]string{}
		for _, c := range infos {
			e.css[c.ID] = c.Intent
		}
	}
	if intent, ok := e.css[i]; ok {
		return ShowEnd{ID: i, ChangeSet: true, Intent: intent}, nil
	}
	return ShowEnd{ID: i}, nil
}

// showEndName names an end for people: "<handle> (<short ID>)", "change
// set <short ID>: <intent>", or the bare ID.
func showEndName(e ShowEnd) string {
	switch {
	case e.ChangeSet:
		return "change set " + shortID(e.ID) + ": " + e.Intent
	case e.Handle != "":
		return e.Handle + " (" + shortID(e.ID) + ")"
	}
	return e.ID.String()
}

// showLinkSuffix is a link line's note and retired mark.
func showLinkSuffix(note string, retired bool) string {
	s := ""
	if note != "" {
		s += " — " + strings.Join(strings.Fields(note), " ")
	}
	if retired {
		s += " [retired]"
	}
	return s
}

// Human shows the node or link as a heading and a field list, then for a
// node its Fields, Links and Children sections, each only when it has
// something in it.
func (r ShowResponse) Human(w io.Writer) {
	if l := r.Link; l != nil {
		fmt.Fprintf(w, "# link %s\n\n", shortID(l.ID))
		fmt.Fprintf(w, "- ID: %s\n", l.ID)
		fmt.Fprintf(w, "- Type: %s\n", l.Type)
		fmt.Fprintf(w, "- From: %s\n", showEndName(l.From))
		fmt.Fprintf(w, "- To: %s\n", showEndName(l.To))
		if l.Note != "" {
			fmt.Fprintf(w, "- Note: %s\n", strings.Join(strings.Fields(l.Note), " "))
		}
		if !l.Decision.IsZero() {
			fmt.Fprintf(w, "- Decision: %s\n", l.Decision)
		}
		if l.Retired {
			fmt.Fprintln(w, "- Retired: yes")
		}
		fmt.Fprintf(w, "- Version: %s\n", l.Version)
		fmt.Fprintf(w, "- Created: %s\n", l.Created)
		fmt.Fprintf(w, "- View: %s\n", r.View)
		return
	}
	n := r.Node
	if n == nil {
		return
	}
	fmt.Fprintf(w, "# %s %s\n\n", n.Kind, n.Handle)
	fmt.Fprintf(w, "- ID: %s\n", n.ID)
	if len(n.Path) > 0 {
		hs := make([]string, len(n.Path))
		for i, p := range n.Path {
			hs[i] = p.Handle
		}
		fmt.Fprintf(w, "- Path: %s\n", strings.Join(hs, " / "))
	}
	if n.Retired {
		fmt.Fprintf(w, "- Retired: %s\n", cmp.Or(strings.Join(strings.Fields(n.Reason), " "), "yes"))
	}
	fmt.Fprintf(w, "- Version: %s\n", n.Version)
	fmt.Fprintf(w, "- Created: %s\n", n.Created)
	if n.Actor != "" {
		fmt.Fprintf(w, "- Actor: %s\n", n.Actor)
	}
	fmt.Fprintf(w, "- View: %s\n", r.View)
	if len(n.Fields) > 0 {
		fmt.Fprintf(w, "\n## Fields\n")
		names := make([]string, 0, len(n.Fields))
		for name := range n.Fields {
			names = append(names, name)
		}
		slices.Sort(names)
		for _, name := range names {
			fmt.Fprintf(w, "\n### %s\n", name)
			if text := strings.TrimRight(fieldText(n.Fields[name]), "\n"); text != "" {
				fmt.Fprintf(w, "\n%s\n", text)
			}
		}
	}
	if len(n.LinksOut)+len(n.LinksIn) > 0 {
		fmt.Fprintf(w, "\n## Links\n\n")
		for _, ls := range [][]ShowNodeLink{n.LinksOut, n.LinksIn} {
			for _, l := range ls {
				fmt.Fprintf(w, "- %s %s%s\n", l.Name, showEndName(l.Other), showLinkSuffix(l.Note, l.Retired))
			}
		}
	}
	if len(n.Children) > 0 {
		fmt.Fprintf(w, "\n## Children\n\n")
		for _, c := range n.Children {
			fmt.Fprintf(w, "- %s %s (%s)%s\n", c.Kind, c.Handle, shortID(c.ID), showLinkSuffix("", c.Retired))
		}
	}
	writeNotes(w, n.Notes)
}
