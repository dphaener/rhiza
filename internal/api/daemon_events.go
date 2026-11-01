package api

// This file holds the daemon's event hook: the owner's command, run once
// per change set lifecycle move.

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
)

// eventNames are the event names of the lifecycle moves.
var eventNames = map[entry.Type]string{
	entry.TypePropose:   "proposed",
	entry.TypeAccept:    "accepted",
	entry.TypeReject:    "rejected",
	entry.TypeWithdraw:  "withdrawn",
	entry.TypeSupersede: "superseded",
	entry.TypeLand:      "landed",
	entry.TypeReopen:    "reopened",
}

// lifecycleEvent is one lifecycle move, as the event command gets it on stdin.
type lifecycleEvent struct {
	Event  string `json:"event"`
	Change string `json:"change"`
	Entry  string `json:"entry"`
	Title  string `json:"title"`
	At     string `json:"at"`
	Actor  string `json:"actor"`
}

// daemonEvents runs command once for each lifecycle move in the log that
// state does not list as handled, in log (entry ID) order, and records it.
// The first run, with or without a command, records every move already in
// the log and runs nothing, so history is not replayed. A failing or
// timed-out command is logged with the entry ID and not retried. Without a
// command nothing is read nor recorded, so a command set later still
// starts from the moves of that time.
func daemonEvents(ctx *Ctx, state *daemon.State, command string) error {
	if command == "" {
		return nil
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return err
	}
	type move struct {
		ev lifecycleEvent
		id id.ID
	}
	m0 := w.model(id.ID{})
	var moves []move
	for _, c := range w.changeSets() {
		for _, mv := range c.Lifecycle {
			name, ok := eventNames[mv.Type]
			if !ok {
				continue
			}
			moves = append(moves, move{id: mv.Entry, ev: lifecycleEvent{
				Event: name, Change: c.ID.String(), Entry: mv.Entry.String(),
				Title: c.Open.DisplayTitle(), At: rfc3339(mv.At), Actor: feedActor(m0, mv.Actor),
			}})
		}
	}
	slices.SortFunc(moves, func(a, b move) int { return id.Compare(a.id, b.id) })
	seeding := !state.EventsSeeded
	for _, m := range moves {
		if slices.Contains(state.Events, m.ev.Entry) {
			continue
		}
		state.Events = append(state.Events, m.ev.Entry)
		if seeding {
			continue
		}
		in, err := json.Marshal(m.ev)
		if err == nil {
			err = runHookCommand(ctx, command, in, []string{
				"RHIZA_EVENT=" + m.ev.Event,
				"RHIZA_EVENT_CHANGE=" + m.ev.Change,
				"RHIZA_EVENT_ENTRY=" + m.ev.Entry,
				"RHIZA_EVENT_TITLE=" + m.ev.Title,
			})
		}
		if err != nil && ctx.Stderr != nil {
			fmt.Fprintf(ctx.Stderr, "rhiza: daemon: event command for entry %s: %v\n", m.ev.Entry, strings.TrimSpace(err.Error()))
		}
	}
	state.EventsSeeded = true
	return nil
}
