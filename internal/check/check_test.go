package check

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/fs"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/render"
	"github.com/dphaener/rhiza/internal/store"
)

const repoJSON = "{\n  \"format\": 1,\n  \"id\": \"01M5N90F40GQXYEAV0CGM90155\",\n  \"remotes\": []\n}\n"

// fid is fixture ID n; IDs grow with n.
func fid(n int) id.ID {
	var i id.ID
	i[0] = 1
	binary.BigEndian.PutUint32(i[12:], uint32(n))
	return i
}

var actor = fid(9999)

// repo builds a .rhiza folder in memory, entry by entry, through
// entry.Encode, so tests read it with store.ReadLog as check's caller will.
type repo struct {
	t     *testing.T
	n     int
	files fstest.MapFS
}

func newRepo(t *testing.T) *repo {
	return &repo{t: t, n: 100, files: fstest.MapFS{store.RepoFile: {Data: []byte(repoJSON)}}}
}

func (r *repo) next() id.ID {
	r.n++
	return fid(r.n)
}

// set is one change set being written.
type set struct {
	r    *repo
	id   id.ID
	last id.ID
	// by is the actor of the set's entries; zero means actor.
	by id.ID
}

func (r *repo) open(kind changeset.Kind, intent string) *set {
	s := &set{r: r, id: r.next()}
	s.add(changeset.Open{Kind: kind, Intent: intent}, nil)
	return s
}

// record opens a record change set with ops and lands it.
func (r *repo) record(intent string, ops ...op.Op) *set {
	s := r.open(changeset.KindRecord, intent)
	s.then(changeset.Ops{Ops: ops})
	s.then(changeset.Land{})
	return s
}

// recordBy is record with every entry written by actor who.
func (r *repo) recordBy(who id.ID, intent string, ops ...op.Op) *set {
	s := &set{r: r, id: r.next(), by: who}
	s.add(changeset.Open{Kind: changeset.KindRecord, Intent: intent}, nil)
	s.then(changeset.Ops{Ops: ops})
	s.then(changeset.Land{})
	return s
}

// then adds body after the last entry added.
func (s *set) then(b changeset.Body) id.ID { return s.add(b, []id.ID{s.last}) }

// add adds body with the given after list and returns the entry's ID.
func (s *set) add(b changeset.Body, after []id.ID) id.ID {
	s.r.t.Helper()
	body, err := changeset.EncodeBody(b)
	if err != nil {
		s.r.t.Fatal(err)
	}
	who := actor
	if !s.by.IsZero() {
		who = s.by
	}
	e := entry.Entry{V: changeset.EntryVersion(b), ID: s.r.next(), ChangeSet: s.id, Actor: who, For: who,
		At: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC), After: after, Type: b.Type(), Body: body}
	data, err := entry.Encode(e)
	if err != nil {
		s.r.t.Fatal(err)
	}
	s.r.files[entryPath(s.id, e.ID)] = &fstest.MapFile{Data: data}
	s.last = e.ID
	return e.ID
}

func entryPath(cs, e id.ID) string {
	return fmt.Sprintf("%s/%s/%s.json", store.LogDir, cs, e)
}

func statement(s string) map[string]json.RawMessage {
	b, _ := json.Marshal(s)
	return map[string]json.RawMessage{"statement": b}
}

func mkNode(o, node id.ID, k op.Kind, handle string) op.Op {
	x := op.Op{ID: o, Type: op.TypeNodeCreate, Node: node, Kind: k, Handle: handle, Fields: statement(handle)}
	if op.HasParent(k) {
		x.Key = "i"
	}
	return x
}

func mkLink(o, link id.ID, t op.LinkType, from, to id.ID) op.Op {
	return op.Op{ID: o, Type: op.TypeLinkCreate, Link: link, LinkType: t, From: from, To: to}
}

func setField(o, node, base id.ID, value string) op.Op {
	return op.Op{ID: o, Type: op.TypeNodeUpdate, Node: node, Base: base, Set: statement(value)}
}

// rendered is what render writes for log l.
func rendered(l store.Log) fstest.MapFS {
	sets, _ := changeset.FoldLog(l)
	m, _ := model.AsIs(sets)
	out := fstest.MapFS{}
	for p, b := range render.Render(m, sets) {
		out[p] = &fstest.MapFile{Data: b}
	}
	return out
}

// input reads r's log and pairs it with a current render.
func (r *repo) input() Input {
	l := store.ReadLog(r.files)
	return Input{Log: l, Rendered: rendered(l)}
}

// clean writes a landed intent, a capability serving it and a serves link
// from the change set itself; it checks clean.
func clean(r *repo) (intent, capab id.ID) {
	intent, capab, _ = cleanOps(r)
	return intent, capab
}

// cleanOps is clean, also returning the capability's create op.
func cleanOps(r *repo) (intent, capab, capOp id.ID) {
	intent, capab = r.next(), r.next()
	o1, o2, o3, l1 := r.next(), r.next(), r.next(), r.next()
	cs := r.open(changeset.KindRecord, "Start")
	cs.then(changeset.Ops{Ops: []op.Op{
		mkNode(o1, intent, op.KindIntent, "goal"),
		mkNode(o2, capab, op.KindCapability, "thing"),
		mkLink(o3, l1, op.LinkServes, capab, intent),
	}})
	o4, l2 := r.next(), r.next()
	cs.then(changeset.Ops{Ops: []op.Op{mkLink(o4, l2, op.LinkServes, cs.id, intent)}})
	cs.then(changeset.Land{})
	return intent, capab, o2
}

// short is a finding as severity, rule and path, for comparing.
func short(fs []Finding) []string {
	out := []string{}
	for _, f := range fs {
		out = append(out, strings.TrimSpace(f.Severity+" "+f.Rule+" "+f.Path))
	}
	return out
}

func TestRun(t *testing.T) {
	cases := []struct {
		name string
		// build returns the input and the expected findings in order.
		build func(t *testing.T) (Input, []string, func(t *testing.T, r Report))
	}{
		{"clean", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			return r.input(), nil, nil
		}},
		{"empty", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			r.files[store.LogDir] = &fstest.MapFile{Mode: fs.ModeDir | 0o755}
			return r.input(), nil, nil
		}},
		{"log", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			r.files["log/junk.txt"] = &fstest.MapFile{Data: []byte("x")}
			return r.input(), []string{"fail log .rhiza/log/junk.txt"}, func(t *testing.T, rep Report) {
				if !strings.Contains(rep.Findings[0].Message, "unexpected") {
					t.Errorf("message %q", rep.Findings[0].Message)
				}
			}
		}},
		{"fold", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			s := r.open(changeset.KindIntention, "Bad move")
			bad := s.then(changeset.Accept{Rule: "owner"})
			return r.input(), []string{"fail fold " + ".rhiza/" + entryPath(s.id, bad)}, func(t *testing.T, rep Report) {
				f := rep.Findings[0]
				if f.ChangeSet != s.id.String() || f.Entry != bad.String() || !strings.Contains(f.Message, "illegal lifecycle move") {
					t.Errorf("finding %+v", f)
				}
			}
		}},
		{"concurrent-warn", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			s := r.open(changeset.KindIntention, "Two writers")
			start := s.last
			o, n := r.next(), r.next()
			s.add(changeset.Ops{Ops: []op.Op{mkNode(o, n, op.KindDecision, "a")}}, []id.ID{start})
			later := s.add(changeset.Ops{Ops: []op.Op{mkNode(o, n, op.KindDecision, "b")}}, []id.ID{start})
			return r.input(), []string{"warn fold .rhiza/" + entryPath(s.id, later)}, nil
		}},
		{"model", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			n, o1 := r.next(), r.next()
			r.record("Create", mkNode(o1, n, op.KindDecision, "pick"))
			o2 := r.next()
			r.record("Edit one", setField(o2, n, o1, "one"))
			o3 := r.next()
			c := r.record("Edit two", setField(o3, n, o1, "two"))
			return r.input(), []string{"fail model"}, func(t *testing.T, rep Report) {
				f := rep.Findings[0]
				if f.ChangeSet != c.id.String() || f.Op != o3.String() || !strings.Contains(f.Message, "conflict") {
					t.Errorf("finding %+v", f)
				}
				want := fmt.Sprintf(`rhiza change new "Resolve the conflict on %s"`, n)
				if f.Next != want {
					t.Errorf("next %q, want %q", f.Next, want)
				}
			}
		}},
		{"in-flight", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			s := r.open(changeset.KindIntention, "Edit a ghost")
			o := r.next()
			s.then(changeset.Ops{Ops: []op.Op{setField(o, r.next(), r.next(), "x")}})
			s.then(changeset.Propose{})
			s.then(changeset.Accept{Rule: "owner"})
			// Accepted: the problem is in ToBe and in WhatIf, reported once.
			return r.input(), []string{"warn in-flight"}, func(t *testing.T, rep Report) {
				f := rep.Findings[0]
				if f.ChangeSet != s.id.String() || f.Op != o.String() || !strings.Contains(f.Message, "not found") {
					t.Errorf("finding %+v", f)
				}
			}
		}},
		{"link", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			intent, capab := clean(r)
			o, l := r.next(), r.next()
			s := r.record("Bad link", mkLink(o, l, op.LinkVerifies, intent, capab))
			return r.input(), []string{"fail link"}, func(t *testing.T, rep Report) {
				f := rep.Findings[0]
				if f.Link != l.String() || f.ChangeSet != s.id.String() || f.Op != o.String() {
					t.Errorf("finding %+v", f)
				}
				if !strings.Contains(f.Message, `"verifies"`) || !strings.Contains(f.Message, "intent") ||
					!strings.Contains(f.Message, "capability") {
					t.Errorf("message %q", f.Message)
				}
			}
		}},
		{"link-retired", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			intent, capab, capOp := cleanOps(r)
			o, l := r.next(), r.next()
			r.record("Bad link", mkLink(o, l, op.LinkVerifies, intent, capab))
			o2 := r.next()
			r.record("Drop it", op.Op{ID: o2, Type: op.TypeLinkRetire, Link: l, Base: o})
			// A live link to a retired node is fine: only kinds matter.
			ev, o3, o4, l2 := r.next(), r.next(), r.next(), r.next()
			r.record("Evidence", mkNode(o3, ev, op.KindEvidence, "proof"), mkLink(o4, l2, op.LinkVerifies, ev, capab))
			o5 := r.next()
			r.record("Retire it", op.Op{ID: o5, Type: op.TypeNodeRetire, Node: capab, Base: capOp, Reason: "gone"})
			return r.input(), nil, nil
		}},
		{"render-differ", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			in := r.input()
			in.Rendered.(fstest.MapFS)["index.md"] = &fstest.MapFile{Data: []byte("# Old\n")}
			return in, []string{"fail render .rhiza/rendered/index.md"}, func(t *testing.T, rep Report) {
				if rep.Findings[0].Next != "rhiza render" {
					t.Errorf("next %q", rep.Findings[0].Next)
				}
			}
		}},
		{"render-missing", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			in := r.input()
			delete(in.Rendered.(fstest.MapFS), "capability/thing.md")
			return in, []string{"fail render .rhiza/rendered/capability/thing.md"}, nil
		}},
		{"render-extra", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			in := r.input()
			in.Rendered.(fstest.MapFS)[".DS_Store"] = &fstest.MapFile{Data: []byte("x")}
			in.Rendered.(fstest.MapFS)["old/page.md"] = &fstest.MapFile{Data: []byte("x")}
			return in, []string{"fail render .rhiza/rendered/.DS_Store", "fail render .rhiza/rendered/old/page.md"}, nil
		}},
		{"render-nil", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			in := r.input()
			in.Rendered = nil
			return in, []string{"fail render .rhiza/rendered"}, func(t *testing.T, rep Report) {
				if rep.Findings[0].Next != "rhiza render" {
					t.Errorf("next %q", rep.Findings[0].Next)
				}
			}
		}},
		{"append-only", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			in := r.input()
			in.Rev = "abc123"
			in.Base = []Entry{
				{ID: "d", Path: ".rhiza/log/a/d.json", Blob: "b1"},
				{ID: "m", Path: ".rhiza/log/b/m.json", Blob: "b2"},
				{ID: "k", Path: ".rhiza/log/b/k.json", Blob: "b3"},
			}
			in.Current = []Entry{
				{ID: "m", Path: ".rhiza/log/b/m.json", Blob: "b2-edited"},
				{ID: "k", Path: ".rhiza/log/b/k.json", Blob: "b3"},
				{ID: "a", Path: ".rhiza/log/b/a.json", Blob: "b4"},
			}
			return in, []string{"fail append-only .rhiza/log/a/d.json", "fail append-only .rhiza/log/b/m.json"},
				func(t *testing.T, rep Report) {
					if rep.Base != "abc123" {
						t.Errorf("base %q", rep.Base)
					}
					if want := "git checkout abc123 -- .rhiza/log/a/d.json"; rep.Findings[0].Next != want {
						t.Errorf("next %q, want %q", rep.Findings[0].Next, want)
					}
				}
		}},
		{"skipped", func(t *testing.T) (Input, []string, func(*testing.T, Report)) {
			r := newRepo(t)
			clean(r)
			in := r.input()
			in.Skipped = "not a git work tree"
			return in, []string{"warn append-only"}, func(t *testing.T, rep Report) {
				if !strings.Contains(rep.Findings[0].Message, "not a git work tree") {
					t.Errorf("message %q", rep.Findings[0].Message)
				}
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in, want, more := c.build(t)
			rep := Run(in)
			if want == nil {
				want = []string{}
			}
			if got := short(rep.Findings); !reflect.DeepEqual(got, want) {
				t.Fatalf("findings\n got %q\nwant %q\nall %+v", got, want, rep.Findings)
			}
			wantOK := !slices.ContainsFunc(want, func(s string) bool { return strings.HasPrefix(s, SeverityFail) })
			if rep.OK != wantOK {
				t.Errorf("OK = %v, want %v", rep.OK, wantOK)
			}
			if more != nil {
				more(t, rep)
			}
			if again := Run(in); !reflect.DeepEqual(again, rep) {
				t.Errorf("second run differs")
			}
		})
	}
}

// TestRunOrder puts one finding of each rule in a single input and checks
// they come in rule order.
func TestRunOrder(t *testing.T) {
	r := newRepo(t)
	intent, capab := clean(r)
	r.files["log/junk.txt"] = &fstest.MapFile{Data: []byte("x")}
	bad := r.open(changeset.KindIntention, "Bad move")
	bad.then(changeset.Accept{Rule: "owner"})
	o, l := r.next(), r.next()
	r.record("Bad link", mkLink(o, l, op.LinkVerifies, intent, capab))
	n, o1 := r.next(), r.next()
	r.record("Create", mkNode(o1, n, op.KindDecision, "pick"))
	o2, o3 := r.next(), r.next()
	r.record("Edit one", setField(o2, n, o1, "one"))
	r.record("Edit two", setField(o3, n, o1, "two"))
	g := r.open(changeset.KindIntention, "Ghost")
	g.then(changeset.Ops{Ops: []op.Op{setField(r.next(), r.next(), r.next(), "x")}})
	la := r.open(changeset.KindIntention, "Loop A")
	lb := r.open(changeset.KindIntention, "Loop B")
	r.record("Loop", dependsOn(r, la.id, lb.id), dependsOn(r, lb.id, la.id))
	leaf, ev := r.next(), r.next()
	leafSet := r.record("Leaf", mkNode(r.next(), leaf, op.KindCapability, "leaf"))
	r.record("Verify", mkNode(r.next(), ev, op.KindEvidence, "seen"), mkLink(r.next(), r.next(), op.LinkVerifies, ev, leaf))
	in := r.input()
	in.Rendered = nil
	in.Skipped = "no git"
	in.Commits = []Commit{{Hash: "abc", Changes: []string{leafSet.id.String()}, Actors: []string{actor.String()}}}

	rep := Run(in)
	var rules []string
	for _, f := range rep.Findings {
		rules = append(rules, f.Rule)
	}
	rules = slices.Compact(rules)
	want := []string{RuleLog, RuleFold, RuleModel, RuleInFlight, RuleLink, RuleDependsOn, RuleEvidenceActor,
		RuleRender, RuleAppendOnly}
	if !reflect.DeepEqual(rules, want) {
		t.Fatalf("rules %q, want %q\n%+v", rules, want, rep.Findings)
	}
	if !slices.Equal(Rules(), want) {
		t.Errorf("Rules() = %q", Rules())
	}
	if rep.OK {
		t.Error("OK with failures")
	}
}

func TestRunZeroInput(t *testing.T) {
	rep := Run(Input{})
	if rep.OK || len(rep.Findings) != 1 || rep.Findings[0].Rule != RuleRender {
		t.Errorf("Run(Input{}) = %+v", rep)
	}
}

// dependsOn is a depends-on link from change set from to change set to.
func dependsOn(r *repo, from, to id.ID) op.Op {
	return mkLink(r.next(), r.next(), op.LinkDependsOn, from, to)
}

// ofRule is the findings of rule in rep.
func ofRule(rep Report, rule string) []Finding {
	var out []Finding
	for _, f := range rep.Findings {
		if f.Rule == rule {
			out = append(out, f)
		}
	}
	return out
}

func TestRunCycle(t *testing.T) {
	t.Run("CycleTwo", func(t *testing.T) {
		r := newRepo(t)
		clean(r)
		a := r.open(changeset.KindIntention, "A")
		b := r.open(changeset.KindIntention, "B")
		a.then(changeset.Ops{Ops: []op.Op{dependsOn(r, a.id, b.id)}})
		b.then(changeset.Ops{Ops: []op.Op{dependsOn(r, b.id, a.id)}})
		rep := Run(r.input())
		fs := ofRule(rep, RuleDependsOn)
		if len(fs) != 1 {
			t.Fatalf("depends-on findings %+v\nall %+v", fs, rep.Findings)
		}
		f := fs[0]
		if f.Severity != SeverityFail || f.ChangeSet != a.id.String() ||
			!strings.Contains(f.Message, a.id.String()) || !strings.Contains(f.Message, b.id.String()) {
			t.Errorf("finding %+v", f)
		}
		if want := "rhiza change show " + a.id.String(); f.Next != want {
			t.Errorf("next %q, want %q", f.Next, want)
		}
		if rep.OK {
			t.Error("OK with a loop")
		}
	})
	t.Run("CycleLanded", func(t *testing.T) {
		r := newRepo(t)
		clean(r)
		a := r.open(changeset.KindIntention, "A")
		b := r.open(changeset.KindIntention, "B")
		r.record("Links", dependsOn(r, a.id, b.id), dependsOn(r, b.id, a.id))
		if fs := ofRule(Run(r.input()), RuleDependsOn); len(fs) != 1 || fs[0].Severity != SeverityFail {
			t.Errorf("depends-on findings %+v", fs)
		}
	})
	t.Run("CycleSelf", func(t *testing.T) {
		// A log cannot hold a link from a change set to itself (op
		// validation refuses it), so the change set is built folded.
		a := fid(1)
		var c checker
		c.dependsOn([]changeset.ChangeSet{{ID: a, State: changeset.StateDraft,
			Ops: []op.Op{mkLink(fid(2), fid(3), op.LinkDependsOn, a, a)}}})
		if len(c.findings) != 1 || c.findings[0].Severity != SeverityFail || c.findings[0].ChangeSet != a.String() ||
			!strings.Contains(c.findings[0].Message, "itself") {
			t.Fatalf("findings %+v", c.findings)
		}
	})
	t.Run("CycleNone", func(t *testing.T) {
		r := newRepo(t)
		clean(r)
		a := r.open(changeset.KindIntention, "A")
		b := r.open(changeset.KindIntention, "B")
		b.then(changeset.Ops{Ops: []op.Op{dependsOn(r, b.id, a.id)}})
		if fs := ofRule(Run(r.input()), RuleDependsOn); len(fs) != 0 {
			t.Errorf("depends-on findings %+v", fs)
		}
	})
}

func TestRunRestsOn(t *testing.T) {
	// build has A draft a node and B update it, then ends A with end; it
	// returns B, B's update and A.
	build := func(t *testing.T, end func(r *repo, a *set)) (*repo, *set, id.ID, *set) {
		r := newRepo(t)
		clean(r)
		a := r.open(changeset.KindIntention, "A")
		n, a1 := r.next(), r.next()
		a.then(changeset.Ops{Ops: []op.Op{mkNode(a1, n, op.KindDecision, "pick")}})
		b := r.open(changeset.KindIntention, "B")
		b1 := r.next()
		b.then(changeset.Ops{Ops: []op.Op{dependsOn(r, b.id, a.id), setField(b1, n, a1, "two")}})
		if end != nil {
			end(r, a)
		}
		return r, b, b1, a
	}
	// resting is the in-flight findings that name a resting op.
	resting := func(rep Report) []Finding {
		var out []Finding
		for _, f := range ofRule(rep, RuleInFlight) {
			if strings.Contains(f.Message, "rests on") {
				out = append(out, f)
			}
		}
		return out
	}
	for _, c := range []struct {
		name  string
		state changeset.State
		end   func(r *repo, a *set)
	}{
		{"DroppedWithdrawn", changeset.StateWithdrawn, func(r *repo, a *set) { a.then(changeset.Withdraw{}) }},
		{"DroppedRejected", changeset.StateRejected, func(r *repo, a *set) {
			a.then(changeset.Propose{})
			a.then(changeset.Reject{Rule: "owner", Reason: "no"})
		}},
		{"DroppedSuperseded", changeset.StateSuperseded, func(r *repo, a *set) {
			a.then(changeset.Supersede{By: r.next()})
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			r, b, b1, a := build(t, c.end)
			rep := Run(r.input())
			fs := resting(rep)
			if len(fs) != 1 {
				t.Fatalf("resting findings %+v\nall %+v", fs, rep.Findings)
			}
			f := fs[0]
			if f.Severity != SeverityWarn || f.ChangeSet != b.id.String() || f.Op != b1.String() ||
				!strings.Contains(f.Message, a.id.String()) || !strings.Contains(f.Message, string(c.state)) {
				t.Errorf("finding %+v", f)
			}
			if want := "rhiza change show " + b.id.String(); f.Next != want {
				t.Errorf("next %q, want %q", f.Next, want)
			}
		})
	}
	t.Run("RestsOnLanded", func(t *testing.T) {
		r, _, _, _ := build(t, func(r *repo, a *set) {
			a.then(changeset.Propose{})
			a.then(changeset.Accept{Rule: "owner"})
			a.then(changeset.Land{})
		})
		if fs := resting(Run(r.input())); len(fs) != 0 {
			t.Errorf("resting findings %+v", fs)
		}
	})
	t.Run("RestsOnInFlight", func(t *testing.T) {
		r, _, _, _ := build(t, nil)
		if fs := resting(Run(r.input())); len(fs) != 0 {
			t.Errorf("resting findings %+v", fs)
		}
	})
	t.Run("RestsOnRevised", func(t *testing.T) {
		// A revises its create in place after B's update wrote the same
		// field: a revised conflict, warned with a next.
		r := newRepo(t)
		clean(r)
		a := r.open(changeset.KindIntention, "A")
		n, a1 := r.next(), r.next()
		a.then(changeset.Ops{Ops: []op.Op{mkNode(a1, n, op.KindDecision, "pick")}})
		b := r.open(changeset.KindIntention, "B")
		b1 := r.next()
		b.then(changeset.Ops{Ops: []op.Op{dependsOn(r, b.id, a.id), setField(b1, n, a1, "two")}})
		revised := mkNode(a1, n, op.KindDecision, "pick")
		revised.Fields = statement("changed")
		a.then(changeset.Ops{Ops: []op.Op{revised}})
		rep := Run(r.input())
		var fs []Finding
		for _, f := range ofRule(rep, RuleInFlight) {
			if strings.Contains(f.Message, "revised") {
				fs = append(fs, f)
			}
		}
		if len(fs) != 1 {
			t.Fatalf("revised findings %+v\nall %+v", fs, rep.Findings)
		}
		f := fs[0]
		if f.Severity != SeverityWarn || f.ChangeSet != b.id.String() || f.Op != b1.String() {
			t.Errorf("finding %+v", f)
		}
		if want := fmt.Sprintf("rhiza slice --change %s --view what-if", b.id); f.Next != want {
			t.Errorf("next %q, want %q", f.Next, want)
		}
		if !rep.OK {
			t.Errorf("not OK on warnings: %+v", rep.Findings)
		}
	})
}

// appendOnly runs the append-only rule alone on base and current entries at
// base revision "abc".
func appendOnly(base, current []Entry) []Finding {
	c := checker{in: Input{Rev: "abc", Base: base, Current: current}}
	c.appendOnly()
	return c.findings
}

func TestAppendOnlyIdentity(t *testing.T) {
	const (
		old = ".old/log/cs1/e1.json"
		at  = ".rhiza/log/cs1/e1.json"
	)
	for _, c := range []struct {
		name          string
		base, current []Entry
		want          []Finding
	}{
		{name: "Unchanged",
			base:    []Entry{{ID: "e1", Path: at, Blob: "x"}},
			current: []Entry{{ID: "e1", Path: at, Blob: "x"}}},
		{name: "MovedFolderPasses",
			base:    []Entry{{ID: "e1", Path: old, Blob: "x"}, {ID: "e2", Path: ".old/log/cs1/e2.json", Blob: "y"}},
			current: []Entry{{ID: "e1", Path: at, Blob: "x"}, {ID: "e2", Path: ".rhiza/log/cs1/e2.json", Blob: "y"}}},
		{name: "OtherChangeSetFolderPasses",
			base:    []Entry{{ID: "e1", Path: at, Blob: "x"}},
			current: []Entry{{ID: "e1", Path: ".rhiza/log/cs2/e1.json", Blob: "x"}}},
		{name: "AddedPasses",
			base:    []Entry{{ID: "e1", Path: at, Blob: "x"}},
			current: []Entry{{ID: "e1", Path: at, Blob: "x"}, {ID: "e9", Path: ".rhiza/log/cs1/e9.json", Blob: "z"}}},
		{name: "OneCopyWithSameBytesPasses",
			base:    []Entry{{ID: "e1", Path: at, Blob: "x"}},
			current: []Entry{{ID: "e1", Path: at, Blob: "edited"}, {ID: "e1", Path: ".rhiza/log/cs2/e1.json", Blob: "x"}}},
		{name: "EditedFails",
			base:    []Entry{{ID: "e1", Path: at, Blob: "x"}},
			current: []Entry{{ID: "e1", Path: at, Blob: "edited"}},
			want: []Finding{{Severity: SeverityFail, Rule: RuleAppendOnly, Path: at, Entry: "e1",
				Message: "log entry e1 modified since the base; entries may only be added",
				Next:    "git checkout abc -- " + at}}},
		{name: "MovedAndEditedFails",
			base:    []Entry{{ID: "e1", Path: old, Blob: "x"}},
			current: []Entry{{ID: "e1", Path: at, Blob: "edited"}},
			want: []Finding{{Severity: SeverityFail, Rule: RuleAppendOnly, Path: at, Entry: "e1",
				Message: "log entry e1 modified since the base, where it was " + old + "; entries may only be added",
				Next:    "git checkout abc -- " + old + " && git mv -f " + old + " " + at}}},
		{name: "MovedWithinStoreAndEditedFails",
			base:    []Entry{{ID: "e1", Path: at, Blob: "x"}},
			current: []Entry{{ID: "e1", Path: ".rhiza/log/cs2/e1.json", Blob: "edited"}},
			want: []Finding{{Severity: SeverityFail, Rule: RuleAppendOnly, Path: ".rhiza/log/cs2/e1.json", Entry: "e1",
				Message: "log entry e1 modified since the base, where it was " + at + "; entries may only be added",
				Next:    "git checkout abc -- " + at}}},
		{name: "EditedPrefersBasePath",
			base: []Entry{{ID: "e1", Path: ".rhiza/log/cs2/e1.json", Blob: "x"}},
			current: []Entry{{ID: "e1", Path: ".rhiza/log/cs1/e1.json", Blob: "p"},
				{ID: "e1", Path: ".rhiza/log/cs2/e1.json", Blob: "q"}},
			want: []Finding{{Severity: SeverityFail, Rule: RuleAppendOnly, Path: ".rhiza/log/cs2/e1.json", Entry: "e1",
				Message: "log entry e1 modified since the base; entries may only be added",
				Next:    "git checkout abc -- .rhiza/log/cs2/e1.json"}}},
		{name: "DeletedFails",
			base: []Entry{{ID: "e1", Path: at, Blob: "x"}},
			want: []Finding{{Severity: SeverityFail, Rule: RuleAppendOnly, Path: at, Entry: "e1",
				Message: "log entry e1 deleted since the base; entries may only be added",
				Next:    "git checkout abc -- " + at}}},
		{name: "DeletedFromOldFolderFails",
			base:    []Entry{{ID: "e1", Path: old, Blob: "x"}},
			current: []Entry{{ID: "e2", Path: ".rhiza/log/cs1/e2.json", Blob: "x"}},
			want: []Finding{{Severity: SeverityFail, Rule: RuleAppendOnly, Path: old, Entry: "e1",
				Message: "log entry e1 deleted since the base; entries may only be added",
				Next:    "git checkout abc -- " + old + " && git mv " + old + " " + at}}},
		{name: "SameBaseEntryTwiceOnce",
			base: []Entry{{ID: "e1", Path: at, Blob: "x"}, {ID: "e1", Path: at, Blob: "x"}},
			want: []Finding{{Severity: SeverityFail, Rule: RuleAppendOnly, Path: at, Entry: "e1",
				Message: "log entry e1 deleted since the base; entries may only be added",
				Next:    "git checkout abc -- " + at}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := appendOnly(c.base, c.current)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("findings\n got %+v\nwant %+v", got, c.want)
			}
		})
	}
}

// TestAppendOnlyNoRev reports a deleted entry without a Next when there is
// no base revision to restore it from.
func TestAppendOnlyNoRev(t *testing.T) {
	rep := Run(Input{Base: []Entry{{ID: "e1", Path: ".rhiza/log/cs/e1.json", Blob: "x"}}})
	f := ofRule(rep, RuleAppendOnly)
	if len(f) != 1 || f[0].Next != "" || rep.OK {
		t.Errorf("findings %+v", rep.Findings)
	}
}

// TestLoops finds a graph's loops: groups that reach one another and
// self-loops, not chains.
func TestLoops(t *testing.T) {
	a, b, c, d := fid(1), fid(2), fid(3), fid(4)
	got := loops(map[id.ID][]id.ID{a: {b}, b: {c, a}, c: {d}, d: {d}})
	if len(got) != 2 || len(got[0]) != 1 || got[0][0] != d || len(got[1]) != 2 {
		t.Errorf("loops = %v", got)
	}
}

func TestRunCapabilityLoop(t *testing.T) {
	// caps lands capabilities a, b and c and returns their IDs.
	caps := func(r *repo) (a, b, c id.ID) {
		a, b, c = r.next(), r.next(), r.next()
		r.record("Caps", mkNode(r.next(), a, op.KindCapability, "a"),
			mkNode(r.next(), b, op.KindCapability, "b"), mkNode(r.next(), c, op.KindCapability, "c"))
		return a, b, c
	}
	dep := func(r *repo, from, to id.ID) op.Op {
		return mkLink(r.next(), r.next(), op.LinkDependsOn, from, to)
	}

	t.Run("Concurrent", func(t *testing.T) {
		// Each change set alone is a chain; landed together they loop.
		r := newRepo(t)
		clean(r)
		a, b, c := caps(r)
		r.record("Chain", dep(r, a, b), dep(r, b, c))
		r.record("Back", dep(r, c, a))
		rep := Run(r.input())
		fs := ofRule(rep, RuleDependsOn)
		if len(fs) != 1 {
			t.Fatalf("depends-on findings %+v\nall %+v", fs, rep.Findings)
		}
		want := Finding{Severity: SeverityFail, Rule: RuleDependsOn,
			Message: "depends-on loop among capabilities a, b, c", Next: "rhiza show a"}
		if fs[0] != want {
			t.Errorf("finding %+v, want %+v", fs[0], want)
		}
		if rep.OK {
			t.Error("OK with a loop")
		}
	})
	t.Run("Retired", func(t *testing.T) {
		r := newRepo(t)
		clean(r)
		a, b, _ := caps(r)
		ab, l := r.next(), r.next()
		r.record("Chain", mkLink(ab, l, op.LinkDependsOn, a, b), dep(r, b, a))
		r.record("Unlink", op.Op{ID: r.next(), Type: op.TypeLinkRetire, Link: l, Base: ab})
		if fs := ofRule(Run(r.input()), RuleDependsOn); len(fs) != 0 {
			t.Errorf("depends-on findings %+v", fs)
		}
	})
	t.Run("Chain", func(t *testing.T) {
		r := newRepo(t)
		clean(r)
		a, b, c := caps(r)
		r.record("Chain", dep(r, a, b), dep(r, b, c))
		r.record("Short", dep(r, a, c))
		rep := Run(r.input())
		if fs := ofRule(rep, RuleDependsOn); len(fs) != 0 || !rep.OK {
			t.Errorf("depends-on findings %+v\nall %+v", fs, rep.Findings)
		}
	})
}

// The actors of the evidence-actor fixture, each recorded as an actor node
// whose ID is the actor ID its entries carry.
var (
	builder  = fid(9001)
	verifier = fid(9002)
)

// evidenceRepo lands actors builder and verifier, capabilities leaf and
// other in change set cs, then Evidence ev verifying leaf over link l,
// written by who. ev's create op is evOp and l's evLinkOp.
type evidenceRepo struct {
	r                                   *repo
	cs, leaf, ev, l, evOp, evLinkOp, by id.ID
}

func newEvidenceRepo(t *testing.T, who id.ID) evidenceRepo {
	r := newRepo(t)
	clean(r)
	r.record("Actors", mkNode(r.next(), builder, op.KindActor, "builder"),
		mkNode(r.next(), verifier, op.KindActor, "verifier"))
	x := evidenceRepo{r: r, leaf: r.next(), by: who}
	x.cs = r.record("Leaf", mkNode(r.next(), x.leaf, op.KindCapability, "leaf"),
		mkNode(r.next(), r.next(), op.KindCapability, "other")).id
	x.ev, x.l, x.evOp, x.evLinkOp = r.next(), r.next(), r.next(), r.next()
	r.recordBy(who, "Verify", mkNode(x.evOp, x.ev, op.KindEvidence, "leaf-verified"),
		mkLink(x.evLinkOp, x.l, op.LinkVerifies, x.ev, x.leaf))
	return x
}

func TestRunEvidenceActor(t *testing.T) {
	const hash = "0123456789abcdef0123456789abcdef01234567"
	commit := func(cs id.ID, actors, leaves []string) Commit {
		return Commit{Hash: hash, Changes: []string{cs.String()}, Actors: actors, Leaves: leaves}
	}
	t.Run("Builder", func(t *testing.T) {
		x := newEvidenceRepo(t, builder)
		in := x.r.input()
		in.Commits = []Commit{commit(x.cs, []string{"builder"}, nil)}
		rep := Run(in)
		fs := ofRule(rep, RuleEvidenceActor)
		want := Finding{Severity: SeverityWarn, Rule: RuleEvidenceActor, ChangeSet: x.cs.String(), Link: x.l.String(),
			Message: "evidence leaf-verified verifies capability leaf, but its actor builder wrote that " +
				"capability in commit 0123456789ab",
			Next: "rhiza show leaf-verified"}
		if len(fs) != 1 || fs[0] != want {
			t.Fatalf("evidence-actor findings %+v, want %+v\nall %+v", fs, want, rep.Findings)
		}
		if !rep.OK {
			t.Errorf("not OK on a warning: %+v", rep.Findings)
		}
	})
	cases := []struct {
		name    string
		who     id.ID
		commits func(x evidenceRepo) []Commit
		want    int
	}{
		{"OtherActor", verifier, func(x evidenceRepo) []Commit {
			return []Commit{commit(x.cs, []string{"builder"}, nil)}
		}, 0},
		{"NoActor", builder, func(x evidenceRepo) []Commit { return []Commit{commit(x.cs, nil, nil)} }, 0},
		{"OtherChange", builder, func(x evidenceRepo) []Commit {
			return []Commit{commit(x.ev, []string{"builder"}, nil)}
		}, 0},
		{"OtherLeaf", builder, func(x evidenceRepo) []Commit {
			return []Commit{commit(x.cs, []string{"builder"}, []string{"other"})}
		}, 0},
		{"ThisLeaf", builder, func(x evidenceRepo) []Commit {
			return []Commit{commit(x.cs, []string{"builder"}, []string{"other", "leaf"})}
		}, 1},
		{"ThisLeafByID", builder, func(x evidenceRepo) []Commit {
			return []Commit{commit(x.cs, []string{"builder"}, []string{x.leaf.String()})}
		}, 1},
		{"ActorByID", builder, func(x evidenceRepo) []Commit {
			return []Commit{commit(x.cs, []string{builder.String()}, nil)}
		}, 1},
		{"ManyCommits", builder, func(x evidenceRepo) []Commit {
			return []Commit{commit(x.cs, []string{"builder"}, nil), commit(x.cs, []string{"builder"}, nil)}
		}, 1},
		{"NoGit", builder, func(evidenceRepo) []Commit { return nil }, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			x := newEvidenceRepo(t, tc.who)
			in := x.r.input()
			in.Commits = tc.commits(x)
			rep := Run(in)
			if fs := ofRule(rep, RuleEvidenceActor); len(fs) != tc.want || !rep.OK {
				t.Errorf("evidence-actor findings %+v, want %d\nall %+v", fs, tc.want, rep.Findings)
			}
		})
	}
	t.Run("Oldest", func(t *testing.T) {
		// Commits come newest first; the finding names the oldest match.
		x := newEvidenceRepo(t, builder)
		in := x.r.input()
		newer, older := commit(x.cs, []string{"builder"}, nil), commit(x.cs, []string{"builder"}, nil)
		newer.Hash, older.Hash = strings.Repeat("a", 40), strings.Repeat("b", 40)
		in.Commits = []Commit{newer, older}
		fs := ofRule(Run(in), RuleEvidenceActor)
		if len(fs) != 1 || !strings.HasSuffix(fs[0].Message, "commit bbbbbbbbbbbb") {
			t.Errorf("evidence-actor findings %+v", fs)
		}
	})
	t.Run("RetiredEvidence", func(t *testing.T) {
		x := newEvidenceRepo(t, builder)
		x.r.record("Retire", op.Op{ID: x.r.next(), Type: op.TypeNodeRetire, Node: x.ev, Base: x.evOp, Reason: "stale"})
		in := x.r.input()
		in.Commits = []Commit{commit(x.cs, []string{"builder"}, nil)}
		if fs := ofRule(Run(in), RuleEvidenceActor); len(fs) != 0 {
			t.Errorf("evidence-actor findings %+v", fs)
		}
	})
	t.Run("RetiredLink", func(t *testing.T) {
		x := newEvidenceRepo(t, builder)
		x.r.record("Unlink", op.Op{ID: x.r.next(), Type: op.TypeLinkRetire, Link: x.l, Base: x.evLinkOp})
		in := x.r.input()
		in.Commits = []Commit{commit(x.cs, []string{"builder"}, nil)}
		if fs := ofRule(Run(in), RuleEvidenceActor); len(fs) != 0 {
			t.Errorf("evidence-actor findings %+v", fs)
		}
	})
}

func TestRunUnlinkedConstraint(t *testing.T) {
	r := newRepo(t)
	_, capab := clean(r)
	linked, related, bare := r.next(), r.next(), r.next()
	relates := mkLink(r.next(), r.next(), op.LinkRelatesTo, related, capab)
	relates.Note = "context only"
	r.record("Constraints",
		mkNode(r.next(), linked, op.KindConstraint, "linked"),
		mkNode(r.next(), related, op.KindConstraint, "related"),
		mkNode(r.next(), bare, op.KindConstraint, "bare"),
		mkLink(r.next(), r.next(), op.LinkConstrains, linked, capab),
		relates,
	)
	rep := Run(r.input())
	if !rep.OK {
		t.Fatalf("not OK: %+v", rep.Findings)
	}
	got := ofRule(rep, RuleLink)
	if len(got) != 2 {
		t.Fatalf("link findings = %+v, want 2", got)
	}
	for i, name := range []string{"bare", "related"} {
		if got[i].Severity != SeverityWarn || !strings.Contains(got[i].Message, "constraint "+name+" ") {
			t.Errorf("finding %d = %+v, want warn on %s", i, got[i], name)
		}
	}
	if len(rep.Findings) != 2 {
		t.Errorf("findings = %+v, want only the two", rep.Findings)
	}
}
