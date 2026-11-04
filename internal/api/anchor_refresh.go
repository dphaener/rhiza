package api

// This file holds anchor.refresh: it re-fingerprints an existing anchor from
// the working tree into the current change set, the deliberate step after the
// anchored code's change is reviewed. An anchor whose fingerprint already
// matches writes nothing and says so.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/anchor"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

const anchorRefreshNext = "rhiza describe anchor.refresh"

func init() {
	register(Operation{
		Name:     "anchor.refresh",
		Summary:  "Re-fingerprint an anchor from the working tree into the current change set",
		Write:    true,
		Request:  reflect.TypeFor[AnchorRefreshRequest](),
		Response: reflect.TypeFor[AnchorRefreshResponse](),
		Run:      runAnchorRefresh,
	})
}

// AnchorRefreshRequest names the anchor.
type AnchorRefreshRequest struct {
	Anchor string `json:"anchor" arg:"0" desc:"The anchor: an ID, a unique ID prefix of 4+ characters, or a live node's handle"`
}

// AnchorRefreshResponse reports the refresh.
type AnchorRefreshResponse struct {
	Anchor      id.ID  `json:"anchor" desc:"The anchor"`
	Unchanged   bool   `json:"unchanged" desc:"True if the fingerprint already matched the working tree, so nothing was written"`
	Previous    string `json:"previous,omitempty" desc:"The fingerprint the anchor held; absent if it had none"`
	Fingerprint string `json:"fingerprint" desc:"The fingerprint of the working tree"`
	Entry       id.ID  `json:"entry,omitzero" desc:"The ops entry written; absent when unchanged and in a dry run"`
	ChangeSet   id.ID  `json:"change_set,omitzero" desc:"The change set written to; absent when unchanged"`
	DryRun      bool   `json:"dry_run,omitempty" desc:"True if nothing was written"`

	handle string
}

// Human says what was refreshed, or that nothing needed it.
func (r AnchorRefreshResponse) Human(w io.Writer) {
	name := r.handle
	if name == "" {
		name = r.Anchor.String()
	}
	switch {
	case r.Unchanged:
		fmt.Fprintf(w, "Anchor %s is unchanged (%s); nothing was written.\n", name, r.Fingerprint)
	case r.DryRun:
		fmt.Fprintf(w, "Dry run: would refresh anchor %s to %s; nothing was written.\n", name, r.Fingerprint)
	default:
		fmt.Fprintf(w, "Refreshed anchor %s to %s (entry %s, change set %s).\n", name, r.Fingerprint, r.Entry, r.ChangeSet)
	}
}

// unchangedAnchor carries the no-op result out of the write pipeline.
type unchangedAnchor struct{ resp AnchorRefreshResponse }

func (unchangedAnchor) Error() string { return "anchor unchanged" }

func runAnchorRefresh(ctx *Ctx, req any) (any, error) {
	r := req.(AnchorRefreshRequest)
	ref := strings.TrimSpace(r.Anchor)
	if ref == "" {
		return nil, Usagef("anchor refresh needs an anchor").WithNext(anchorRefreshNext)
	}
	var resp AnchorRefreshResponse
	out, err := runWrite(ctx, "anchor.refresh", commandLine("anchor refresh", ref), func(b writeBase) (edit, error) {
		n, err := writableNode(b, ref)
		if err != nil {
			return edit{}, err
		}
		if n.Kind != op.KindAnchor {
			return edit{}, Usagef("%s is %s %s, not an anchor", ref, article(string(n.Kind)), n.Kind).WithNext(anchorRefreshNext)
		}
		path := textField(n.Fields, "path")
		if path == "" {
			return edit{}, Usagef("anchor %s has no path to fingerprint", ref).
				WithNext("rhiza set " + ref + " --field path=<path relative to the project root>")
		}
		fp, err := anchor.Fingerprint(b.root, path, textField(n.Fields, "symbol"))
		if err != nil {
			return edit{}, anchorError(err)
		}
		prev := textField(n.Fields, "fingerprint")
		resp = AnchorRefreshResponse{Anchor: n.ID, Previous: prev, Fingerprint: fp, handle: n.Handle}
		if prev == fp {
			resp.Unchanged = true
			return edit{}, unchangedAnchor{resp}
		}
		return buildSet(b, SetRequest{Ref: ref, Fields: []string{"fingerprint=" + fp}})
	})
	if u, ok := err.(unchangedAnchor); ok {
		return u.resp, nil
	}
	if err != nil {
		return nil, err
	}
	wr := out.(WriteResponse)
	resp.Entry, resp.ChangeSet, resp.DryRun = wr.Entry, wr.ChangeSet, wr.DryRun
	return resp, nil
}
