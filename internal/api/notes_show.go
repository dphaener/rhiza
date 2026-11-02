package api

// This file holds the notes section of show and change show: the notes
// that discuss a node or a change set, oldest first, with who wrote each
// and when.

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// ShowNote is one note that discusses the node or change set shown.
type ShowNote struct {
	ID    id.ID  `json:"id" desc:"The note's ID"`
	Text  string `json:"text" desc:"The note"`
	Op    id.ID  `json:"op,omitzero" desc:"The operation of the change set the note is about, if it names one"`
	Actor string `json:"actor,omitempty" desc:"Who wrote it: the actor's handle (its ID when it has none)"`
	At    string `json:"at,omitempty" desc:"When it was written (RFC 3339, UTC)"`
}

// writeNotes prints the Notes section, only when there are notes.
func writeNotes(w io.Writer, notes []ShowNote) {
	if len(notes) == 0 {
		return
	}
	fmt.Fprintf(w, "\n## Notes\n\n")
	for _, n := range notes {
		head := "- " + shortID(n.ID)
		if n.Actor != "" {
			head += " by " + n.Actor
		}
		if n.At != "" {
			head += " at " + n.At
		}
		if !n.Op.IsZero() {
			head += " (on op " + shortID(n.Op) + ")"
		}
		fmt.Fprintf(w, "%s: %s\n", head, strings.Join(strings.Fields(n.Text), " "))
	}
}

// sortNotes puts notes oldest first, ties by ID.
func sortNotes(notes []ShowNote) {
	slices.SortStableFunc(notes, func(a, b ShowNote) int {
		if c := strings.Compare(a.At, b.At); c != 0 {
			return c
		}
		return strings.Compare(a.ID.String(), b.ID.String())
	})
}

func noteOp(raw string) id.ID {
	o, _ := id.Parse(raw)
	return o
}

// showNotes lists the live notes with a live discusses link to node n.
func showNotes(r *reader, n id.ID, links []ShowNodeLink) ([]ShowNote, error) {
	notes := []ShowNote{}
	for _, l := range links {
		if l.Type != op.LinkDiscusses || l.Retired {
			continue
		}
		nn, ok, err := r.Node(l.Other.ID)
		if err != nil {
			return nil, err
		}
		if !ok || nn.Kind != op.KindNote || nn.Retired {
			continue
		}
		note := ShowNote{ID: nn.ID, Text: fieldText(nn.Fields["statement"]), Op: noteOp(fieldText(nn.Fields["op"]))}
		if len(nn.History) > 0 {
			cs := nn.History[0].ChangeSet
			es, err := r.Entries(cs)
			if err != nil {
				return nil, err
			}
			c, _ := changeset.Fold(cs, es)
			if rs := c.History[nn.History[0].Op]; len(rs) > 0 {
				note.Actor, note.At = noteWho(rs[0].Actor, rs[0].Entry, es, func(a id.ID) string {
					an, ok, err := r.Node(a)
					if err != nil || !ok {
						return ""
					}
					return an.Handle
				})
			}
		}
		notes = append(notes, note)
	}
	sortNotes(notes)
	return notes, nil
}

// noteWho names the actor (by handle, else ID) and the time of entry e.
func noteWho(actor, e id.ID, es []entry.Entry, handle func(id.ID) string) (string, string) {
	who := ""
	if !actor.IsZero() {
		who = handle(actor)
		if who == "" {
			who = actor.String()
		}
	}
	at := ""
	for _, x := range es {
		if x.ID == e {
			at = x.At.UTC().Format(lineageAtLayout)
		}
	}
	return who, at
}

// changeNotes lists the live notes discussing change set cs, finding each
// note's author and time among the change sets of w. It looks in cs's own
// view and, when it is another change set, the current one's: a note about
// cs is written into the writer's current change set, so a reviewer's note
// shows before their change set lands.
func changeNotes(w *workspace, cs, current id.ID) []ShowNote {
	notes := []ShowNote{}
	sets := w.changeSets()
	views := []Model{w.model(cs)}
	if !current.IsZero() && current != cs {
		views = append(views, w.model(current))
	}
	seen := map[id.ID]bool{}
	for _, m := range views {
		handle := func(a id.ID) string {
			if n, ok := m.Node(a); ok {
				return n.Handle
			}
			return ""
		}
		for _, l := range m.LinksTo(cs, op.LinkDiscusses) {
			nn, ok := m.Node(l.From)
			if !ok || nn.Kind != op.KindNote || nn.Retired || seen[nn.ID] {
				continue
			}
			seen[nn.ID] = true
			note := ShowNote{ID: nn.ID, Text: fieldText(nn.Fields["statement"]), Op: noteOp(fieldText(nn.Fields["op"]))}
		find:
			for _, s := range sets {
				for _, rs := range s.History {
					if len(rs) == 0 || rs[0].Op.Type != op.TypeNodeCreate || rs[0].Op.Node != nn.ID {
						continue
					}
					note.Actor, note.At = noteWho(rs[0].Actor, rs[0].Entry, entriesOf(w.log(), s.ID), handle)
					break find
				}
			}
			notes = append(notes, note)
		}
	}
	sortNotes(notes)
	return notes
}
