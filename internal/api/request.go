package api

// This file holds the plan request queue's operations: rhiza request queues
// a feature, bug or milestone to plan in the clone (.git/rhiza/requests.json),
// and rhiza request ls shows each request's state. The daemon turns the
// oldest queued request into a plan run (daemonPlanRequests).

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "request",
		Summary:  "Queue a plan request in this clone: the daemon starts a plan run for it",
		Request:  reflect.TypeFor[RequestRequest](),
		Response: reflect.TypeFor[RequestResponse](),
		Run:      runRequest,
	})
	register(Operation{
		Name:     "request.ls",
		Summary:  "List this clone's plan requests, oldest first, each with its state, run and planned change set",
		Request:  reflect.TypeFor[RequestLsRequest](),
		Response: reflect.TypeFor[RequestLsResponse](),
		Run:      runRequestLs,
	})
}

// RequestRequest is a plan request to queue.
type RequestRequest struct {
	Text string `json:"text" arg:"0" desc:"What to plan, in one line"`
	Kind string `json:"kind,omitempty" enum:"feature,bug,milestone" desc:"feature, bug or milestone; default feature. A bug's plan run opens its change set with the intent prefixed Fix:, a milestone's a record change set"`
}

// RequestInfo is one plan request.
type RequestInfo struct {
	ID     string `json:"id" desc:"The request's ID"`
	Kind   string `json:"kind" desc:"feature, bug or milestone"`
	Text   string `json:"text" desc:"What to plan"`
	At     string `json:"at" desc:"When it was queued, RFC 3339 UTC"`
	State  string `json:"state" desc:"queued, planning (its plan run started), planned (its run succeeded) or failed (its run failed)"`
	Run    string `json:"run,omitempty" desc:"Its plan run; absent while queued"`
	Change string `json:"change,omitempty" desc:"The change set its plan run plans in; absent until the run's journal names it"`
}

// RequestResponse is the request queued, or in a dry run the one that
// would be.
type RequestResponse struct {
	Request RequestInfo `json:"request" desc:"The request"`
	DryRun  bool        `json:"dry_run,omitempty" desc:"True if nothing was written"`
}

// Human writes the response as text.
func (r RequestResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Queued %s request %s: %s\n", r.Request.Kind, r.Request.ID, r.Request.Text)
}

// RequestLsRequest takes nothing.
type RequestLsRequest struct{}

// RequestLsResponse is every plan request of the clone, oldest first.
type RequestLsResponse struct {
	Requests []RequestInfo `json:"requests" desc:"The requests, oldest first"`
}

// Human writes the list as text.
func (r RequestLsResponse) Human(w io.Writer) {
	if len(r.Requests) == 0 {
		fmt.Fprintln(w, "No plan requests")
		return
	}
	for _, q := range r.Requests {
		fmt.Fprintf(w, "- %s %s %s: %s", q.ID, q.Kind, q.State, q.Text)
		if q.Run != "" {
			fmt.Fprintf(w, " (run %s", q.Run)
			if q.Change != "" {
				fmt.Fprintf(w, ", change set %s", q.Change)
			}
			fmt.Fprint(w, ")")
		}
		fmt.Fprintln(w)
	}
}

// requestKinds are the kinds a request takes.
var requestKinds = []string{daemon.KindFeature, daemon.KindBug, daemon.KindMilestone}

func runRequest(ctx *Ctx, req any) (any, error) {
	r := req.(RequestRequest)
	text := strings.Join(strings.Fields(r.Text), " ")
	if text == "" {
		return nil, Usagef("a request needs its text").WithNext(`rhiza request "<what to plan>"`)
	}
	kind := strings.TrimSpace(r.Kind)
	if kind == "" {
		kind = daemon.KindFeature
	}
	if !slices.Contains(requestKinds, kind) {
		return nil, Usagef("--kind %q: one of %s", r.Kind, strings.Join(requestKinds, ", "))
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	reqs, err := daemon.ReadRequests(common)
	if err != nil {
		return nil, Errorf("read requests: %v", err)
	}
	rid, err := newRequestID(ctx, reqs)
	if err != nil {
		return nil, err
	}
	q := daemon.Request{ID: rid, Kind: kind, Text: text, At: ctx.Now().UTC()}
	if !ctx.DryRun {
		if err := daemon.WriteRequests(common, append(reqs, q)); err != nil {
			return nil, Errorf("write requests: %v", err)
		}
	}
	return RequestResponse{Request: requestInfo(q, nil), DryRun: ctx.DryRun}, nil
}

// newRequestID draws an ID no request of reqs has.
func newRequestID(ctx *Ctx, reqs []daemon.Request) (string, error) {
	for range maxRunIDDraws {
		u, err := ctx.IDs.New()
		if err != nil {
			return "", Errorf("%s", err.Error())
		}
		if !slices.ContainsFunc(reqs, func(q daemon.Request) bool { return q.ID == u.String() }) {
			return u.String(), nil
		}
	}
	return "", Errorf("no free request ID after %d draws", maxRunIDDraws)
}

func runRequestLs(ctx *Ctx, req any) (any, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	reqs, err := daemon.ReadRequests(common)
	if err != nil {
		return nil, Errorf("read requests: %v", err)
	}
	runs, err := run.NewStore(run.Dir(common)).List()
	if err != nil {
		return nil, Errorf("read runs: %v", err)
	}
	resp := RequestLsResponse{Requests: []RequestInfo{}}
	for _, q := range reqs {
		resp.Requests = append(resp.Requests, requestInfo(q, runs))
	}
	return resp, nil
}

func requestInfo(q daemon.Request, runs []run.Journal) RequestInfo {
	state, change := daemon.RequestState(q, runs)
	return RequestInfo{ID: q.ID, Kind: q.Kind, Text: q.Text, At: q.At.UTC().Format(time.RFC3339), State: state, Run: q.Run, Change: change}
}

// daemonPlanRequests starts a plan run, in a new Rex tab, for each of the
// oldest queued requests while fewer than max runs run: running counts the
// running runs and the starts not journaled yet. A request is marked with
// the run ID it draws before its tab opens, and run start --request takes
// that ID; a tab that fails to open puts the request back in the queue. The
// run is autonomous when plan autonomy is auto. It returns the requests
// started.
func daemonPlanRequests(ctx *Ctx, rs *run.Store, runs []run.Journal, running, max int, mode string) ([]string, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	reqs, err := daemon.ReadRequests(common)
	if err != nil {
		return nil, Errorf("read requests: %v", err)
	}
	for _, q := range reqs {
		// A started request whose run has not journaled yet takes a slot.
		if q.Run == "" || slices.ContainsFunc(runs, func(j run.Journal) bool { return j.ID == q.Run }) {
			continue
		}
		if ctx.Now().Sub(q.StartedAt) < daemonStartGrace {
			running++
		}
	}
	planMode, err := autonomyStage(ctx, "plan")
	if err != nil {
		return nil, err
	}
	var started []string
	for i := range reqs {
		if running >= max {
			break
		}
		if reqs[i].Run != "" {
			continue
		}
		rid, err := newRequestRunID(ctx, rs, reqs)
		if err != nil {
			return started, err
		}
		reqs[i].Run, reqs[i].StartedAt = rid, ctx.Now().UTC()
		if err := daemon.WriteRequests(common, reqs); err != nil {
			return started, Errorf("write requests: %v", err)
		}
		if err := daemonStartPlan(ctx, reqs[i], planMode == AutonomyAuto, mode); err != nil {
			reqs[i].Run, reqs[i].StartedAt = "", time.Time{}
			if werr := daemon.WriteRequests(common, reqs); werr != nil {
				return started, Errorf("write requests: %v", werr)
			}
			return started, err
		}
		started = append(started, reqs[i].ID)
		running++
	}
	return started, nil
}

// newRequestRunID draws a run ID that no run has and no request of reqs
// holds for its run yet.
func newRequestRunID(ctx *Ctx, rs *run.Store, reqs []daemon.Request) (string, error) {
	for range maxRunIDDraws {
		rid, err := newRunID(ctx, rs)
		if err != nil {
			return "", err
		}
		if !slices.ContainsFunc(reqs, func(q daemon.Request) bool { return q.Run == rid || q.ID == rid }) {
			return rid, nil
		}
	}
	return "", Errorf("no free run ID after %d draws", maxRunIDDraws)
}

// daemonStartPlan opens a new Rex tab running the plan run of request q:
// run start --workflow plan --new with its intent, --milestone for a
// milestone, and --request so the run takes the ID the request names.
func daemonStartPlan(ctx *Ctx, q daemon.Request, autonomous bool, mode string) error {
	if ctx.Exec == nil || ctx.Self == "" {
		return Errorf("daemon cannot start the plan run of request %s: no rex or rhiza path", q.ID)
	}
	args := []string{"window", "new"}
	if ctx.RexSession != "" {
		args = append(args, "-s", ctx.RexSession)
	}
	args = append(args, "plan-"+q.ID[:min(len(q.ID), 9)], "-c", ctx.Root, "--json", "--",
		ctx.Self, "run", "start", "--workflow", planWorkflow, "--new", q.Intent(), "--request", q.ID, "--permission-mode", mode)
	if q.Kind == daemon.KindMilestone {
		args = append(args, "--milestone")
	}
	if autonomous {
		args = append(args, "--autonomous")
	}
	if _, err := ctx.Exec("rex", args, ctx.Root, nil); err != nil {
		return Errorf("rex window new for request %s: %v", q.ID, err)
	}
	return nil
}

// requestRun is the run ID run start --request <id> takes: the one the
// daemon drew for the request. A request that is unknown, has no run yet,
// or whose run already journaled is a usage error.
func requestRun(common, reqID string, rs *run.Store) (string, error) {
	reqs, err := daemon.ReadRequests(common)
	if err != nil {
		return "", Errorf("read requests: %v", err)
	}
	for _, q := range reqs {
		if q.ID != reqID {
			continue
		}
		if q.Run == "" {
			return "", Usagef("request %s is queued: the daemon starts its plan run", reqID).WithNext("rhiza request ls")
		}
		_, err := rs.Read(q.Run)
		var nf *run.NotFoundError
		switch {
		case err == nil:
			return "", Usagef("request %s already has run %s", reqID, q.Run).WithNext("rhiza request ls")
		case !errors.As(err, &nf):
			return "", Errorf("read run %s: %v", q.Run, err)
		}
		return q.Run, nil
	}
	return "", Usagef("no request %s", reqID).WithNext("rhiza request ls")
}
