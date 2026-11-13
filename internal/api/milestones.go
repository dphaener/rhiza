package api

// This file holds the milestones read: every milestone of the project with
// its done points (its child intents), each with whether it is met and the
// change sets serving it in every state, and how many milestones are
// finished. roadmap stays the plan of proposed and accepted work; this is
// the whole picture the web roadmap shows.
//
// A done point is met when at least one change set serves it (directly or
// through a child intent) and every one of those that is not dropped
// (withdrawn, rejected or superseded) is landed; one with only dropped
// change sets is not met. A milestone is finished when it has done points
// and all are met; else building when it has an accepted or a landed change
// set; else planned when it has a proposed one; else designing.

import (
	"reflect"
	"slices"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "milestones",
		ReadOnly: true,
		Summary:  "Show every milestone with its done points and change sets in every state",
		Request:  reflect.TypeFor[MilestonesRequest](),
		Response: reflect.TypeFor[MilestonesResponse](),
		Run:      runMilestones,
	})
}

// MilestonesRequest asks for the milestones.
type MilestonesRequest struct {
	View string `json:"view,omitempty" desc:"The view: as-is, to-be or what-if:<change set>; default to-be, whatever the current change set"`
}

// MilestonesResponse is the milestones in intent-tree order.
type MilestonesResponse struct {
	Milestones []Milestone `json:"milestones" desc:"Every milestone of the project, in intent-tree order"`
	Finished   int         `json:"finished" desc:"How many of the milestones are finished"`
	View       string      `json:"view" desc:"The view read: as-is, to-be or what-if:<ID>"`
}

// Milestone is one milestone.
type Milestone struct {
	ID         id.ID                `json:"id" desc:"The milestone intent's ID"`
	Handle     string               `json:"handle,omitempty" desc:"The milestone intent's handle"`
	Statement  string               `json:"statement" desc:"The milestone intent's statement"`
	State      string               `json:"state" desc:"finished when it has done points and all are met; else building when it has an accepted or a landed change set; else planned when it has a proposed one; else designing"`
	DonePoints []MilestoneDonePoint `json:"done_points" desc:"Its child intents, in tree order"`
	ChangeSets []MilestoneChangeSet `json:"change_sets" desc:"Every change set serving the milestone or one of its intents, in every state, by ID"`
}

// MilestoneDonePoint is one done point of a milestone.
type MilestoneDonePoint struct {
	ID         id.ID                `json:"id" desc:"The done point intent's ID"`
	Handle     string               `json:"handle,omitempty" desc:"The done point intent's handle"`
	Statement  string               `json:"statement" desc:"The done point intent's statement"`
	Met        bool                 `json:"met" desc:"True when at least one change set serves it and every one that is not dropped is landed"`
	ChangeSets []MilestoneChangeSet `json:"change_sets" desc:"Every change set serving the done point or one of its child intents, in every state, by ID"`
}

// MilestoneChangeSet is one change set serving a milestone or a done point.
type MilestoneChangeSet struct {
	ID    id.ID           `json:"id" desc:"The change set's ID"`
	Title string          `json:"title" desc:"The change set's short title"`
	State changeset.State `json:"state" desc:"draft, proposed, accepted, landed, rejected, withdrawn or superseded"`
}

// milestoneStateOf is the state word of a milestone with done points dps and
// change sets css.
func milestoneStateOf(dps []MilestoneDonePoint, css []MilestoneChangeSet) string {
	met := len(dps) > 0
	for _, d := range dps {
		met = met && d.Met
	}
	if met {
		return "finished"
	}
	proposed := false
	for _, c := range css {
		switch c.State {
		case changeset.StateAccepted, changeset.StateLanded:
			return "building"
		case changeset.StateProposed:
			proposed = true
		}
	}
	if proposed {
		return "planned"
	}
	return "designing"
}

// donePointMet reports whether css make a done point met.
func donePointMet(css []MilestoneChangeSet) bool {
	landed := false
	for _, c := range css {
		switch {
		case c.State == changeset.StateLanded:
			landed = true
		case changeset.Dropped(c.State):
		default:
			return false
		}
	}
	return landed
}

func runMilestones(ctx *Ctx, req any) (any, error) {
	q := req.(MilestonesRequest)
	r, err := openReader(ctx, "milestones", q.View, toBeView)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// The live tree, with the milestones among its intents.
	in := planInputs{nodes: map[id.ID]gmodel.Node{}, intentPath: map[id.ID][]id.ID{},
		milestones: map[id.ID]bool{}, milestoneNums: map[id.ID]float64{}, statements: map[id.ID]string{}}
	tree, err := r.Subtree(id.ID{}, 0, false)
	if err != nil {
		return nil, err
	}
	var stack []id.ID
	var intents []id.ID
	for _, t := range tree {
		n := t.Node
		stack = append(stack[:min(t.Depth, len(stack))], n.ID)
		if n.Kind != op.KindIntent {
			continue
		}
		in.nodes[n.ID] = lightNode(n)
		in.intentPath[n.ID] = slices.Clone(stack[:len(stack)-1])
		planNoteStatement(in, n.ID, n.Fields)
		if planMilestoneField(n.Fields) {
			in.milestones[n.ID] = true
			planNoteMilestoneNum(in, n.ID, n.Fields)
		}
		intents = append(intents, n.ID)
	}

	// The change sets, and the in-flight ones folded: their titles and the
	// serves links their ops make.
	infos, err := r.ChangeSets()
	if err != nil {
		return nil, err
	}
	states := map[id.ID]changeset.State{}
	titles := map[id.ID]string{}
	for _, c := range infos {
		states[c.ID] = c.State
		titles[c.ID] = changeset.Open{Intent: c.Intent}.DisplayTitle()
	}
	flown := map[id.ID]changeset.ChangeSet{}
	if err := r.x.EachInFlight(func(c id.ID, es []entry.Entry) {
		s, _ := changeset.FoldCanonical(c, es)
		flown[c] = s
	}); err != nil {
		return nil, readError(err)
	}
	for c, s := range flown {
		if t := s.Open.DisplayTitle(); t != "" {
			titles[c] = t
		}
	}
	// A dropped change set's links are not in the view: its own ops make them.
	var dropped []id.ID
	for _, c := range infos {
		if changeset.Dropped(c.State) {
			dropped = append(dropped, c.ID)
		}
	}
	if len(dropped) > 0 {
		entries, err := r.EntriesOf(dropped)
		if err != nil {
			return nil, err
		}
		for _, s := range planFold(dropped, entries) {
			flown[s.ID] = s
		}
	}
	view, err := r.LinksTo(intents, op.LinkServes)
	if err != nil {
		return nil, err
	}
	inView := map[id.ID]bool{}
	var serves []gmodel.Link
	for _, l := range view {
		inView[l.ID] = true
		if !l.Retired {
			serves = append(serves, l)
		}
	}
	for c, s := range flown {
		for _, l := range planOwnLinks(c, s.Ops) {
			if l.Type == op.LinkServes && l.From == c && !inView[l.ID] {
				inView[l.ID] = true
				serves = append(serves, l)
			}
		}
	}
	slices.SortFunc(serves, func(a, b gmodel.Link) int { return id.Compare(a.ID, b.ID) })
	// servedBy is each live intent's change sets, once each.
	servedBy := map[id.ID][]id.ID{}
	for _, l := range serves {
		if _, ok := in.nodes[l.To]; !ok {
			continue
		}
		if _, ok := states[l.From]; !ok || slices.Contains(servedBy[l.To], l.From) {
			continue
		}
		servedBy[l.To] = append(servedBy[l.To], l.From)
	}
	row := func(cs id.ID) MilestoneChangeSet {
		return MilestoneChangeSet{ID: cs, Title: titles[cs], State: states[cs]}
	}

	var ms []id.ID
	for _, i := range intents {
		if in.milestones[i] {
			ms = append(ms, i)
		}
	}
	slices.SortFunc(ms, func(a, b id.ID) int { return planCompareIntents(a, b, in) })
	resp := MilestonesResponse{Milestones: []Milestone{}, View: r.View()}
	for _, m := range ms {
		out := Milestone{ID: m, Handle: in.nodes[m].Handle, Statement: in.statements[m],
			DonePoints: []MilestoneDonePoint{}, ChangeSets: []MilestoneChangeSet{}}
		var kids []id.ID
		for _, i := range intents {
			if in.nodes[i].Parent == m {
				kids = append(kids, i)
			}
		}
		slices.SortFunc(kids, func(a, b id.ID) int { return planCompareIntents(a, b, in) })
		for _, d := range kids {
			dp := MilestoneDonePoint{ID: d, Handle: in.nodes[d].Handle, Statement: in.statements[d], ChangeSets: []MilestoneChangeSet{}}
			var ids []id.ID
			for _, i := range intents {
				if i == d || slices.Contains(in.intentPath[i], d) {
					for _, c := range servedBy[i] {
						if !slices.Contains(ids, c) {
							ids = append(ids, c)
						}
					}
				}
			}
			slices.SortFunc(ids, id.Compare)
			for _, c := range ids {
				dp.ChangeSets = append(dp.ChangeSets, row(c))
			}
			dp.Met = donePointMet(dp.ChangeSets)
			out.DonePoints = append(out.DonePoints, dp)
		}
		var ids []id.ID
		for _, i := range intents {
			if mm, ok := planMilestoneOf(i, in); ok && mm == m {
				for _, c := range servedBy[i] {
					if !slices.Contains(ids, c) {
						ids = append(ids, c)
					}
				}
			}
		}
		slices.SortFunc(ids, id.Compare)
		for _, c := range ids {
			out.ChangeSets = append(out.ChangeSets, row(c))
		}
		out.State = milestoneStateOf(out.DonePoints, out.ChangeSets)
		if out.State == "finished" {
			resp.Finished++
		}
		resp.Milestones = append(resp.Milestones, out)
	}
	return resp, nil
}
