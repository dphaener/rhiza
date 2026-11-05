package api

// This file holds note.add: a note is written into the current change set
// as a note node holding the text, linked discusses to the node or change
// set it is about, in one entry. With --op the note also names one
// operation of that change set.

import (
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "note.add",
		Summary:  "Add a note discussing a node or change set, optionally naming one of that change set's operations, into the current change set",
		Write:    true,
		Request:  reflect.TypeFor[NoteAddRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Run:      runNoteAdd,
	})
}

const noteAddNext = "rhiza describe note.add"

// NoteAddRequest names the target and the note's text.
type NoteAddRequest struct {
	Target string `json:"target" arg:"0" desc:"What the note discusses: a node (an ID, a unique ID prefix of 4+ characters, or a live node's handle) or a change set"`
	Text   string `json:"text" arg:"1" desc:"The note; must not be blank"`
	Op     string `json:"op,omitempty" flag:"op" desc:"An operation of the target change set the note is about (the ID change show lists, or the operation's own ID; a unique prefix of 4+ characters)"`
}

// noteHandle is the automatic handle of note node n: note-<its short ID>.
func noteHandle(n id.ID) string {
	return "note-" + handleSlug(shortID(n))
}

func runNoteAdd(ctx *Ctx, req any) (any, error) {
	r := req.(NoteAddRequest)
	ref := strings.TrimSpace(r.Target)
	if ref == "" {
		return nil, Usagef("note add needs a target").WithNext(noteAddNext)
	}
	text := strings.TrimSpace(r.Text)
	if text == "" {
		return nil, Usagef("the note is blank: say something").WithNext(noteAddNext)
	}
	opRef := strings.TrimSpace(r.Op)
	return runWriteOn(ctx, "note.add", commandLine("note add", ref, text), func(w *workspace, cs id.ID) (edit, error) {
		m := w.model(cs)
		var to string
		if n, err := m.Resolve(ref); err == nil {
			if opRef != "" {
				return edit{}, Usagef("--op names an operation of a change set, but %s is a node", ref).WithNext(noteAddNext)
			}
			to = n.ID.String()
		} else if e := AsError(err); e == nil || e.Code != CodeNotFound {
			return edit{}, err
		} else if target, csErr := w.resolveChangeSet(ref); csErr == nil {
			to = target.String()
			if opRef != "" {
				full, err := changeSetOp(w, target, opRef)
				if err != nil {
					return edit{}, err
				}
				opRef = full
			}
		} else {
			return edit{}, Usagef("no node or change set %q", ref).WithNext("rhiza ls")
		}
		var handle string
		fields := []string{"statement=" + text}
		if opRef != "" {
			fields = append(fields, "op="+opRef)
		}
		items := []applyItem{
			{op: "add", build: func(b writeBase) (edit, error) {
				return buildAdd(b, AddRequest{Kind: string(op.KindNote), Fields: fields, lenient: true, handleFor: func(n id.ID) string {
					handle = noteHandle(n)
					return handle
				}})
			}},
			{op: "link", build: func(b writeBase) (edit, error) {
				return buildLink(b, LinkRequest{From: handle, Type: string(op.LinkDiscusses), To: to})
			}},
		}
		return buildApply(w, cs, ctx.IDs, items)
	})
}

// changeSetOp resolves ref to the full ID of one of change set cs's
// operations: the ID of the operation or of what it acts on, or a unique prefix of 4+ characters.
func changeSetOp(w *workspace, cs id.ID, ref string) (string, error) {
	c, ok := w.changeSet(cs)
	if !ok {
		return "", Usagef("no change set %s", cs)
	}
	// change show prints an operation by the node or link it acts on, so
	// either that or the operation's own ID names it.
	var found []string
	for _, o := range c.Ops {
		for _, s := range []string{o.ID.String(), o.Target().String()} {
			if strings.EqualFold(s, ref) || (len(ref) >= 4 && strings.HasPrefix(s, strings.ToUpper(ref))) {
				found = append(found, o.ID.String())
				break
			}
		}
	}
	slices.Sort(found)
	found = slices.Compact(found)
	switch len(found) {
	case 0:
		return "", Usagef("--op %s is not an operation of change set %s", ref, shortID(cs)).WithNext("rhiza change show " + shortID(cs))
	case 1:
		return found[0], nil
	}
	return "", Usagef("--op %s matches several operations: %s", ref, strings.Join(found, ", "))
}
