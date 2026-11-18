package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dphaener/rhiza/internal/run"
)

// RequestsFile is the plan request queue's name, under Dir: one per clone,
// shared by its worktrees, never committed.
const RequestsFile = "requests.json"

// RequestsLockFile is the lock every writer of the queue holds around its
// read and write, under Dir, so no writer's update overwrites another's.
const RequestsLockFile = "requests.lock"

// The kinds of a plan request.
const (
	KindFeature   = "feature"
	KindBug       = "bug"
	KindMilestone = "milestone"
)

// The states of a plan request, derived by RequestState.
const (
	RequestQueued   = "queued"
	RequestPlanning = "planning"
	RequestPlanned  = "planned"
	RequestFailed   = "failed"
)

// Request is one plan request in the queue. Run is set when the daemon
// starts its plan run: the run's ID, drawn by the daemon and taken by run
// start --request. Everything else about it (planning, planned, failed, the
// change set) is read from that run's journal. RetryAt is set when a start
// of the daemon's never journaled and the request went back in the queue:
// the daemon does not start it again before then. Notes are free text the
// request carries for its planner, oldest first, each normalised by
// NormalizeRequestNote; a request without notes writes no notes key.
type Request struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Text      string `json:"text"`
	Milestone string `json:"milestone,omitempty"`
	// Workflow is the workflow the request named for its change set's ready
	// pick; the change set its plan run opens records it.
	Workflow  string    `json:"workflow,omitempty"`
	Ref       string    `json:"ref,omitempty"`
	At        time.Time `json:"at"`
	Run       string    `json:"run,omitempty"`
	StartedAt time.Time `json:"started_at,omitzero"`
	RetryAt   time.Time `json:"retry_at,omitzero"`
	Notes     []string  `json:"notes,omitempty"`
}

// MaxRequestNote is the most bytes a request note holds, once normalised.
const MaxRequestNote = 8000

// NormalizeRequestNote is note as a request stores it: CRLF turned into LF
// and leading and trailing blank space trimmed, its lines
// kept. A note that is not valid UTF-8, is empty or blank, or is over
// MaxRequestNote bytes is an error naming why.
func NormalizeRequestNote(note string) (string, error) {
	if !utf8.ValidString(note) {
		return "", errors.New("is not valid UTF-8")
	}
	note = strings.ReplaceAll(note, "\r\n", "\n")
	note = strings.TrimSpace(note)
	if note == "" {
		return "", errors.New("is empty")
	}
	if len(note) > MaxRequestNote {
		return "", fmt.Errorf("is %d bytes, at most %d", len(note), MaxRequestNote)
	}
	return note, nil
}

// Intent is the one-line intent the request's plan run opens its change set
// with: a bug's text prefixed "Fix " (no colon, so the change set's display
// title, its first clause, names the bug rather than being "Fix" for every
// bug, and two bugs planned at once get distinct branches and worktrees).
func (r Request) Intent() string {
	if r.Kind == KindBug {
		return "Fix " + r.Text
	}
	return r.Text
}

// requestsDoc is the file's shape.
type requestsDoc struct {
	Requests []Request `json:"requests"`
}

// ReadRequests is the queue, oldest first; none when there is no file.
func ReadRequests(common string) ([]Request, error) {
	b, err := os.ReadFile(filepath.Join(Dir(common), RequestsFile))
	if errors.Is(err, os.ErrNotExist) {
		return []Request{}, nil
	}
	if err != nil {
		return nil, err
	}
	var d requestsDoc
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	if d.Requests == nil {
		d.Requests = []Request{}
	}
	return d.Requests, nil
}

// UpdateRequests reads the queue, passes it to update and writes back what
// update returns, holding the queue's lock throughout (waiting for it while
// another writer holds it). update returning an error writes nothing.
func UpdateRequests(common string, update func([]Request) ([]Request, error)) error {
	if err := os.MkdirAll(Dir(common), 0o755); err != nil {
		return err
	}
	unlock, err := openWait(filepath.Join(Dir(common), RequestsLockFile))
	if err != nil {
		return err
	}
	defer unlock()
	reqs, err := ReadRequests(common)
	if err != nil {
		return err
	}
	reqs, err = update(reqs)
	if err != nil {
		return err
	}
	return WriteRequests(common, reqs)
}

// WriteRequests replaces the queue, whole. Writers that read the queue
// first go through UpdateRequests instead.
func WriteRequests(common string, reqs []Request) error {
	if err := os.MkdirAll(Dir(common), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(requestsDoc{Requests: reqs}, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(Dir(common), RequestsFile)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// RequestState is the request's state and the change set it planned, read
// from runs (every run of the repository): queued with no run; planning
// while its run has no journal yet or has not succeeded; planned, with the
// change set its run's journal names, once the run succeeded; failed when
// the run failed.
func RequestState(r Request, runs []run.Journal) (state, change string) {
	if r.Run == "" {
		return RequestQueued, ""
	}
	for _, j := range runs {
		if j.ID != r.Run {
			continue
		}
		st := j.State()
		switch {
		case st.Ended && st.Result == run.Failed:
			return RequestFailed, st.Change
		case st.Ended && st.Result == run.Succeeded:
			return RequestPlanned, st.Change
		default:
			return RequestPlanning, st.Change
		}
	}
	return RequestPlanning, ""
}
