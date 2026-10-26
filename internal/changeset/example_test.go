package changeset_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// ExampleFold folds a change set through its whole life: opened, given two
// ops, revised (one op replaced by repeating its ID, the other dropped),
// proposed, accepted and landed.
func ExampleFold() {
	now := time.Date(2026, 10, 23, 21, 0, 0, 0, time.UTC)
	gen := id.NewGenerator(func() time.Time { return now }, rand.New(rand.NewSource(1)))
	next := func() id.ID {
		now = now.Add(time.Second)
		v, err := gen.New()
		if err != nil {
			panic(err)
		}
		return v
	}

	owner, cs := next(), next()
	var entries []entry.Entry
	add := func(b changeset.Body) {
		body, err := changeset.EncodeBody(b)
		if err != nil {
			panic(err)
		}
		e := entry.Entry{V: changeset.EntryVersion(b), ID: next(), ChangeSet: cs, Actor: owner, For: owner,
			At: now, Type: b.Type(), Body: body}
		if len(entries) > 0 {
			e.After = []id.ID{entries[len(entries)-1].ID}
		}
		entries = append(entries, e)
	}
	statement := func(s string) map[string]json.RawMessage {
		b, _ := json.Marshal(s)
		return map[string]json.RawMessage{"statement": b}
	}

	tz, tzOp := next(), next()
	log, logOp := next(), next()
	add(changeset.Open{Kind: changeset.KindIntention, Intent: "Show times in the reader's zone"})
	add(changeset.Ops{Ops: []op.Op{
		{ID: tzOp, Type: op.TypeNodeCreate, Node: tz, Kind: op.KindCapability, Handle: "local-times", Key: "i",
			Fields: statement("Show every time in UTC.")},
		{ID: logOp, Type: op.TypeNodeCreate, Node: log, Kind: op.KindCapability, Handle: "time-log", Key: "r",
			Fields: statement("Log each conversion.")},
	}})
	add(changeset.Ops{
		Ops: []op.Op{{ID: tzOp, Type: op.TypeNodeCreate, Node: tz, Kind: op.KindCapability, Handle: "local-times",
			Key: "i", Fields: statement("Show every time in the reader's zone.")}},
		Drop: []id.ID{logOp},
	})
	add(changeset.Propose{})
	add(changeset.Accept{Rule: "owner"})
	add(changeset.Land{})

	got, problems := changeset.Fold(cs, entries)
	fmt.Println("state:", got.State)
	for _, o := range got.Ops {
		fmt.Printf("op: %s %s %s\n", o.Type, o.Handle, o.Fields["statement"])
	}
	fmt.Println("rule:", got.Decided.Rule)
	fmt.Println("problems:", len(problems))
	// Output:
	// state: landed
	// op: node.create local-times "Show every time in the reader's zone."
	// rule: owner
	// problems: 0
}
