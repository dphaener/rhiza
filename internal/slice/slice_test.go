package slice

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

// mkID is fixture ID n of family f: IDs of one family grow with n.
func mkID(f byte, n int) id.ID {
	var i id.ID
	i[0] = f
	binary.BigEndian.PutUint32(i[12:], uint32(n))
	return i
}

func csID(n int) id.ID { return mkID(2, n) }
func opID(n int) id.ID { return mkID(4, n) }
func nd(n int) id.ID   { return mkID(5, n) }
func ln(n int) id.ID   { return mkID(6, n) }

// fx builds change sets op by op, tracking each target's version so that
// updates, moves and retires name the right base.
type fx struct {
	n    int
	ver  map[id.ID]id.ID
	sets []changeset.ChangeSet
}

// newFx starts with one landed change set.
func newFx() *fx {
	f := &fx{ver: map[id.ID]id.ID{}}
	f.begin(changeset.StateLanded, "base")
	return f
}

// begin starts change set csID(k) in state s; later ops go into it.
func (f *fx) begin(s changeset.State, intent string) id.ID {
	cs := csID(len(f.sets) + 1)
	f.sets = append(f.sets, changeset.ChangeSet{ID: cs, State: s, Open: changeset.Open{Intent: intent}})
	return cs
}

func (f *fx) add(o op.Op) {
	f.n++
	o.ID = opID(f.n)
	if o.Type != op.TypeNodeCreate && o.Type != op.TypeLinkCreate {
		o.Base = f.ver[o.Target()]
	}
	last := &f.sets[len(f.sets)-1]
	last.Ops = append(last.Ops, o)
	f.ver[o.Target()] = o.ID
}

// ref is nd(n), or zero for 0.
func ref(n int) id.ID {
	if n == 0 {
		return id.ID{}
	}
	return nd(n)
}

func fields(kv ...string) map[string]json.RawMessage {
	if len(kv) == 0 {
		return nil
	}
	m := map[string]json.RawMessage{}
	for i := 0; i < len(kv); i += 2 {
		m[kv[i]] = json.RawMessage(kv[i+1])
	}
	return m
}

func q(s string) string { return fmt.Sprintf("%q", s) }

// node creates nd(n) of kind k; parent applies only to tree kinds.
func (f *fx) node(n int, k op.Kind, handle string, parent int, kv ...string) {
	o := op.Op{Type: op.TypeNodeCreate, Node: nd(n), Kind: k, Handle: handle, Fields: fields(kv...)}
	if op.HasParent(k) {
		o.Parent = ref(parent)
		o.Key = handle
	}
	f.add(o)
}

// cap creates capability nd(n) under nd(parent) with a statement.
func (f *fx) cap(n int, handle string, parent int) {
	f.node(n, op.KindCapability, handle, parent, "statement", q(handle+" does things"))
}

func (f *fx) link(n int, t op.LinkType, from, to int, note string, decision int) {
	f.add(op.Op{Type: op.TypeLinkCreate, Link: ln(n), LinkType: t, From: nd(from), To: nd(to), Note: note,
		Decision: ref(decision)})
}

func (f *fx) set(n int, kv ...string) {
	f.add(op.Op{Type: op.TypeNodeUpdate, Node: nd(n), Set: fields(kv...)})
}

func (f *fx) move(n, parent int) {
	f.add(op.Op{Type: op.TypeNodeMove, Node: nd(n), Parent: ref(parent), Key: "k"})
}

func (f *fx) retire(n int, reason string) {
	f.add(op.Op{Type: op.TypeNodeRetire, Node: nd(n), Reason: reason})
}

func (f *fx) retireLink(n int) { f.add(op.Op{Type: op.TypeLinkRetire, Link: ln(n)}) }

func (f *fx) asIs(t *testing.T) Reader {
	t.Helper()
	m, ps := model.AsIs(f.sets)
	if ps != nil {
		t.Fatalf("fixture problems: %v", ps)
	}
	return FromModel(m)
}

func (f *fx) toBe(t *testing.T) Reader {
	t.Helper()
	m, ps := model.ToBe(f.sets)
	if ps != nil {
		t.Fatalf("fixture problems: %v", ps)
	}
	return FromModel(m)
}

func (f *fx) whatIf(t *testing.T, cs id.ID) Reader {
	t.Helper()
	m, ps := model.WhatIf(f.sets, cs)
	if ps != nil {
		t.Fatalf("fixture problems: %v", ps)
	}
	return FromModel(m)
}

// cs returns the change set with ID i.
func (f *fx) cs(i id.ID) *changeset.ChangeSet {
	for k := range f.sets {
		if f.sets[k].ID == i {
			return &f.sets[k]
		}
	}
	return nil
}

// tree is the fixture most cases start from: runtime (1) > governor (2) >
// pause (3), with pause's siblings resume (4).
func tree() *fx {
	f := newFx()
	f.cap(1, "runtime", 0)
	f.cap(2, "governor", 1)
	f.cap(3, "pause", 2)
	f.cap(4, "resume", 2)
	return f
}

// show writes items as "handle[ from=..][ status=..][ rel=..][ note=..]".
func show(items []Item) []string {
	var out []string
	for _, it := range items {
		s := it.Handle
		if it.From != "" {
			s += " from=" + it.From
		}
		if it.Status != "" {
			s += " status=" + it.Status
		}
		if it.Relation != "" {
			s += " rel=" + it.Relation
		}
		if it.Note != "" {
			s += " note=" + it.Note
		}
		out = append(out, s)
	}
	return out
}

func eq(t *testing.T, what string, got, want []string) {
	t.Helper()
	if len(got) == 0 && len(want) == 0 {
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s =\n  %q\nwant\n  %q", what, got, want)
	}
}

func build(t *testing.T, r, asIs Reader, req Request) Slice {
	t.Helper()
	s, err := Build(r, asIs, req)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	return s
}

func nodeSlice(t *testing.T, f *fx, focus int) Slice {
	t.Helper()
	r := f.asIs(t)
	return build(t, r, r, Request{Focus: nd(focus), View: "as-is"})
}

func TestBuild(t *testing.T) {
	t.Run("inherit constraint direct and from root", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/boundary", 0, "statement", q("only the provider layer names a provider"))
		f.node(11, op.KindConstraint, "con/reserve", 0, "statement", q("keep 10%"))
		f.node(12, op.KindConstraint, "con/unlinked", 0)
		f.node(13, op.KindConstraint, "con/sibling", 0)
		f.link(1, op.LinkConstrains, 10, 1, "", 0)
		f.link(2, op.LinkConstrains, 11, 3, "", 0)
		f.link(3, op.LinkConstrains, 13, 4, "", 0)
		s := nodeSlice(t, f, 3)
		eq(t, "constraints", show(s.Constraints), []string{"con/boundary from=runtime", "con/reserve"})
		if got := s.Constraints[0].Statement; got != "only the provider layer names a provider" {
			t.Errorf("statement = %q", got)
		}
	})

	t.Run("inherit constraint keeps the nearest from", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/a", 0)
		f.link(1, op.LinkConstrains, 10, 1, "", 0)
		f.link(2, op.LinkConstrains, 10, 2, "", 0)
		eq(t, "constraints", show(nodeSlice(t, f, 3).Constraints), []string{"con/a from=governor"})
		eq(t, "constraints on governor", show(nodeSlice(t, f, 2).Constraints), []string{"con/a"})
	})

	t.Run("inherit drops retired constraints and links", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/a", 0)
		f.node(11, op.KindConstraint, "con/b", 0)
		f.link(1, op.LinkConstrains, 10, 1, "", 0)
		f.link(2, op.LinkConstrains, 11, 1, "", 0)
		f.retire(10, "obsolete")
		f.retireLink(2)
		eq(t, "constraints", show(nodeSlice(t, f, 3).Constraints), nil)
	})

	t.Run("inherit decisions and intents", func(t *testing.T) {
		f := tree()
		f.node(20, op.KindDecision, "dec/no-overage", 0, "statement", q("refuse to run"))
		f.node(21, op.KindDecision, "dec/local", 0)
		f.node(30, op.KindIntent, "int/reserve", 0, "statement", q("keep a reserve"))
		f.node(31, op.KindIntent, "int/fast", 0)
		f.link(1, op.LinkDecidedBy, 2, 20, "", 0)
		f.link(2, op.LinkDecidedBy, 3, 21, "", 0)
		f.link(3, op.LinkServes, 1, 30, "", 0)
		f.link(4, op.LinkServes, 3, 31, "", 0)
		s := nodeSlice(t, f, 3)
		eq(t, "decisions", show(s.Decisions), []string{"dec/local", "dec/no-overage from=governor"})
		eq(t, "context", show(s.Context), []string{"runtime rel=ancestor", "governor rel=ancestor",
			"int/fast rel=serves", "int/reserve from=runtime rel=serves"})
	})

	t.Run("exempt with a live decision removes the constraint", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/a", 0)
		f.node(11, op.KindConstraint, "con/b", 0)
		f.node(20, op.KindDecision, "dec/allow", 0)
		f.link(1, op.LinkConstrains, 10, 1, "", 0)
		f.link(2, op.LinkConstrains, 11, 1, "", 0)
		f.link(3, op.LinkExempts, 2, 10, "allowed here", 20)
		eq(t, "constraints of pause", show(nodeSlice(t, f, 3).Constraints), []string{"con/b from=runtime"})
		eq(t, "constraints of runtime", show(nodeSlice(t, f, 1).Constraints), []string{"con/a", "con/b"})
	})

	t.Run("exempt without a decision exempts nothing", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/a", 0)
		f.link(1, op.LinkConstrains, 10, 1, "", 0)
		f.link(2, op.LinkExempts, 3, 10, "please", 0)
		eq(t, "constraints", show(nodeSlice(t, f, 3).Constraints), []string{"con/a from=runtime"})
	})

	t.Run("exempt naming a retired decision exempts nothing", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/a", 0)
		f.node(20, op.KindDecision, "dec/allow", 0)
		f.link(1, op.LinkConstrains, 10, 1, "", 0)
		f.link(2, op.LinkExempts, 3, 10, "allowed", 20)
		f.retire(20, "reversed")
		eq(t, "constraints", show(nodeSlice(t, f, 3).Constraints), []string{"con/a from=runtime"})
	})

	t.Run("superseded decision dropped", func(t *testing.T) {
		f := tree()
		f.node(20, op.KindDecision, "dec/old", 0)
		f.node(21, op.KindDecision, "dec/new", 0)
		f.node(22, op.KindDecision, "dec/kept", 0)
		f.node(23, op.KindDecision, "dec/dead-replacement", 0)
		f.link(1, op.LinkDecidedBy, 3, 20, "", 0)
		f.link(2, op.LinkDecidedBy, 3, 21, "", 0)
		f.link(3, op.LinkDecidedBy, 3, 22, "", 0)
		f.link(4, op.LinkSupersedes, 21, 20, "", 0)
		f.link(5, op.LinkSupersedes, 23, 22, "", 0)
		f.retire(23, "withdrawn")
		eq(t, "decisions", show(nodeSlice(t, f, 3).Decisions), []string{"dec/kept", "dec/new"})
	})

	t.Run("answered question dropped and blocking before relates-to", func(t *testing.T) {
		f := tree()
		f.node(40, op.KindQuestion, "q/z-blocks", 0)
		f.node(41, op.KindQuestion, "q/a-related", 0)
		f.node(42, op.KindQuestion, "q/answered", 0)
		f.node(43, op.KindQuestion, "q/m-blocks-root", 0)
		f.node(44, op.KindQuestion, "q/dead-answer", 0)
		f.node(20, op.KindDecision, "dec/answer", 0)
		f.node(21, op.KindDecision, "dec/retired-answer", 0)
		f.link(1, op.LinkBlocks, 40, 3, "", 0)
		f.link(2, op.LinkRelatesTo, 3, 41, "unclear limit", 0)
		f.link(3, op.LinkBlocks, 42, 3, "", 0)
		f.link(4, op.LinkAnswers, 20, 42, "", 0)
		f.link(5, op.LinkBlocks, 43, 1, "", 0)
		f.link(6, op.LinkBlocks, 44, 3, "", 0)
		f.link(7, op.LinkAnswers, 21, 44, "", 0)
		f.retire(21, "wrong")
		eq(t, "questions", show(nodeSlice(t, f, 3).Questions), []string{
			"q/dead-answer rel=blocks", "q/m-blocks-root from=runtime rel=blocks", "q/z-blocks rel=blocks",
			"q/a-related rel=relates-to note=unclear limit"})
	})

	t.Run("anchors and evidence status combos", func(t *testing.T) {
		cases := []struct {
			name     string
			statuses []string
			want     string
		}{
			{"none", nil, "none"},
			{"unverified", []string{`"draft"`}, "unverified"},
			{"verified", []string{`"verified"`, `"draft"`}, "verified"},
			{"failing wins", []string{`"verified"`, `"failing"`, `1`}, "failing"},
			{"non-string", []string{`true`}, "unverified"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				f := tree()
				f.node(50, op.KindAnchor, "anc/decide", 0, "statement", q("internal/governor/decide.go · Decide"))
				f.link(1, op.LinkImplements, 50, 3, "", 0)
				for i, st := range c.statuses {
					f.node(60+i, op.KindEvidence, fmt.Sprintf("ev/%d", i), 0, "status", st)
					f.link(10+i, op.LinkVerifies, 60+i, 3, "", 0)
				}
				s := nodeSlice(t, f, 3)
				eq(t, "anchors", show(s.Anchors), []string{"anc/decide"})
				if s.Anchors[0].Statement != "internal/governor/decide.go · Decide" {
					t.Errorf("anchor statement = %q", s.Anchors[0].Statement)
				}
				if got := s.Focus[0].Status; got != c.want {
					t.Errorf("focus evidence = %q, want %q", got, c.want)
				}
				var want []string
				for i, st := range c.statuses {
					want = append(want, fmt.Sprintf("ev/%d status=%s", i, wantEvidence(st)))
				}
				eq(t, "evidence", show(s.Evidence), want)
			})
		}
	})

	t.Run("children with statement and evidence status", func(t *testing.T) {
		f := tree()
		f.cap(5, "audit", 2)
		f.node(60, op.KindEvidence, "ev/pause", 0, "status", q("failing"))
		f.link(1, op.LinkVerifies, 60, 3, "", 0)
		f.retire(5, "merged")
		s := nodeSlice(t, f, 2)
		eq(t, "children", show(s.Children), []string{"pause status=failing", "resume status=none"})
		if s.Children[0].Statement != "pause does things" {
			t.Errorf("child statement = %q", s.Children[0].Statement)
		}
	})

	t.Run("neighbours both ways", func(t *testing.T) {
		f := tree()
		f.cap(5, "metering", 0)
		f.cap(6, "scheduler", 0)
		f.cap(7, "billing", 0)
		f.cap(8, "dead", 0)
		f.link(1, op.LinkDependsOn, 3, 5, "", 0)
		f.link(2, op.LinkDependsOn, 6, 3, "", 0)
		f.link(3, op.LinkRelatesTo, 7, 3, "shares the window", 0)
		f.link(4, op.LinkDependsOn, 3, 8, "", 0)
		f.retire(8, "gone")
		eq(t, "neighbours", show(nodeSlice(t, f, 3).Neighbours), []string{
			"billing rel=relates-to note=shares the window", "metering rel=needs", "scheduler rel=needed by"})
	})

	t.Run("in flight lists other change sets targeting focus", func(t *testing.T) {
		f := tree()
		r := f.asIs(t)
		req := Request{Focus: nd(3), View: "as-is", InFlight: []InFlight{
			{ID: csID(9), State: changeset.StateProposed, Intent: "Other", Targets: []id.ID{nd(4)}},
			{ID: csID(8), State: changeset.StateDraft, Intent: "Pause at 90%", Targets: []id.ID{nd(1), nd(3)}},
			{ID: csID(7), State: changeset.StateAccepted, Intent: "Also", Targets: []id.ID{nd(3)}},
		}}
		s := build(t, r, r, req)
		var got []string
		for _, c := range s.InFlight {
			got = append(got, fmt.Sprintf("%s %s %s", c.ID, c.State, c.Intent))
		}
		eq(t, "in flight", got, []string{csID(7).String() + " accepted Also", csID(8).String() + " draft Pause at 90%"})
		req.Change = &changeset.ChangeSet{ID: csID(8), State: changeset.StateDraft}
		if s := build(t, r, r, req); len(s.InFlight) != 1 || s.InFlight[0].ID != csID(7) {
			t.Errorf("in flight with change 8 = %+v, want only change 7", s.InFlight)
		}
	})

	t.Run("change union de-duplicated keeps the first from", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/a", 0)
		f.link(1, op.LinkConstrains, 10, 2, "", 0)
		f.node(20, op.KindDecision, "dec/d", 0)
		cs := f.begin(changeset.StateDraft, "Pause at 90%")
		f.set(3, "statement", q("stop at 90%"))
		f.set(2, "statement", q("governs"))
		f.link(2, op.LinkDecidedBy, 4, 20, "", 0)
		r, asIs := f.whatIf(t, cs), f.asIs(t)
		s := build(t, r, asIs, Request{Change: f.cs(cs), View: "what-if"})
		var focus []string
		for _, fi := range s.Focus {
			focus = append(focus, fi.Handle)
		}
		// The link's endpoints are both touched (D8). governor comes before
		// pause, so con/a keeps governor's direct from, not "from governor".
		eq(t, "focus", focus, []string{"dec/d", "governor", "pause", "resume"})
		eq(t, "constraints", show(s.Constraints), []string{"con/a"})
		eq(t, "decisions", show(s.Decisions), []string{"dec/d"})
		eq(t, "context", show(s.Context), []string{"runtime rel=ancestor", "governor rel=ancestor"})
		if s.Change == nil || s.Change.ID != cs || s.Change.Intent != "Pause at 90%" || s.Change.State != changeset.StateDraft {
			t.Errorf("change header = %+v", s.Change)
		}
		if want := fmt.Sprintf("change %s %q", s.Change.Short, "Pause at 90%"); s.Title != want {
			t.Errorf("title = %q, want %q", s.Title, want)
		}
	})

	t.Run("change slice leaves focus nodes out of neighbours and children", func(t *testing.T) {
		f := tree()
		f.cap(5, "metering", 0)
		f.link(1, op.LinkDependsOn, 3, 5, "", 0)
		cs := f.begin(changeset.StateDraft, "Tangle")
		f.set(2, "statement", q("governs"))
		f.link(2, op.LinkDependsOn, 3, 2, "", 0)
		r := f.whatIf(t, cs)
		s := build(t, r, f.asIs(t), Request{Change: f.cs(cs), View: "what-if"})
		var focus []string
		for _, fi := range s.Focus {
			focus = append(focus, fi.Handle)
		}
		// pause is governor's child and depends on it; both are Focus
		// items, so neither is listed again under the other.
		eq(t, "focus", focus, []string{"governor", "pause"})
		eq(t, "children", show(s.Children), []string{"resume status=none"})
		eq(t, "neighbours", show(s.Neighbours), []string{"metering rel=needs"})
		eq(t, "context", show(s.Context), []string{"runtime rel=ancestor", "governor rel=ancestor"})
	})

	t.Run("change touches link.retire endpoints", func(t *testing.T) {
		f := tree()
		f.link(1, op.LinkDependsOn, 3, 4, "", 0)
		cs := f.begin(changeset.StateAccepted, "Untangle")
		f.retireLink(1)
		r := f.toBe(t)
		s := build(t, r, f.asIs(t), Request{Change: f.cs(cs), View: "to-be"})
		var focus []string
		for _, fi := range s.Focus {
			focus = append(focus, fi.Handle)
		}
		eq(t, "focus", focus, []string{"pause", "resume"})
	})

	t.Run("proposed new", func(t *testing.T) {
		f := tree()
		cs := f.begin(changeset.StateAccepted, "Add audit")
		f.cap(5, "audit", 2)
		s := build(t, f.toBe(t), f.asIs(t), Request{Change: f.cs(cs), View: "to-be"})
		eq(t, "proposed", s.Focus[0].Proposed, []string{"new"})
	})

	t.Run("proposed retired with reason, retired focus kept", func(t *testing.T) {
		f := tree()
		cs := f.begin(changeset.StateAccepted, "Drop resume")
		f.retire(4, "folded into pause")
		s := build(t, f.toBe(t), f.asIs(t), Request{Change: f.cs(cs), View: "to-be"})
		if len(s.Focus) != 1 || s.Focus[0].Handle != "resume" || !s.Focus[0].Retired {
			t.Fatalf("focus = %+v", s.Focus)
		}
		eq(t, "proposed", s.Focus[0].Proposed, []string{"retired: folded into pause"})
	})

	t.Run("proposed moved", func(t *testing.T) {
		f := tree()
		cs := f.begin(changeset.StateDraft, "Promote pause")
		f.move(3, 1)
		f.move(4, 0)
		s := build(t, f.whatIf(t, cs), f.asIs(t), Request{Change: f.cs(cs), View: "what-if"})
		eq(t, "proposed pause", s.Focus[0].Proposed, []string{"moved: governor → runtime"})
		eq(t, "proposed resume", s.Focus[1].Proposed, []string{"moved: governor → (root)"})
	})

	t.Run("proposed fields", func(t *testing.T) {
		f := tree()
		f.set(3, "threshold", `95`, "owner", q("ops"))
		cs := f.begin(changeset.StateDraft, "Pause at 90%")
		f.set(3, "threshold", `90`, "statement", q("stop\n at  90%"), "extra", `{"a":1}`)
		f.add(op.Op{Type: op.TypeNodeUpdate, Node: nd(3), Unset: []string{"owner"}})
		s := build(t, f.whatIf(t, cs), f.asIs(t), Request{Focus: nd(3), Change: f.cs(cs), View: "what-if"})
		eq(t, "proposed", s.Focus[0].Proposed, []string{
			`extra: (none) → {"a":1}`, "owner: ops → (none)", "statement: pause does things → stop at 90%",
			"threshold: 95 → 90"})
		if got := s.Focus[0].Statement; got != "stop at 90%" {
			t.Errorf("statement = %q", got)
		}
		if got := string(s.Focus[0].Fields["threshold"]); got != "90" {
			t.Errorf("fields threshold = %q", got)
		}
	})

	t.Run("as-is has no proposed change", func(t *testing.T) {
		f := tree()
		s := nodeSlice(t, f, 3)
		if s.Focus[0].Proposed != nil {
			t.Errorf("proposed = %q, want nil", s.Focus[0].Proposed)
		}
		if s.Title != "pause" || s.View != "as-is" || s.Change != nil {
			t.Errorf("title, view, change = %q, %q, %+v", s.Title, s.View, s.Change)
		}
	})

	t.Run("statement falls back to the handle", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/bare", 0)
		f.node(11, op.KindConstraint, "con/number", 0, "statement", `7`)
		f.link(1, op.LinkConstrains, 10, 3, "", 0)
		f.link(2, op.LinkConstrains, 11, 3, "", 0)
		s := nodeSlice(t, f, 3)
		if s.Constraints[0].Statement != "con/bare" || s.Constraints[1].Statement != "con/number" {
			t.Errorf("statements = %+v", s.Constraints)
		}
	})

	t.Run("unknown focus is not found", func(t *testing.T) {
		r := tree().asIs(t)
		_, err := Build(r, r, Request{Focus: nd(99), View: "as-is"})
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("err = %v, want ErrNotFound", err)
		}
	})

	t.Run("unknown nothing to slice", func(t *testing.T) {
		r := tree().asIs(t)
		if _, err := Build(r, r, Request{View: "as-is"}); err == nil {
			t.Fatal("no focus and no change: want an error")
		}
	})

	t.Run("retired focus of a node slice is an error", func(t *testing.T) {
		f := tree()
		f.retire(4, "gone")
		r := f.asIs(t)
		if _, err := Build(r, r, Request{Focus: nd(4), View: "as-is"}); !errors.Is(err, ErrRetired) {
			t.Fatalf("err = %v, want ErrRetired", err)
		}
	})

	t.Run("short IDs of change sets are among change sets", func(t *testing.T) {
		mk := func(s string) id.ID {
			i, err := id.Parse(s + strings.Repeat("0", 26-len(s)))
			if err != nil {
				t.Fatal(err)
			}
			return i
		}
		a, b, far := mk("01JB7QAAAA"), mk("01JB7QAAAB"), mk("01K")
		f := tree()
		r := f.asIs(t)
		req := Request{Focus: nd(3), View: "as-is", Change: &changeset.ChangeSet{ID: a, State: changeset.StateDraft},
			InFlight:   []InFlight{{ID: b, State: changeset.StateProposed, Targets: []id.ID{nd(3)}}},
			ChangeSets: []id.ID{a, b, far}}
		s := build(t, r, r, req)
		if s.Change.Short != "01JB7QAAAA" {
			t.Errorf("header short = %q, want 10 characters", s.Change.Short)
		}
		if len(s.InFlight) != 1 || s.InFlight[0].Short != "01JB7QAAAB" {
			t.Errorf("in flight = %+v, want short 01JB7QAAAB", s.InFlight)
		}
		req.ChangeSets = nil
		if s := build(t, r, r, req); s.Change.Short != "01JB7Q" || s.InFlight[0].Short != "01JB7Q" {
			t.Errorf("no change sets: shorts = %q, %q, want 6 characters", s.Change.Short, s.InFlight[0].Short)
		}
	})

	t.Run("short IDs on items", func(t *testing.T) {
		s := nodeSlice(t, tree(), 3)
		r := tree().asIs(t)
		want, _ := r.ShortID(nd(3))
		if s.Focus[0].Short != want || len(want) < 6 {
			t.Errorf("focus short = %q, want %q", s.Focus[0].Short, want)
		}
	})

	t.Run("deterministic", func(t *testing.T) {
		f := tree()
		f.node(10, op.KindConstraint, "con/a", 0)
		// Same handle, another kind: the model allows it, and the sort
		// falls back to the ID.
		f.node(11, op.KindSource, "con/a", 0)
		f.link(1, op.LinkConstrains, 11, 1, "", 0)
		f.link(2, op.LinkConstrains, 10, 1, "", 0)
		a, b := nodeSlice(t, f, 3), nodeSlice(t, f, 3)
		if !reflect.DeepEqual(a, b) {
			t.Error("two builds differ")
		}
		if a.Constraints[0].ID != nd(10) {
			t.Errorf("same handle not sorted by ID: %+v", a.Constraints)
		}
	})
}

func wantEvidence(raw string) string {
	switch raw {
	case `"verified"`:
		return "verified"
	case `"failing"`:
		return "failing"
	}
	return "unverified"
}

func TestShortID(t *testing.T) {
	mk := func(s string) model.Node {
		i, err := id.Parse(s + strings.Repeat("0", 26-len(s)))
		if err != nil {
			t.Fatal(err)
		}
		return model.Node{ID: i}
	}
	ids := []model.Node{mk("01JB2QAAAA"), mk("01JB2QAAAB"), mk("01JB7"), mk("01K")}
	short := func(i id.ID) string {
		s, err := shortID(ids, i)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	cases := []struct {
		id   id.ID
		want string
	}{
		{ids[3].ID, "01K000"},               // min 6
		{ids[2].ID, "01JB70"},               // unique at 5, padded to 6
		{ids[0].ID, "01JB2QAAAA"},           // shares 9 with its neighbour
		{ids[1].ID, "01JB2QAAAB"},           // and the other way
		{mk("01JB2QAAAC").ID, "01JB2QAAAC"}, // not a node: against its would-be neighbours
	}
	for _, c := range cases {
		if got := short(c.id); got != c.want {
			t.Errorf("short(%s) = %q, want %q", c.id, got, c.want)
		}
	}
	if s, _ := shortID(nil, ids[0].ID); s != "01JB2Q" {
		t.Errorf("empty view short = %q", s)
	}
}

func TestFromModel(t *testing.T) {
	f := tree()
	f.link(1, op.LinkDependsOn, 3, 4, "", 0)
	f.link(2, op.LinkRelatesTo, 3, 3, "self", 0)
	r := f.asIs(t)
	anc, err := r.Ancestors(nd(3))
	if err != nil || len(anc) != 2 || anc[0].ID != nd(1) || anc[1].ID != nd(2) {
		t.Errorf("ancestors = %+v, %v", anc, err)
	}
	if _, err := r.Ancestors(nd(99)); !errors.Is(err, ErrNotFound) {
		t.Errorf("ancestors of unknown: %v", err)
	}
	both, _ := r.Links(nd(3), Both)
	if len(both) != 2 || both[0].ID != ln(1) || both[1].ID != ln(2) {
		t.Errorf("both = %+v", both)
	}
	to, _ := r.Links(nd(4), To)
	from, _ := r.Links(nd(4), From)
	if len(to) != 1 || len(from) != 0 {
		t.Errorf("to, from = %d, %d", len(to), len(from))
	}
	kids, _ := r.Children(nd(2))
	if len(kids) != 2 {
		t.Errorf("children = %d", len(kids))
	}
}
