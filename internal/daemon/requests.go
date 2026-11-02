package daemon

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// RequestsFile is the plan request queue's name, under Dir: one per clone,
// shared by its worktrees, never committed.
const RequestsFile = "requests.json"

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
// change set) is read from that run's journal.
type Request struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Text      string    `json:"text"`
	At        time.Time `json:"at"`
	Run       string    `json:"run,omitempty"`
	StartedAt time.Time `json:"started_at,omitzero"`
}

// Intent is the one-line intent the request's plan run opens its change set
// with: a bug's text prefixed "Fix: ".
func (r Request) Intent() string {
	if r.Kind == KindBug {
		return "Fix: " + r.Text
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

// WriteRequests replaces the queue, whole.
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
