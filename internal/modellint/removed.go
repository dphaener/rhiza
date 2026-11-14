package modellint

import (
	"io/fs"
	"path"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// KindRemoved is a live capability whose code the branch deleted.
const KindRemoved = "removed-code-capability"

// RemovedCode reports each live capability outside cs's slice (one cs does
// not write) that references at least one repo path, whose every path
// reference is absent from head, and at least one of whose path references
// is among deleted, the paths deleted on the branch (merge-base..HEAD). A
// reference to a directory counts as deleted when a deleted path is under
// it. Each hit names the capability and, in Ref, the first such deleted
// path reference as written. Hits are sorted by node.
func RemovedCode(m model.Model, cs changeset.ChangeSet, head fs.FS, deleted []string) []Hit {
	written := map[string]bool{}
	for _, i := range writtenNodes(cs) {
		written[i.String()] = true
	}
	c := &checker{head: head}
	isDeleted := func(ref string) bool {
		ref = path.Clean(ref)
		for _, d := range deleted {
			d = path.Clean(d)
			if d == ref || strings.HasPrefix(d, ref+"/") {
				return true
			}
		}
		return false
	}

	var out []Hit
	for _, n := range m.Nodes() {
		if n.Kind != op.KindCapability || n.Retired || written[n.ID.String()] {
			continue
		}
		paths, gone := 0, 0
		var first Hit
		seen := map[string]bool{}
		for _, f := range fields {
			for _, ref := range pathRefs(fieldText(n, f)) {
				if seen[ref] || !fs.ValidPath(ref) {
					continue
				}
				seen[ref] = true
				// A reference is a path when its first segment is in head,
				// or was deleted (a whole top-level directory may be gone).
				if !c.isPath(ref) && !isDeleted(ref) {
					continue
				}
				paths++
				if c.exists(ref) {
					continue
				}
				gone++
				if first.Kind == "" && isDeleted(ref) {
					first = Hit{Node: n.ID, Handle: n.Handle, Field: f, Kind: KindRemoved, Ref: ref}
				}
			}
		}
		if paths > 0 && gone == paths && first.Kind != "" {
			out = append(out, first)
		}
	}
	return out
}
