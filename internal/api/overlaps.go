package api

// This file holds the overlaps operation: the overlaps between in-flight
// change sets (model.Overlaps), with the named draft counted as in flight.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	gmodel "github.com/dphaener/rhiza/internal/model"
)

func init() {
	register(Operation{
		Name:     "overlaps",
		Summary:  "List the overlaps between in-flight change sets",
		Request:  reflect.TypeFor[OverlapsRequest](),
		Response: reflect.TypeFor[OverlapsResponse](),
		Run:      runOverlaps,
	})
}

// OverlapsRequest takes nothing of its own: --change names the draft.
type OverlapsRequest struct{}

// OverlapInfo is one overlap.
type OverlapInfo struct {
	A       string   `json:"a" desc:"The first change set's ID"`
	B       string   `json:"b" desc:"The second change set's ID"`
	Target  string   `json:"target" desc:"The node or link both write, its handle when it has one"`
	Aspects []string `json:"aspects" desc:"The aspects that take part in the overlap, sorted"`
}

// OverlapsResponse is the overlaps.
type OverlapsResponse struct {
	Overlaps []OverlapInfo `json:"overlaps" desc:"The overlaps, by first change set, second change set and target; empty when none"`
}

// Human writes the overlaps as text.
func (r OverlapsResponse) Human(w io.Writer) {
	if len(r.Overlaps) == 0 {
		fmt.Fprintln(w, "No overlaps")
		return
	}
	for _, o := range r.Overlaps {
		fmt.Fprintf(w, "change sets %s and %s both write %s of %s\n", o.A, o.B, strings.Join(o.Aspects, ", "), o.Target)
	}
}

// runOverlaps folds the log and lists the overlaps. A change set given with
// --change counts as in flight even when it is a draft; without one no draft
// does.
func runOverlaps(ctx *Ctx, _ any) (any, error) {
	ws, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	var drafts []id.ID
	if !ctx.ChangeSet.IsZero() || strings.TrimSpace(ctx.ChangeRef) != "" {
		cs, err := ws.effectiveChangeSet(ctx)
		if err != nil {
			return nil, err
		}
		drafts = append(drafts, cs)
	}
	m := ws.model(id.ID{})
	resp := OverlapsResponse{Overlaps: []OverlapInfo{}}
	for _, ov := range gmodel.Overlaps(ws.changeSets(), drafts...) {
		target := ov.Target.String()
		if n, ok := m.Node(ov.Target); ok && n.Handle != "" {
			target = n.Handle
		}
		resp.Overlaps = append(resp.Overlaps, OverlapInfo{
			A: ov.A.String(), B: ov.B.String(), Target: target, Aspects: ov.Aspects,
		})
	}
	return resp, nil
}
