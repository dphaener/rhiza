package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGates(t *testing.T) {
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	hist := RunHistory{"abc": now.Add(-30 * time.Hour), "fresh": now.Add(-2 * time.Hour)}
	base := GateInputs{Now: now, Commit: "abc", Kind: KindPatch, CI: "success", Runs: hist}
	status := func(gs []Gate, name string) Gate {
		for _, g := range gs {
			if g.Name == name {
				return g
			}
		}
		t.Fatalf("no gate %s", name)
		return Gate{}
	}
	mod := func(f func(*GateInputs)) GateInputs { in := base; f(&in); return in }
	for _, c := range []struct {
		name string
		in   GateInputs
		gate string
		pass bool
	}{
		{"all pass", base, "", true},
		{"ci red", mod(func(i *GateInputs) { i.CI = "failure" }), "ci", false},
		{"ci missing", mod(func(i *GateInputs) { i.CI = "" }), "ci", false},
		{"commit missing from history", mod(func(i *GateInputs) { i.Commit = "zzz" }), "soak", false},
		{"soak too short", mod(func(i *GateInputs) { i.Commit = "fresh" }), "soak", false},
		{"bug filed since", mod(func(i *GateInputs) { i.BugRequests = []time.Time{now.Add(-3 * time.Hour)} }), "soak", false},
		{"bug before the soak", mod(func(i *GateInputs) { i.BugRequests = []time.Time{now.Add(-40 * time.Hour)} }), "soak", true},
		{"milestone in flight", mod(func(i *GateInputs) { i.InFlight = []string{"CS1"} }), "milestone", false},
		{"milestone blocked", mod(func(i *GateInputs) { i.Blocking = []string{"Q1"} }), "milestone", false},
		{"second patch in 24h", mod(func(i *GateInputs) { i.LastPatch = now.Add(-5 * time.Hour) }), "patch", false},
		{"patch after 24h", mod(func(i *GateInputs) { i.LastPatch = now.Add(-25 * time.Hour) }), "patch", true},
		{"severe skips soak", mod(func(i *GateInputs) { i.Severe, i.Commit = true, "zzz" }), "soak", true},
		{"severe skips batching", mod(func(i *GateInputs) { i.Severe, i.LastPatch = true, now.Add(-time.Hour) }), "patch", true},
		{"severe never skips ci", mod(func(i *GateInputs) { i.Severe, i.CI = true, "failure" }), "ci", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			gs := Gates(c.in)
			if c.gate == "" {
				if !AllPass(gs) {
					t.Fatalf("%+v", gs)
				}
				return
			}
			g := status(gs, c.gate)
			if g.Pass != c.pass || g.Reason == "" {
				t.Fatalf("%+v", g)
			}
			if AllPass(gs) != (c.pass && AllPass(gs)) {
				t.Fatalf("AllPass disagrees: %+v", gs)
			}
		})
	}
}

func TestReadRunHistory(t *testing.T) {
	p := filepath.Join(t.TempDir(), "runs.jsonl")
	os.WriteFile(p, []byte(`{"commit":"a","at":"2026-11-10T00:00:00Z"}`+"\n\n"+`{"commit":"a","at":"2026-11-09T00:00:00Z"}`+"\n"), 0o644)
	h, err := ReadRunHistory(p)
	if err != nil || !h["a"].Equal(time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("%v %v", h, err)
	}
	if h, err := ReadRunHistory(p + ".none"); err != nil || len(h) != 0 {
		t.Fatalf("%v %v", h, err)
	}
	os.WriteFile(p, []byte("nope\n"), 0o644)
	if _, err := ReadRunHistory(p); err == nil {
		t.Fatal("bad line accepted")
	}
}
