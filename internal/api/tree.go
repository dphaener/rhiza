package api

// This file holds the tree operation: the capability tree (or the intent
// tree), or the subtree under one node, read with one recursive index query.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "tree",
		ReadOnly: true,
		Summary:  "Show the capability tree, or the subtree under a node",
		Request:  reflect.TypeFor[TreeRequest](),
		Response: reflect.TypeFor[TreeResponse](),
		Run:      runTree,
	})
}

// TreeRequest names the tree, or the node whose subtree to show.
type TreeRequest struct {
	Ref     string `json:"ref,omitempty" arg:"0" desc:"The node whose subtree to show: an ID, ID prefix or handle; default every root of the tree"`
	Kind    string `json:"kind,omitempty" enum:"capability,intent" desc:"Which tree: capability or intent; default capability, or the node's own kind with a ref"`
	Depth   int    `json:"depth,omitempty" desc:"How many levels below the top to show; 0 for all"`
	Retired bool   `json:"retired,omitempty" desc:"Show retired nodes too, and what is under them"`
	View    string `json:"view,omitempty" desc:"as-is, to-be or what-if:<change set>; default follows the effective change set"`
}

// TreeResponse is the tree in pre-order, siblings by key.
type TreeResponse struct {
	Root  id.ID          `json:"root,omitzero" desc:"The node whose subtree this is; absent for the whole tree"`
	Kind  op.Kind        `json:"kind" desc:"capability or intent"`
	Nodes []TreeNodeInfo `json:"nodes" desc:"The nodes, depth first in pre-order, siblings by key then ID"`
	View  string         `json:"view" desc:"The view read: as-is, to-be or what-if:<change set ID>"`
}

// TreeNodeInfo is one node of a tree.
type TreeNodeInfo struct {
	ID       id.ID  `json:"id" desc:"The node's ID"`
	Handle   string `json:"handle" desc:"The node's handle"`
	Depth    int    `json:"depth" desc:"0 for the top, one more per level down"`
	Retired  bool   `json:"retired,omitempty" desc:"True if the node is retired"`
	Children int    `json:"children" desc:"How many of its direct children are in the output"`
}

// treeKinds are the kinds the tree command walks.
var treeKinds = map[string]op.Kind{
	string(op.KindCapability): op.KindCapability,
	string(op.KindIntent):     op.KindIntent,
}

// runTree resolves the ref (if any), checks kind and depth, and reads the
// subtree with one SubtreeLight call, since it shows no fields. The whole
// tree is every root of both kinds, kept to the asked kind root by root: a
// node's parent is always of its own kind.
func runTree(ctx *Ctx, req any) (any, error) {
	q := req.(TreeRequest)
	kind, ok := treeKinds[q.Kind]
	if q.Kind != "" && !ok {
		return nil, Usagef("unknown tree kind %q: give capability or intent", q.Kind).
			WithNext("rhiza describe tree")
	}
	if q.Depth < 0 {
		return nil, Usagef("depth %d is negative: give 0 for the whole tree, or more", q.Depth).
			WithNext("rhiza describe tree")
	}
	r, err := openReader(ctx, "tree", q.View, defaultView)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var root id.ID
	if ref := strings.TrimSpace(q.Ref); ref != "" {
		n, err := r.Resolve(ref)
		if err != nil {
			return nil, err
		}
		if !op.HasParent(n.Kind) {
			return nil, Usagef("%s is %s %s, which is not in a tree", displayName(n), article(string(n.Kind)), n.Kind).
				WithNext("rhiza show " + ref)
		}
		if q.Kind != "" && kind != n.Kind {
			return nil, Usagef("%s is %s %s, not %s %s", displayName(n), article(string(n.Kind)), n.Kind, article(q.Kind), q.Kind).
				WithNext("rhiza tree " + ref)
		}
		root, kind = n.ID, n.Kind
	} else if kind == "" {
		kind = op.KindCapability
	}

	ts, err := r.SubtreeLight(root, q.Depth, q.Retired)
	if err != nil {
		return nil, err
	}
	resp := TreeResponse{Root: root, Kind: kind, Nodes: []TreeNodeInfo{}, View: r.View()}
	// open[d] is the index in resp.Nodes of the last node kept at depth d,
	// whose children count the next node at depth d+1 adds to.
	var open []int
	keep := true
	for _, t := range ts {
		if t.Depth == 0 {
			keep = t.Node.Kind == kind
		}
		if !keep {
			continue
		}
		open = append(open[:t.Depth], len(resp.Nodes))
		if t.Depth > 0 {
			resp.Nodes[open[t.Depth-1]].Children++
		}
		resp.Nodes = append(resp.Nodes, TreeNodeInfo{ID: t.Node.ID, Handle: t.Node.Handle,
			Depth: t.Depth, Retired: t.Node.Retired})
	}
	return resp, nil
}

// treePlural names the nodes of kind k for the empty tree's message.
func treePlural(k op.Kind) string {
	if k == op.KindCapability {
		return "capabilities"
	}
	return string(k) + "s"
}

// Human prints the tree as an indented list, two spaces per level.
func (r TreeResponse) Human(w io.Writer) {
	if len(r.Nodes) == 0 {
		if !r.Root.IsZero() {
			fmt.Fprintf(w, "Nothing to show: %s is retired.\n\nnext: rhiza tree %s --retired\n", r.Root, r.Root)
			return
		}
		fmt.Fprintf(w, "No %s yet.\n\nnext: rhiza add %s <handle>\n", treePlural(r.Kind), r.Kind)
		return
	}
	for _, n := range r.Nodes {
		line := strings.Repeat("  ", n.Depth) + "- " + n.Handle + " (" + shortID(n.ID) + ")"
		if n.Retired {
			line += " [retired]"
		}
		fmt.Fprintln(w, line)
	}
}
