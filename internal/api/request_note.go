package api

// This file holds request.note: rhiza request note adds a note for the
// planner to a request still in the queue (no plan run drawn for it).

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
)

func init() {
	register(Operation{
		Name:     "request.note",
		Summary:  "Add a note for the planner to a queued plan request of this clone",
		Request:  reflect.TypeFor[RequestNoteRequest](),
		Response: reflect.TypeFor[RequestNoteResponse](),
		Run:      runRequestNote,
	})
}

// minRequestPrefix is the shortest prefix of a request ID request note
// takes.
const minRequestPrefix = 4

// RequestNoteRequest is a note to add to a queued request.
type RequestNoteRequest struct {
	Request string `json:"request" arg:"0" desc:"The request, by its full ID or a unique prefix of at least 4 characters; it must still be queued (no plan run drawn for it)"`
	Text    string `json:"text" arg:"1" desc:"The note for the request's planner: free text, its lines kept, leading and trailing blank space trimmed and CRLF turned into LF; an empty one, one over 8000 bytes or one not valid UTF-8 is refused"`
}

// RequestNoteResponse is the request with the note added, or in a dry run
// as it would be.
type RequestNoteResponse struct {
	Request RequestInfo `json:"request" desc:"The request, its notes oldest first, the new one last"`
	DryRun  bool        `json:"dry_run,omitempty" desc:"True if nothing was written"`
}

// Human writes the request and its notes as text.
func (r RequestNoteResponse) Human(w io.Writer) {
	verb := "Noted"
	if r.DryRun {
		verb = "Would note"
	}
	fmt.Fprintf(w, "%s %s request %s: %s\n", verb, r.Request.Kind, r.Request.ID, r.Request.Text)
	for i, n := range r.Request.Notes {
		fmt.Fprintf(w, "Note %d:\n", i+1)
		for _, line := range strings.Split(n, "\n") {
			fmt.Fprintf(w, "  %s\n", line)
		}
	}
}

func runRequestNote(ctx *Ctx, req any) (any, error) {
	r := req.(RequestNoteRequest)
	ref := strings.TrimSpace(r.Request)
	if ref == "" {
		return nil, Usagef("request note needs the request").WithNext(`rhiza request note <request> "<note>"`)
	}
	note, err := daemon.NormalizeRequestNote(r.Text)
	if err != nil {
		return nil, Usagef("the note %v", err)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	var q daemon.Request
	add := func(reqs []daemon.Request) ([]daemon.Request, error) {
		i, err := queuedRequest(reqs, ref)
		if err != nil {
			return nil, err
		}
		reqs[i].Notes = append(append([]string(nil), reqs[i].Notes...), note)
		q = reqs[i]
		return reqs, nil
	}
	if ctx.DryRun {
		reqs, err := daemon.ReadRequests(common)
		if err != nil {
			return nil, Errorf("read requests: %v", err)
		}
		if _, err := add(reqs); err != nil {
			return nil, err
		}
	} else if err := daemon.UpdateRequests(common, add); err != nil {
		return nil, requestsErr(err)
	}
	return RequestNoteResponse{Request: requestInfo(q, nil), DryRun: ctx.DryRun}, nil
}

// queuedRequest is the index in reqs of the request ref names, by its full
// ID or a unique prefix of at least minRequestPrefix characters. An unknown
// or ambiguous ref, a short prefix, or a request whose plan run has been
// drawn is a usage error.
func queuedRequest(reqs []daemon.Request, ref string) (int, error) {
	match := -1
	for i, q := range reqs {
		if q.ID == ref {
			match = i
			break
		}
	}
	if match < 0 {
		if len(ref) < minRequestPrefix {
			return 0, Usagef("no request %s (a prefix needs at least %d characters)", ref, minRequestPrefix).WithNext("rhiza request ls")
		}
		var ids []string
		for i, q := range reqs {
			if strings.HasPrefix(q.ID, ref) {
				match = i
				ids = append(ids, q.ID)
			}
		}
		switch {
		case len(ids) == 0:
			return 0, Usagef("no request %s", ref).WithNext("rhiza request ls")
		case len(ids) > 1:
			return 0, Usagef("request prefix %s is ambiguous: %s", ref, strings.Join(ids, ", ")).WithNext("rhiza request ls")
		}
	}
	if q := reqs[match]; q.Run != "" {
		return 0, Usagef("request %s is no longer queued: its plan run %s has been drawn, so its planner may already have read its brief", q.ID, q.Run).WithNext("rhiza request ls")
	}
	return match, nil
}
