package api

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
)

// modelConflictFixture is a ship land fixture whose change set updates
// store's statement on a base another change set has since superseded, and
// the other change set's ID.
func modelConflictFixture(t *testing.T) (*sessionFixture, string, string) {
	t.Helper()
	f, _, state := shipFixture(t)
	recordWrite(t, f.ctx, "add", `{"kind":"capability","handle":"store"}`)
	mustWrite(t, f.ctx, "set", `{"ref":"store","fields":["statement=mine"]}`)
	other := recordWrite(t, f.ctx, "set", `{"ref":"store","fields":["statement=theirs"]}`)
	if len(other.Nodes) != 1 {
		t.Fatalf("other = %+v", other)
	}
	for _, v := range []string{"propose", "accept"} {
		if _, err := callOp(t, f.ctx, v, `{}`); err != nil {
			t.Fatalf("%s: %v", v, err)
		}
	}
	if err := writeShipState(state, shipState{PR: 7, URL: "https://github.com/o/r/pull/7", Branch: "cs/ship-it", Head: "abc1234def",
		Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"}}}); err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{
		"symbolic-ref --short HEAD":           "cs/ship-it\n",
		"add -- .rhiza":                       "",
		"diff --cached --name-only -- .rhiza": ".rhiza/rendered/x.md\n",
		"commit -m " + shipLandMessage + " -m Change: " + currentOf(t, f.ctx).String() + "\nActor: " + shipTestActor + " -- .rhiza": "",
		"rev-parse HEAD": "beef000111\n",
		"rev-parse --abbrev-ref --symbolic-full-name @{upstream}": "origin/cs/ship-it\n",
		"rev-list --count @{upstream}..HEAD":                      "1\n",
		"push":                                                    "",
	} {
		f.git.out[k] = v
	}
	f.ctx.ActorRef = shipTestActor
	return f, state, ""
}

func wantRecorded(t *testing.T, state string, want bool) *ShipModelConflict {
	t.Helper()
	st, err := readShipState(state)
	if err != nil {
		t.Fatal(err)
	}
	if (st.ModelConflict != nil) != want {
		t.Fatalf("recorded model conflict = %+v, want recorded %v", st.ModelConflict, want)
	}
	return st.ModelConflict
}

func TestShipLandRecordsModelConflict(t *testing.T) {
	f, state, _ := modelConflictFixture(t)
	_, err := callOp(t, f.ctx, "ship.land", `{}`)
	wantExit(t, err, 3)
	e := AsError(err)
	if e.Conflict == nil || e.Conflict.Node == "" || e.Conflict.Op == "" || e.Conflict.ChangeSet == "" ||
		!slices.Equal(e.Conflict.Aspects, []string{"field:statement"}) {
		t.Fatalf("conflict = %+v", e.Conflict)
	}
	if body := ErrorJSON(err).Error; body.Conflict == nil || *body.Conflict.ptr() != *e.Conflict.ptr() {
		t.Errorf("--json body = %+v", body)
	}
	mc := wantRecorded(t, state, true)
	if mc.Head != "beef000111" || mc.Node != e.Conflict.Node || mc.Op != e.Conflict.Op || mc.ChangeSet != e.Conflict.ChangeSet {
		t.Errorf("recorded = %+v, error %+v", mc, e.Conflict)
	}
	if !strings.Contains(e.Message, mc.ChangeSet) {
		t.Errorf("message %q does not name %s", e.Message, mc.ChangeSet)
	}
}

// ptr gives a comparable copy of c (Aspects joined).
func (c ModelConflict) ptr() *[4]string {
	return &[4]string{c.Node, c.Op, c.ChangeSet, strings.Join(c.Aspects, ",")}
}

// conflictHead is rhizaRepoFixture's head build answering land with the
// body --json prints for the conflict landing in-process gives.
func TestShipLandHeadBuildRecordsModelConflict(t *testing.T) {
	f, state, _ := modelConflictFixture(t)
	// What rhiza land --json prints for this conflict.
	lctx := *f.ctx
	lctx.ChangeSet, lctx.ChangeRef, lctx.Record = currentOf(t, f.ctx), "", false
	land, _ := Lookup("land")
	_, lerr := land.Run(&lctx, LifecycleRequest{})
	if lerr == nil {
		t.Fatal("land succeeded")
	}
	body, err := json.Marshal(ErrorJSON(lerr))
	if err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(t.TempDir(), "rhiza")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	f.ctx.Getenv = func(k string) string {
		if k == envShipRhiza {
			return bin
		}
		return ""
	}
	inner := f.ctx.Exec
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if name == bin {
			return body, errors.New("exit status 3")
		}
		return inner(name, args, dir, env)
	}
	_, err = callOp(t, f.ctx, "ship.land", `{}`)
	wantExit(t, err, 1)
	mc := wantRecorded(t, state, true)
	want := AsError(lerr).Conflict
	if mc.Head != "beef000111" || mc.Node != want.Node || mc.Op != want.Op || mc.ChangeSet != want.ChangeSet || !slices.Equal(mc.Aspects, want.Aspects) {
		t.Errorf("recorded = %+v, want %+v", mc, want)
	}
}

func TestShipLandClearsModelConflict(t *testing.T) {
	f, state, _ := modelConflictFixture(t)
	cs := currentOf(t, f.ctx)
	if _, err := callOp(t, f.ctx, "ship.land", `{}`); err == nil {
		t.Fatal("land succeeded")
	}
	wantRecorded(t, state, true)

	// A land refused for a missing dependency leaves none recorded.
	dep := newChange(t, f.ctx, `{"intent":"Dep"}`).ChangeSet
	c := *f.ctx
	c.ChangeRef = cs.String()
	mustWrite(t, &c, "link", `{"from":"`+cs.String()+`","type":"depends-on","to":"`+dep.String()+`"}`)
	f.ctx.ChangeRef = cs.String()
	_, err := callOp(t, f.ctx, "ship.land", `{}`)
	wantExit(t, err, 3)
	if !strings.Contains(err.Error(), "it depends on") {
		t.Fatalf("err = %v", err)
	}
	wantRecorded(t, state, false)
}

func TestShipLandSucceedsAfterRebase(t *testing.T) {
	f, state, _ := modelConflictFixture(t)
	cs := currentOf(t, f.ctx)
	if _, err := callOp(t, f.ctx, "ship.land", `{}`); err == nil {
		t.Fatal("land succeeded")
	}
	wantRecorded(t, state, true)
	mustWrite(t, f.ctx, "set", `{"ref":"store","fields":["statement=merged"],"rebase":true}`)
	resp, err := callOp(t, f.ctx, "ship.land", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.(ShipLandResponse).Landed {
		t.Fatalf("land = %+v", resp)
	}
	wantRecorded(t, state, false)
	if c, _ := ctxWorkspace(t, f.ctx).changeSet(cs); c.State != changeset.StateLanded {
		t.Errorf("state = %s", c.State)
	}
}

func TestShipCheckModelConflict(t *testing.T) {
	f, state, _ := modelConflictFixture(t)
	head, fixes := "beef000111\n", ""
	inner := f.ctx.Exec
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if name == "git" {
			switch strings.Join(args, " ") {
			case "rev-parse HEAD":
				return []byte(head), nil
			case "log --format=%H -E --grep=^Review: RUN1$ --grep=^Model-Conflict: --all-match HEAD":
				return []byte(fixes), nil
			}
		}
		return inner(name, args, dir, env)
	}
	getenv := func(k string) string {
		if k == "RHIZA_RUN" {
			return "RUN1"
		}
		return ""
	}
	f.ctx.Getenv = getenv
	const req = `{"model_conflict":true}`
	fails := func(why, detail string) {
		t.Helper()
		r := shipCheck(t, f, req)
		if r.Passed || r.ExitCode() != 1 || !strings.Contains(r.Detail, detail) {
			t.Errorf("%s: %+v", why, r)
		}
	}

	fails("nothing recorded", "no model conflict is recorded")

	cs := currentOf(t, f.ctx)
	if _, err := callOp(t, f.ctx, "ship.land", `{}`); err == nil {
		t.Fatal("land succeeded")
	}
	mc := wantRecorded(t, state, true)

	r := shipCheck(t, f, req)
	if !r.Passed || r.ExitCode() != 0 || r.ModelConflict == nil {
		t.Fatalf("recorded: %+v", r)
	}
	v := r.ModelConflict
	if v.ChangeSetID != cs.String() || v.Intent != "Ship it" || v.OtherChangeSet != mc.ChangeSet || v.OtherChangeIntent == "" || v.Head != "beef000111" {
		t.Errorf("view = %+v", v)
	}
	if b, _ := json.Marshal(r); !strings.Contains(string(b), `"other_intent"`) || !strings.Contains(string(b), `"refused_change_set"`) {
		t.Errorf("json = %s", b)
	}

	fixes = "a1\nb2\n"
	if r := shipCheck(t, f, req); !r.Passed {
		t.Errorf("two fixes: %+v", r)
	}
	fixes = "a1\nb2\nc3\n"
	fails("three fixes", "3 model-conflict fixes")
	fixes = ""

	head = "cafe000222\n"
	fails("branch moved", "moved past")
	head = "beef000111\n"

	// A land refused for a missing dependency clears it.
	dep := newChange(t, f.ctx, `{"intent":"Dep"}`).ChangeSet
	c := *f.ctx
	c.ChangeRef = cs.String()
	mustWrite(t, &c, "link", `{"from":"`+cs.String()+`","type":"depends-on","to":"`+dep.String()+`"}`)
	f.ctx.ChangeRef = cs.String()
	if _, err := callOp(t, f.ctx, "ship.land", `{}`); err == nil {
		t.Fatal("land succeeded")
	}
	fails("missing dependency", "no model conflict is recorded")
}
