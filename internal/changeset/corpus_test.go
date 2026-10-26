package changeset_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

// corpusDir holds sample logs from every released format; see its README.
const corpusDir = "../../testdata/corpus"

// corpusWant is one change set's expected fold result, found by its intent.
type corpusWant struct {
	intent string
	state  changeset.State
	// supersededBy is the intent of the change set that superseded it.
	supersededBy string
	// handles lists the handles of the current ops, in order ("" for a link).
	handles []string
	// decided is the type of the deciding entry, if any.
	decided entry.Type
}

// corpusExpect lists every corpus case's change sets, in change set ID order.
var corpusExpect = map[string][]corpusWant{
	"format-1/single-change-set": {
		{intent: "Add a café menu to the landing page", state: changeset.StateLanded,
			handles: []string{"cafe-menu"}, decided: entry.TypeAccept},
	},
	"format-1/lifecycle": {
		{intent: "Rename the config file", state: changeset.StateRejected,
			handles: []string{"config-file-name"}, decided: entry.TypeReject},
		{intent: "Bump the Go version", state: changeset.StateWithdrawn},
		{intent: "Write the README", state: changeset.StateSuperseded, supersededBy: "Write the README and the guide"},
		{intent: "Write the README and the guide", state: changeset.StateDraft},
	},
	"format-1/concurrent": {
		{intent: "Fix the flaky test", state: changeset.StateProposed, handles: []string{"stable-tests"}},
	},
	"format-1/revisions": {
		{intent: "Let people export reports as CSV", state: changeset.StateLanded,
			handles: []string{"report-export", "share-reports"}, decided: entry.TypeAccept},
	},
	"format-2/carve": {
		{intent: "Record actor ada@example.com", state: changeset.StateLanded, handles: []string{"ada"}},
		{intent: "Plan the reports area", state: changeset.StateDraft, handles: []string{"report-schedule", ""}},
		{intent: "Ship report export", state: changeset.StateProposed,
			handles: []string{"reports", "report-export", "csv-first", ""}},
	},
}

func TestCorpusFolds(t *testing.T) {
	cases, err := filepath.Glob(filepath.Join(corpusDir, "format-*", "*"))
	if err != nil {
		t.Fatal(err)
	}
	cases = slices.DeleteFunc(cases, func(p string) bool {
		info, err := os.Stat(p)
		return err != nil || !info.IsDir()
	})
	if len(cases) == 0 {
		t.Fatalf("no corpus cases under %s", corpusDir)
	}
	seen := map[string]bool{}
	for _, root := range cases {
		name := filepath.Base(filepath.Dir(root)) + "/" + filepath.Base(root)
		seen[name] = true
		t.Run(name, func(t *testing.T) {
			want, ok := corpusExpect[name]
			if !ok {
				t.Fatalf("no expectation for corpus case %s: add one to corpusExpect", name)
			}
			s, err := store.Open(root)
			if err != nil {
				t.Fatal(err)
			}
			l := s.ReadLog()
			for _, p := range l.Problems {
				t.Errorf("store problem at %s: %v", p.Path, p.Err)
			}
			got, problems := changeset.FoldLog(l)
			for _, p := range problems {
				t.Errorf("fold problem: %v", p)
			}
			if len(got) != len(want) {
				t.Fatalf("%d change sets, want %d", len(got), len(want))
			}
			byIntent := map[string]id.ID{}
			for _, cs := range got {
				byIntent[cs.Open.Intent] = cs.ID
			}
			for i, w := range want {
				cs := got[i]
				if cs.Open.Intent != w.intent {
					t.Errorf("change set %d: intent %q, want %q", i, cs.Open.Intent, w.intent)
					continue
				}
				if cs.State != w.state {
					t.Errorf("%q: state %s, want %s", w.intent, cs.State, w.state)
				}
				var wantBy id.ID
				if w.supersededBy != "" {
					wantBy = byIntent[w.supersededBy]
					if wantBy.IsZero() {
						t.Errorf("%q: no change set %q", w.intent, w.supersededBy)
					}
				}
				if cs.SupersededBy != wantBy {
					t.Errorf("%q: SupersededBy %s, want %s", w.intent, cs.SupersededBy, wantBy)
				}
				var handles []string
				for _, o := range cs.Ops {
					handles = append(handles, o.Handle)
				}
				if !slices.Equal(handles, w.handles) {
					t.Errorf("%q: op handles %q, want %q", w.intent, handles, w.handles)
				}
				var decided entry.Type
				if cs.Decided != nil {
					decided = cs.Decided.Type
				}
				if decided != w.decided {
					t.Errorf("%q: decided by %q, want %q", w.intent, decided, w.decided)
				}
			}
			switch name {
			case "format-1/revisions":
				checkRevisions(t, got[0])
			case "format-2/carve":
				checkCarve(t, l.Entries, got[1], got[2])
			}
		})
	}
	for name := range corpusExpect {
		if !seen[name] {
			t.Errorf("expected corpus case %s is missing", name)
		}
	}
}

// checkRevisions checks the revisions case in detail: the first op was
// replaced in place by its revision and the link was dropped.
func checkRevisions(t *testing.T, cs changeset.ChangeSet) {
	t.Helper()
	mustParse := func(s string) id.ID {
		v, err := id.Parse(s)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	if len(cs.Ops) != 2 {
		t.Fatalf("Ops = %+v, want 2", cs.Ops)
	}
	first, second := cs.Ops[0], cs.Ops[1]
	if first.ID != mustParse("01M5N90M082VD4GFF1D8X0JFWH") || first.Type != op.TypeNodeCreate ||
		first.Node != mustParse("01M5N90K101B4TAGQSP5T1JAHC") {
		t.Errorf("Ops[0] = %+v, want the revised capability create", first)
	}
	if got := string(first.Fields["statement"]); got != `"Export a report as CSV or JSON."` {
		t.Errorf("Ops[0] statement = %s, want the revised text", got)
	}
	if second.ID != mustParse("01M5N90NYRHAGNX9Y9BPR8FH8W") || second.Kind != op.KindIntent {
		t.Errorf("Ops[1] = %+v, want the intent create", second)
	}
	dropped := mustParse("01M5N90QX83KEC0593ENJAKYSA")
	for _, o := range cs.Ops {
		if o.ID == dropped {
			t.Errorf("dropped link op %s is still present", dropped)
		}
	}
	if cs.Decided == nil || cs.Decided.Rule != "owner" {
		t.Errorf("Decided = %+v, want rule owner", cs.Decided)
	}
}

// checkCarve checks the carve case in detail: the source's v2 entry drops
// the moved ops and names the target's, the target's v2 entry re-creates
// them under the same node and link IDs and names the source's, every other
// entry is v1, and the source depends on the target.
func checkCarve(t *testing.T, entries []entry.Entry, src, tgt changeset.ChangeSet) {
	t.Helper()
	var se, te *entry.Entry
	var srcBody, tgtBody changeset.Ops
	for i, e := range entries {
		if e.V == 1 {
			continue
		}
		if e.V != 2 || e.Type != entry.TypeOps {
			t.Errorf("entry %s: v%d %s, want only ops entries at v2", e.ID, e.V, e.Type)
			continue
		}
		b, err := changeset.DecodeBody(e)
		if err != nil {
			t.Fatalf("entry %s: %v", e.ID, err)
		}
		ops := b.(changeset.Ops)
		switch e.ChangeSet {
		case src.ID:
			se, srcBody = &entries[i], ops
		case tgt.ID:
			te, tgtBody = &entries[i], ops
		}
	}
	if se == nil || te == nil {
		t.Fatalf("want a v2 entry in each of %s and %s", src.ID, tgt.ID)
	}
	if c := srcBody.Carve; c == nil || c.To != tgt.ID || c.Entry != te.ID || !c.From.IsZero() {
		t.Errorf("source carve = %+v, want to %s, entry %s", c, tgt.ID, te.ID)
	}
	if c := tgtBody.Carve; c == nil || c.From != src.ID || c.Entry != se.ID || !c.To.IsZero() {
		t.Errorf("target carve = %+v, want from %s, entry %s", c, src.ID, se.ID)
	}
	if len(srcBody.Drop) != len(tgtBody.Ops) || len(tgtBody.Drop) != 0 {
		t.Errorf("source drops %d, target creates %d and drops %d", len(srcBody.Drop), len(tgtBody.Ops), len(tgtBody.Drop))
	}
	for _, o := range tgtBody.Ops {
		if slices.Contains(srcBody.Drop, o.ID) {
			t.Errorf("target op %s reuses a dropped op ID", o.ID)
		}
	}
	if len(srcBody.Ops) != 1 || srcBody.Ops[0].LinkType != op.LinkDependsOn ||
		srcBody.Ops[0].From != src.ID || srcBody.Ops[0].To != tgt.ID {
		t.Errorf("source ops = %+v, want one depends-on link to the target", srcBody.Ops)
	}
	if !slices.Contains(changeset.DependsOn([]changeset.ChangeSet{src, tgt}, src.ID), tgt.ID) {
		t.Errorf("%s does not depend on %s", src.ID, tgt.ID)
	}
}
