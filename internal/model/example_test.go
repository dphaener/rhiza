package model_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

func mustID(s string) id.ID {
	v, err := id.Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Example_landedConflict shows a landed conflict and the small change set
// that repairs it. Change set A creates a decision. B and C both rewrite its
// statement from A's create, and both land: whichever applies second (C, by
// ID, since neither builds on the other) wrote over a change it never saw.
// As-is still applies C, so C's statement wins, and reports C's update as a
// conflict. D, written against the decision's current version, settles the
// statement. It is not stale, so it brings no problem of its own; the B/C
// conflict stays in the history and is still reported.
func Example_landedConflict() {
	csA, csB := mustID("01K8Z0000000000000000000A0"), mustID("01K8Z0000000000000000000B0")
	csC, csD := mustID("01K8Z0000000000000000000C0"), mustID("01K8Z0000000000000000000D0")
	auth := mustID("01K8Z0000000000000000000N0")
	createOp, bOp := mustID("01K8Z0000000000000000000A1"), mustID("01K8Z0000000000000000000B1")
	cOp, dOp := mustID("01K8Z0000000000000000000C1"), mustID("01K8Z0000000000000000000D1")
	names := map[id.ID]string{csA: "A", csB: "B", csC: "C", csD: "D", createOp: "A's create", bOp: "B's update",
		cOp: "C's update", dOp: "D's update"}

	statement := func(s string) map[string]json.RawMessage {
		b, _ := json.Marshal(s)
		return map[string]json.RawMessage{"statement": b}
	}
	landed := func(cs id.ID, o op.Op) changeset.ChangeSet {
		return changeset.ChangeSet{ID: cs, State: changeset.StateLanded, Ops: []op.Op{o}}
	}
	update := func(o, base id.ID, s string) op.Op {
		return op.Op{ID: o, Type: op.TypeNodeUpdate, Node: auth, Base: base, Set: statement(s)}
	}
	show := func(sets []changeset.ChangeSet) {
		m, problems := model.AsIs(sets)
		for _, p := range problems {
			fmt.Printf("problem: %s of change set %s, conflict: %t\n", names[p.Op], names[p.ChangeSet],
				errors.Is(p.Err, model.ErrConflict))
		}
		n, _ := m.Node(auth)
		fmt.Printf("statement: %s (version: %s)\n", n.Fields["statement"], names[n.Version])
	}

	sets := []changeset.ChangeSet{
		landed(csA, op.Op{ID: createOp, Type: op.TypeNodeCreate, Node: auth, Kind: op.KindDecision, Handle: "auth",
			Fields: statement("Use sessions.")}),
		landed(csB, update(bOp, createOp, "Use JWT.")),
		landed(csC, update(cOp, createOp, "Use OAuth.")),
	}
	show(sets)

	// D repairs it from the current version, C's update.
	fmt.Println("after D:")
	show(append(sets, landed(csD, update(dOp, cOp, "Use OAuth with JWT access tokens."))))

	// Output:
	// problem: C's update of change set C, conflict: true
	// statement: "Use OAuth." (version: C's update)
	// after D:
	// problem: C's update of change set C, conflict: true
	// statement: "Use OAuth with JWT access tokens." (version: D's update)
}
