package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Soak is how long the root daemon must have run a commit before it ships.
const Soak = 24 * time.Hour

// PatchSpacing is the least time between two patch releases.
const PatchSpacing = 24 * time.Hour

// Run history file: JSON lines, one per commit the root daemon ran, as
// {"commit":"<sha>","at":"<RFC3339 of the first run>"}. A commit may appear
// more than once; the earliest line counts.

// RunHistory maps a commit to when the root daemon first ran it.
type RunHistory map[string]time.Time

// ReadRunHistory reads the history file at path. A missing file is an empty
// history, so every commit then fails the soak gate.
func ReadRunHistory(path string) (RunHistory, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return RunHistory{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := RunHistory{}
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e struct {
			Commit string `json:"commit"`
			At     string `json:"at"`
		}
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, n, err)
		}
		at, err := time.Parse(time.RFC3339, e.At)
		if err != nil || e.Commit == "" {
			return nil, fmt.Errorf("%s:%d: want commit and an RFC3339 at", path, n)
		}
		if old, ok := h[e.Commit]; !ok || at.Before(old) {
			h[e.Commit] = at
		}
	}
	return h, sc.Err()
}

// GateInputs is everything the gates read about the candidate commit.
type GateInputs struct {
	Now         time.Time
	Commit      string
	Kind        Kind
	Severe      bool        // a severe: fix is in the release
	CI          string      // CI conclusion on Commit: "success" passes
	Runs        RunHistory  // when the root daemon first ran each commit
	BugRequests []time.Time // when bug requests were filed
	InFlight    []string    // change sets serving the finished milestone not yet landed
	Blocking    []string    // open questions blocking one of them
	LastPatch   time.Time   // zero when no patch was released
}

// Gate is one release gate's verdict and why.
type Gate struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

// Gates evaluates every gate for the candidate commit. A severe fix skips the
// soak and the patch batching, never CI.
func Gates(in GateInputs) []Gate {
	return []Gate{gateCI(in), gateSoak(in), gateMilestone(in), gatePatch(in)}
}

// AllPass reports whether every gate passes.
func AllPass(gs []Gate) bool {
	for _, g := range gs {
		if !g.Pass {
			return false
		}
	}
	return true
}

func gateCI(in GateInputs) Gate {
	g := Gate{Name: "ci"}
	switch in.CI {
	case "success":
		g.Pass, g.Reason = true, "CI is green on "+in.Commit
	case "":
		g.Reason = "no CI status for " + in.Commit
	default:
		g.Reason = fmt.Sprintf("CI on %s is %s", in.Commit, in.CI)
	}
	return g
}

func gateSoak(in GateInputs) Gate {
	g := Gate{Name: "soak"}
	if in.Severe {
		g.Pass, g.Reason = true, "skipped: a severe fix is in the release"
		return g
	}
	first, ok := in.Runs[in.Commit]
	if !ok {
		g.Reason = in.Commit + " is not in the daemon run history"
		return g
	}
	if age := in.Now.Sub(first); age < Soak {
		g.Reason = fmt.Sprintf("the daemon first ran %s %s ago, under %s", in.Commit, age.Round(time.Minute), Soak)
		return g
	}
	for _, b := range in.BugRequests {
		if b.After(first) {
			g.Reason = "a bug request was filed at " + b.UTC().Format(time.RFC3339) + ", after the soak began"
			return g
		}
	}
	g.Pass, g.Reason = true, "the daemon has run "+in.Commit+" for over 24h with no bug request since"
	return g
}

func gateMilestone(in GateInputs) Gate {
	g := Gate{Name: "milestone"}
	switch {
	case len(in.InFlight) > 0:
		g.Reason = "change sets still in flight: " + strings.Join(in.InFlight, ", ")
	case len(in.Blocking) > 0:
		g.Reason = "questions blocking the milestone: " + strings.Join(in.Blocking, ", ")
	default:
		g.Pass, g.Reason = true, "no change set in flight and no question blocking the milestone"
	}
	return g
}

func gatePatch(in GateInputs) Gate {
	g := Gate{Name: "patch", Pass: true}
	switch {
	case in.Kind != KindPatch:
		g.Reason = "not a patch release"
	case in.Severe:
		g.Reason = "skipped: a severe fix is released at once"
	case !in.LastPatch.IsZero() && in.Now.Sub(in.LastPatch) < PatchSpacing:
		g.Pass = false
		g.Reason = fmt.Sprintf("the last patch was %s ago, under %s", in.Now.Sub(in.LastPatch).Round(time.Minute), PatchSpacing)
	default:
		g.Reason = "no patch in the last 24h"
	}
	return g
}
