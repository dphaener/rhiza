package api

// This file holds the status operation: the binary's version, the current
// change set, the change sets in flight by state and the installed guides,
// with the commands that would help next. With --hook it is the
// session-start hook's output: brief, silent outside a repository or on any
// error, and always exit 0. Status reads the log and the guide folders
// only; it never writes.

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/store"
)

func init() {
	register(Operation{
		Name:     "status",
		Summary:  "Show the current change set, what is in flight and the installed guide",
		Request:  reflect.TypeFor[StatusRequest](),
		Response: reflect.TypeFor[StatusResponse](),
		Run:      runStatus,
	})
}

// StatusRequest asks for the full status or the hook's brief one.
type StatusRequest struct {
	Hook bool `json:"hook,omitempty" desc:"Brief output for a session-start hook: silent outside a repo, always exits 0"`
}

// StatusResponse is the state of the repository and the guide installs.
// A silent response (the hook outside a repository or after an error)
// prints nothing and encodes as {}.
type StatusResponse struct {
	Version   string           `json:"version" desc:"The running rhiza's version"`
	ChangeSet *StatusChangeSet `json:"change_set" desc:"The current change set, as writes resolve it; null if none"`
	Attention StatusAttention  `json:"attention" desc:"The change sets in flight, counted by state"`
	Guide     []StatusGuide    `json:"guide" desc:"Each installed guide found, user then project"`
	Next      []string         `json:"next" desc:"Commands that would help, in order; empty if none"`

	hook   bool
	silent bool
	notice string
}

// StatusChangeSet is the current change set.
type StatusChangeSet struct {
	ID     id.ID           `json:"id" desc:"The change set's ID"`
	Intent string          `json:"intent" desc:"The change set's intent"`
	Kind   changeset.Kind  `json:"kind" desc:"intention or record; empty without a valid open"`
	State  changeset.State `json:"state" desc:"The lifecycle state; empty without a valid open"`
}

// StatusAttention counts the change sets in flight by state. Terminal
// change sets are not counted.
type StatusAttention struct {
	Draft    int `json:"draft" desc:"Change sets in draft"`
	Proposed int `json:"proposed" desc:"Change sets proposed, waiting for a decision"`
	Accepted int `json:"accepted" desc:"Change sets accepted, waiting to land"`
}

// StatusGuide is one installed guide.
type StatusGuide struct {
	Scope   string `json:"scope" enum:"user,project" desc:"user (in the home directory) or project (in the repository)"`
	Dir     string `json:"dir" desc:"The plugin folder"`
	Version string `json:"version" desc:"The version it was installed by"`
	Current bool   `json:"current" desc:"Whether it is the running rhiza's version and every file it installs is on disk as that version writes it"`
}

// oldFolderNotice is the hook's one line for an old-style root.
const oldFolderNotice = "This repo keeps its model in " + store.OldDir + "/, the folder's name before the rename to Rhiza; rename it with: " + oldFolderNext

// guideInstallNext is the command that installs or updates the user guide.
const guideInstallNext = "rhiza guide install --agent claude"

// hookIntentMax is the most bytes of intent the hook prints.
const hookIntentMax = 60

// runStatus reports the status. With Hook, every error becomes a silent
// response, so the hook never fails a session start.
func runStatus(ctx *Ctx, req any) (any, error) {
	r := req.(StatusRequest)
	resp, err := status(ctx)
	if r.Hook {
		if isOldFolderError(err) {
			return StatusResponse{hook: true, silent: true, notice: oldFolderNotice}, nil
		}
		if err != nil {
			return StatusResponse{hook: true, silent: true}, nil
		}
		resp.hook = true
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// status gathers the full status of ctx.Root.
func status(ctx *Ctx) (StatusResponse, error) {
	w, err := openWorkspace(ctx.Root, ctx.Now, ctx.IDs, true)
	if err != nil {
		return StatusResponse{}, err
	}
	cs, err := w.effectiveChangeSet(ctx)
	if err != nil {
		return StatusResponse{}, err
	}
	resp := StatusResponse{Version: ctx.Version, Guide: []StatusGuide{}, Next: []string{}}
	sets, _ := changeset.FoldLog(w.log())
	for _, c := range sets {
		switch c.State {
		case changeset.StateDraft:
			resp.Attention.Draft++
		case changeset.StateProposed:
			resp.Attention.Proposed++
		case changeset.StateAccepted:
			resp.Attention.Accepted++
		}
		if !cs.IsZero() && c.ID == cs {
			resp.ChangeSet = &StatusChangeSet{ID: c.ID, Intent: c.Open.Intent, Kind: c.Open.Kind, State: c.State}
		}
	}
	if !cs.IsZero() && resp.ChangeSet == nil {
		// Current names a change set the log does not hold.
		resp.ChangeSet = &StatusChangeSet{ID: cs}
	}
	type place struct{ scope, base, next string }
	var places []place
	if ctx.Home != "" {
		places = append(places, place{"user", ctx.Home, guideInstallNext})
	}
	places = append(places, place{"project", ctx.Root, guideInstallNext + " --project"})
	for _, p := range places {
		dir := filepath.Join(p.base, ".claude", "skills", "rhiza")
		v, ok := installedVersion(dir)
		if !ok {
			continue
		}
		g := StatusGuide{Scope: p.scope, Dir: dir, Version: v, Current: v == ctx.Version && pluginCurrent(dir, ctx.Version)}
		resp.Guide = append(resp.Guide, g)
		if !g.Current {
			resp.Next = append(resp.Next, p.next)
		}
	}
	if resp.ChangeSet == nil {
		resp.Next = append(resp.Next, newEntryNext)
	}
	return resp, nil
}

// MarshalJSON encodes a silent response as {} and any other as its fields.
func (r StatusResponse) MarshalJSON() ([]byte, error) {
	if r.silent {
		return []byte("{}"), nil
	}
	type plain StatusResponse
	return json.Marshal(plain(r))
}

// ExitCode is always 0: a status is never a failure the caller acts on,
// and the hook must not fail a session start.
func (r StatusResponse) ExitCode() int { return 0 }

// Human prints the full status, the hook's brief one, or nothing when
// silent.
func (r StatusResponse) Human(w io.Writer) {
	switch {
	case r.notice != "":
		fmt.Fprintln(w, r.notice)
	case r.silent:
		return
	case r.hook:
		r.humanHook(w)
		return
	}
	fmt.Fprintf(w, "rhiza %s\n\n", r.Version)
	if c := r.ChangeSet; c != nil {
		fmt.Fprintf(w, "- Change set: %s (%s): %s\n", c.ID, c.State, c.Intent)
	} else {
		fmt.Fprintln(w, "- Change set: none current")
	}
	fmt.Fprintf(w, "- In flight: %s\n", r.Attention.counts())
	if len(r.Guide) == 0 {
		fmt.Fprintln(w, "- Guide: not installed")
	}
	for _, g := range r.Guide {
		note := "current"
		if !g.Current {
			note = "out of date"
		}
		fmt.Fprintf(w, "- Guide (%s): %s in %s (%s)\n", g.Scope, g.Version, g.Dir, note)
	}
	if len(r.Next) > 0 {
		fmt.Fprintln(w)
	}
	for _, n := range r.Next {
		fmt.Fprintf(w, "next: %s\n", n)
	}
}

// humanHook prints the brief status: at most guide.StartCap tokens. Every
// line is bounded: the intent is truncated and no version is printed.
func (r StatusResponse) humanHook(w io.Writer) {
	fmt.Fprintln(w, "This repo uses Rhiza (load the rhiza skill; `rhiza guide` for more).")
	if c := r.ChangeSet; c != nil {
		fmt.Fprintf(w, "Change set %s (%s): %s\n", c.ID, c.State, truncate(c.Intent, hookIntentMax))
	} else {
		fmt.Fprintln(w, "No current change set.")
	}
	fmt.Fprintf(w, "In flight: %s.\n", r.Attention.counts())
	var updates []string
	for _, n := range r.Next {
		if strings.HasPrefix(n, guideInstallNext) {
			updates = append(updates, "`"+n+"`")
		}
	}
	if len(updates) > 0 {
		// No versions: a version string has no bound, the cap does.
		fmt.Fprintf(w, "The installed rhiza guide is out of date: run %s.\n", strings.Join(updates, ", "))
	}
}

// counts lists the in-flight counts.
func (a StatusAttention) counts() string {
	return fmt.Sprintf("%d draft, %d proposed, %d accepted", a.Draft, a.Proposed, a.Accepted)
}

// truncate shortens s to at most n bytes, cutting at a rune boundary and
// ending with "…" when it cuts. Line breaks become spaces.
func truncate(s string, n int) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= n {
		return s
	}
	const ell = "…"
	cut := n - len(ell)
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut] + ell
}
