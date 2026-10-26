package api

// This file holds the session operations: session.start launches an agent
// for a change set in its own worktree and terminal tab, and session.ls and
// session.show report the sessions Rhiza has recorded (session_run.go holds
// session.run and session.hook, which keep the record's status). Sessions live in
// <git common dir>/rhiza/sessions, shared by every worktree and never
// committed; ls and show read only those files, never the terminal.
//
// Nothing here names an agent: the agent is an adapter.Adapter, the default
// one unless the request names another.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/dphaener/rhiza/internal/adapter"
	_ "github.com/dphaener/rhiza/internal/adapter/claudecode"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
)

func init() {
	register(Operation{
		Name:     "session.start",
		Summary:  "Launch an agent for a change set in its own worktree and terminal tab",
		Write:    true,
		Request:  reflect.TypeFor[SessionStartRequest](),
		Response: reflect.TypeFor[SessionStartResponse](),
		Run:      runSessionStart,
	})
	register(Operation{
		Name:     "session.ls",
		Summary:  "List the agent sessions with their change set, tab, worktree and status",
		Request:  reflect.TypeFor[SessionLsRequest](),
		Response: reflect.TypeFor[SessionLsResponse](),
		Run:      runSessionLs,
	})
	register(Operation{
		Name:     "session.show",
		Summary:  "Show one agent session with every status report",
		Request:  reflect.TypeFor[SessionShowRequest](),
		Response: reflect.TypeFor[SessionInfo](),
		Run:      runSessionShow,
	})
}

// Next commands of the session operations.
const (
	sessionStartNext = "rhiza describe session.start"
	sessionLsNext    = "rhiza session ls"
)

// defaultStep is the step a session works on unless the request names one.
const defaultStep = "build"

var (
	stepPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	runPattern  = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)
)

// Longest step and run the request accepts.
const (
	maxStep = 30
	maxRun  = 64
)

// SessionStartRequest launches an agent for a change set.
type SessionStartRequest struct {
	Change         string `json:"change,omitempty" arg:"0" desc:"The change set: a full ID or a unique prefix of 4+ characters; default --change, else the current one"`
	PermissionMode string `json:"permission_mode" desc:"The agent's permission mode, always explicit: one the agent accepts; modes that bypass its permission checks are refused"`
	Agent          string `json:"agent,omitempty" desc:"The agent adapter to launch; default the built-in default adapter"`
	Step           string `json:"step,omitempty" desc:"The step the session works on, lower-case words joined by dashes; default build, which the tab title leaves out"`
	Run            string `json:"run,omitempty" desc:"The run the launch belongs to, if any"`
	Prompt         string `json:"prompt,omitempty" desc:"The agent's first prompt; default none"`
	RexSession     string `json:"rex_session,omitempty" desc:"The Rex session to open the tab in; default the caller's REX_SESSION, else Rex chooses"`
}

// SessionStartResponse is the session launched, or in a dry run the one
// that would be.
type SessionStartResponse struct {
	Session        SessionInfo `json:"session" desc:"The session as recorded"`
	RexArgv        []string    `json:"rex_argv" desc:"The Rex command that opens the tab"`
	WorktreeReused bool        `json:"worktree_reused" desc:"Whether the worktree already existed on the branch"`
	DryRun         bool        `json:"dry_run,omitempty" desc:"True if nothing was written or run"`
}

// SessionLsRequest takes nothing.
type SessionLsRequest struct{}

// SessionLsResponse lists every recorded session.
type SessionLsResponse struct {
	Sessions []SessionSummary `json:"sessions" desc:"Every session, oldest first"`
}

// SessionSummary is one line of session.ls.
type SessionSummary struct {
	ID             string     `json:"id" desc:"The session ID, also the agent's own session ID"`
	Agent          string     `json:"agent" desc:"The agent adapter"`
	Change         string     `json:"change" desc:"The change set it works on"`
	Title          string     `json:"title" desc:"The change set's title"`
	Step           string     `json:"step" desc:"The step it works on"`
	Tab            SessionTab `json:"tab" desc:"The terminal tab it runs in"`
	Worktree       string     `json:"worktree" desc:"The worktree it runs in"`
	Branch         string     `json:"branch" desc:"The worktree's branch"`
	PermissionMode string     `json:"permission_mode" desc:"The agent's permission mode"`
	Status         string     `json:"status" enum:"starting,working,idle,blocked,done,error,exited" desc:"The resolved status"`
	Source         string     `json:"source" enum:"hook,rex,exit,default" desc:"Who reported the resolved status: hook, rex, exit, or default when no one has"`
	StartedAt      string     `json:"started_at" desc:"When it was launched, RFC 3339 UTC"`
}

// SessionShowRequest names one session.
type SessionShowRequest struct {
	ID string `json:"id" arg:"0" desc:"The session: its full ID or a unique prefix of 4+ characters"`
}

// SessionInfo is one session with every report and the resolved status.
type SessionInfo struct {
	ID             string        `json:"id" desc:"The session ID, also the agent's own session ID"`
	Agent          string        `json:"agent" desc:"The agent adapter"`
	Change         string        `json:"change" desc:"The change set it works on"`
	Title          string        `json:"title" desc:"The change set's title"`
	Step           string        `json:"step" desc:"The step it works on"`
	Run            string        `json:"run,omitempty" desc:"The run the launch belongs to, if any"`
	Actor          string        `json:"actor" desc:"The agent's actor handle in the worktree's model"`
	For            string        `json:"for" desc:"The actor the agent acts for: the launcher"`
	PermissionMode string        `json:"permission_mode" desc:"The agent's permission mode"`
	Prompt         string        `json:"prompt,omitempty" desc:"The agent's first prompt, if any"`
	Worktree       string        `json:"worktree" desc:"The worktree it runs in"`
	Branch         string        `json:"branch" desc:"The worktree's branch"`
	Tab            SessionTab    `json:"tab" desc:"The terminal tab it runs in"`
	StartedAt      string        `json:"started_at" desc:"When it was launched, RFC 3339 UTC"`
	Reports        SessionStatus `json:"reports" desc:"The latest report of each source that has spoken"`
	Status         SessionReport `json:"status" desc:"The resolved status: hooks first, then Rex, then the process exit"`
	Source         string        `json:"source" enum:"hook,rex,exit,default" desc:"Who reported the resolved status: hook, rex, exit, or default when no one has"`
}

// SessionTab is the terminal tab of a session.
type SessionTab struct {
	Session string `json:"session" desc:"The Rex session; empty when Rex chose it"`
	Window  string `json:"window" desc:"The Rex window ID; empty until the tab opened"`
	Block   string `json:"block" desc:"The Rex block running the agent; empty until the tab opened"`
	Title   string `json:"title" desc:"The tab title"`
}

// SessionStatus holds the latest report of each source.
type SessionStatus struct {
	Hook *SessionReport `json:"hook,omitempty" desc:"The agent's hooks"`
	Rex  *SessionReport `json:"rex,omitempty" desc:"Rex's status for the tab"`
	Exit *SessionReport `json:"exit,omitempty" desc:"The process exit, or a launch failure"`
}

// SessionReport is one status report.
type SessionReport struct {
	State  string `json:"state" enum:"starting,working,idle,blocked,done,error,exited" desc:"The state"`
	Kind   string `json:"kind,omitempty" desc:"What a blocked session waits for: permission, question or auth"`
	Detail string `json:"detail,omitempty" desc:"Free text from the source"`
	At     string `json:"at" desc:"When it was reported, RFC 3339 UTC"`
	Code   *int   `json:"code,omitempty" desc:"The process exit code; only the exit source sets it"`
}

// runSessionStart checks the request, then: resolves the change set and its
// title; derives slug, branch, worktree and tab title; creates or reuses the
// worktree and makes the change set current there; records the agent's
// actor in the worktree's model; writes the session record; and opens one
// Rex tab running `rhiza session run <id>`. A dry run derives everything
// and writes and runs nothing.
func runSessionStart(ctx *Ctx, req any) (any, error) {
	r := req.(SessionStartRequest)
	a, err := startAdapter(r)
	if err != nil {
		return nil, err
	}
	step, err := startStep(r)
	if err != nil {
		return nil, err
	}
	t, err := resolveChangeWorktree(ctx, r.Change, "no agent can work on it")
	if err != nil {
		return nil, err
	}
	w, cs, title, repo, worktree, branch := t.w, t.cs, t.title, t.repo, t.worktree, t.branch
	sid, err := newSessionID(ctx.IDs)
	if err != nil {
		return nil, err
	}
	rexSession := r.RexSession
	if rexSession == "" {
		rexSession = ctx.RexSession
	}
	tab := session.Tab{Session: rexSession, Title: session.TabTitle(repo, title, step)}
	reuse := t.reuse
	rec := session.Record{
		ID: sid, Agent: a.Name(), Change: cs.String(), Title: title, Step: step, Run: r.Run,
		Actor: handleSlug(a.Name() + "-" + sid[:8]), PermissionMode: r.PermissionMode, Prompt: r.Prompt,
		Worktree: worktree, Branch: branch, Tab: tab, StartedAt: ctx.Now().UTC().Truncate(time.Second),
	}
	self := ctx.Self
	if self == "" {
		self = "rhiza"
	}
	rexArgs := rexWindowArgs(tab, worktree, self, sid)
	if ctx.DryRun {
		return SessionStartResponse{Session: sessionInfo(rec), RexArgv: append([]string{"rex"}, rexArgs...),
			WorktreeReused: reuse, DryRun: true}, nil
	}
	switch {
	case ctx.Self == "":
		return nil, Errorf("cannot launch a session: the path of the running rhiza is unknown")
	case ctx.Exec == nil:
		return nil, Errorf("cannot launch a session: this caller runs no programs")
	}

	current, err := w.current()
	if err != nil {
		return nil, err
	}
	launcher, forID, err := w.resolveActor(current, ctx)
	if err != nil {
		return nil, err
	}
	ww, err := t.open(ctx)
	if err != nil {
		return nil, err
	}
	if rec.Actor, err = recordAgentActor(ww, a.Name(), sid, cs, launcher, forID); err != nil {
		return nil, err
	}
	rec.For = launcher.String()

	st := session.NewStore(session.Dir(t.common))
	if err := st.Write(rec); err != nil {
		return nil, Errorf("write session %s: %v", sid, err)
	}
	out, err := ctx.Exec("rex", rexArgs, worktree, nil)
	if err == nil {
		tab, err = parseRexWindow(tab, out)
	}
	if err != nil {
		msg := err.Error()
		if _, uerr := st.Update(sid, func(r *session.Record) {
			r.Status.Exit = &session.Report{State: session.StateError, Detail: msg, At: ctx.Now().UTC()}
		}); uerr != nil {
			msg += "; and recording that failed: " + uerr.Error()
		}
		return nil, Errorf("rex could not open the session's tab: %s", msg).WithNext("rhiza session show " + sid)
	}
	if rec, err = st.Update(sid, func(r *session.Record) { r.Tab = tab }); err != nil {
		return nil, Errorf("record the tab of session %s: %v", sid, err)
	}
	return SessionStartResponse{Session: sessionInfo(rec), RexArgv: append([]string{"rex"}, rexArgs...),
		WorktreeReused: reuse}, nil
}

// changeWorktree is a change set and the worktree an agent or a run works
// on it in, derived from its title: <repo>.<slug> on cs/<slug>, beside the
// main checkout.
type changeWorktree struct {
	w        *workspace
	cs       id.ID
	title    string
	common   string // the absolute git common dir
	prefix   string // Root's folder within the git repository, "" or ending in "/"
	repo     string // the main checkout's folder name
	branch   string
	worktree string
	reuse    bool // the worktree already exists on branch
}

// resolveChangeWorktree resolves the change set ref names (default
// --change, else the current one), refuses one in a terminal state with
// "change set X is <state>; <refusal>", and derives its worktree. It
// writes nothing.
func resolveChangeWorktree(ctx *Ctx, ref, refusal string) (changeWorktree, error) {
	var t changeWorktree
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return t, err
	}
	var cs id.ID
	if ref = strings.TrimSpace(ref); ref != "" {
		cs, err = w.resolveChangeSet(ref)
	} else {
		cs, err = w.effectiveChangeSet(ctx)
	}
	if err != nil {
		return t, err
	}
	if cs.IsZero() {
		return t, Usagef("no change set: name one or make one current").WithNext(changeLsNext)
	}
	c, _ := w.changeSet(cs)
	if c.State.Terminal() {
		return t, Conflictf("change set %s is %s; %s", cs, c.State, refusal).WithNext(changeLsNext)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return t, err
	}
	prefix, err := ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return t, Errorf("%v", err)
	}
	t = changeWorktree{w: w, cs: cs, title: c.Open.DisplayTitle(), common: common, prefix: strings.TrimSpace(prefix)}
	mainRoot := filepath.Dir(common)
	t.repo = filepath.Base(mainRoot)
	slug := session.Slug(t.title)
	t.branch = session.Branch(slug)
	t.worktree = session.Worktree(mainRoot, t.repo, slug)
	if t.reuse, err = worktreeReuse(ctx, t.worktree, t.branch); err != nil {
		return t, err
	}
	return t, nil
}

// open creates the worktree unless it is reused, opens its Rhiza
// repository with the change set in its log, and makes the change set
// current there.
func (t changeWorktree) open(ctx *Ctx) (*workspace, error) {
	if !t.reuse {
		if err := addWorktree(ctx, t.prefix, t.worktree, t.branch); err != nil {
			return nil, err
		}
	}
	ww, err := worktreeWorkspace(ctx, t.w, filepath.Join(t.worktree, t.prefix), t.cs)
	if err != nil {
		return nil, err
	}
	if err := ww.setCurrent(t.cs); err != nil {
		return nil, err
	}
	return ww, nil
}

// startAdapter returns the request's adapter after checking its permission
// mode: one is required, it must be the adapter's, and a mode that bypasses
// permission checks is refused, since no machine is trusted for it yet.
func startAdapter(r SessionStartRequest) (adapter.Adapter, error) {
	if r.PermissionMode == "" {
		return nil, Usagef("session start needs --permission-mode: an agent never starts without an explicit one").
			WithNext(sessionStartNext)
	}
	var a adapter.Adapter
	if r.Agent == "" {
		if a = adapter.Default(); a == nil {
			return nil, Errorf("no default agent adapter is registered")
		}
	} else {
		var ok bool
		if a, ok = adapter.Lookup(r.Agent); !ok {
			return nil, Usagef("unknown agent %q: want one of %s", r.Agent, strings.Join(adapter.Names(), ", ")).
				WithNext(sessionStartNext)
		}
	}
	if err := adapter.CheckMode(a, r.PermissionMode); err != nil {
		return nil, Usagef("%s: %v", a.Name(), err).WithNext(sessionStartNext)
	}
	if a.Bypass(r.PermissionMode) {
		return nil, bypassRefused(r.PermissionMode)
	}
	return a, nil
}

// bypassRefused is the error for a permission mode that bypasses the
// agent's permission checks: no machine is trusted for one yet.
func bypassRefused(mode string) *Error {
	return Errorf("permission mode %s bypasses the agent's permission checks: that needs a machine with trust = \"dedicated\", and no machine config declares one yet", mode).
		WithNext(sessionStartNext)
}

// startStep returns the request's step, default build, after checking it
// and the run.
func startStep(r SessionStartRequest) (string, error) {
	step := r.Step
	if step == "" {
		step = defaultStep
	}
	if len(step) > maxStep || !stepPattern.MatchString(step) {
		return "", Usagef("step %q: want lower-case words joined by single dashes, at most %d characters", step, maxStep).
			WithNext(sessionStartNext)
	}
	if r.Run != "" && (len(r.Run) > maxRun || !runPattern.MatchString(r.Run)) {
		return "", Usagef("run %q: want letters, digits, '.', '_', ':' or '-', at most %d characters", r.Run, maxRun).
			WithNext(sessionStartNext)
	}
	if strings.ContainsFunc(r.RexSession, unicode.IsControl) {
		return "", Usagef("rex_session %q has a control character", r.RexSession).WithNext(sessionStartNext)
	}
	return step, nil
}

// gitCommonDir is the absolute git common dir of ctx.Root: the main
// checkout's .git, whichever worktree Root is in.
func gitCommonDir(ctx *Ctx) (string, error) {
	if ctx.Git == nil {
		return "", Errorf("sessions need git: they live in the git common dir")
	}
	out, err := ctx.Git("rev-parse", "--git-common-dir")
	if err != nil {
		return "", Errorf("sessions need a git repository: %v", err)
	}
	d := strings.TrimSpace(out)
	if !filepath.IsAbs(d) {
		root, err := filepath.Abs(ctx.Root)
		if err != nil {
			return "", Errorf("project root: %v", err)
		}
		d = filepath.Join(root, d)
	}
	return filepath.Clean(d), nil
}

// newSessionID draws a UUID v4 from ids. The ID's random part comes first,
// so sessions started within the same minute differ in their leading hex,
// which names the agent's actor.
func newSessionID(ids *id.Generator) (string, error) {
	u, err := ids.New()
	if err != nil {
		return "", Errorf("%s", err.Error())
	}
	b := [16]byte(u)
	var e [16]byte
	copy(e[:10], b[6:])
	copy(e[10:], b[:6])
	return session.NewID(e), nil
}

// worktreeReuse reports whether path is already a worktree on branch, to be
// reused. A missing path is false. A path that exists but is not that
// worktree is a conflict.
func worktreeReuse(ctx *Ctx, path, branch string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, Errorf("worktree %s: %v", path, err)
	}
	out, err := ctx.Git("worktree", "list", "--porcelain")
	if err != nil {
		return false, Errorf("%v", err)
	}
	var wt, br string
	for _, line := range append(strings.Split(out, "\n"), "") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			wt = strings.TrimPrefix(line, "worktree ")
		case strings.HasPrefix(line, "branch "):
			br = strings.TrimPrefix(line, "branch ")
		case line == "":
			if wt != "" && samePath(wt, path) {
				if br == "refs/heads/"+branch {
					return true, nil
				}
				on := strings.TrimPrefix(br, "refs/heads/")
				if on == "" {
					on = "a detached HEAD"
				}
				return false, Conflictf("worktree %s is on %s, not %s", path, on, branch).WithNext("git worktree list")
			}
			wt, br = "", ""
		}
	}
	return false, Conflictf("%s exists but is not a worktree of this repository on %s; move it away first", path, branch).
		WithNext("git worktree list")
}

// samePath reports whether a and b name the same folder, symlinks resolved.
func samePath(a, b string) bool {
	if ra, err := filepath.EvalSymlinks(a); err == nil {
		a = ra
	}
	if rb, err := filepath.EvalSymlinks(b); err == nil {
		b = rb
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

// addWorktree creates the worktree path on branch: a new branch from HEAD,
// or the branch as it is if it already exists. HEAD must hold the Rhiza
// repository (prefix is Root's folder within the git repository), or the
// worktree would have no model to work on.
func addWorktree(ctx *Ctx, prefix, path, branch string) error {
	repoFile := prefix + store.Dir + "/" + store.RepoFile
	if _, err := ctx.Git("cat-file", "-e", "HEAD:"+repoFile); err != nil {
		return Conflictf("HEAD has no %s: commit the Rhiza repository so the worktree gets it", repoFile).
			WithNext("git add " + prefix + store.Dir + " && git commit")
	}
	var err error
	if _, berr := ctx.Git("show-ref", "--verify", "--quiet", "refs/heads/"+branch); berr == nil {
		_, err = ctx.Git("worktree", "add", path, branch)
	} else {
		_, err = ctx.Git("worktree", "add", "-b", branch, path)
	}
	if err != nil {
		return Errorf("create worktree %s: %v", path, err)
	}
	return nil
}

// worktreeWorkspace opens the Rhiza repository at root in the worktree and
// makes sure its log has change set cs, copying the entries from w's log
// that it lacks: a change set opened since HEAD's commit is not checked
// out. Entry files are immutable and named by ID, so a copy is the same
// file a later merge brings.
func worktreeWorkspace(ctx *Ctx, w *workspace, root string, cs id.ID) (*workspace, error) {
	ww, err := openWorkspace(root, ctx.Now, ctx.IDs, false)
	if err != nil {
		if e := AsError(err); e != nil && e.Code == CodeNotFound {
			return nil, Conflictf("the worktree has no Rhiza repository at %s", root).WithNext("git worktree list")
		}
		return nil, err
	}
	if _, ok := ww.changeSet(cs); ok {
		return ww, nil
	}
	src := filepath.Join(w.root, store.Dir, store.LogDir, cs.String())
	dst := filepath.Join(root, store.Dir, store.LogDir, cs.String())
	files, err := os.ReadDir(src)
	if err != nil {
		return nil, Errorf("read change set %s: %v", cs, err)
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return nil, Errorf("copy change set %s: %v", cs, err)
	}
	for _, f := range files {
		name := f.Name()
		if f.IsDir() || strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".json") {
			continue
		}
		if err := copyNewFile(filepath.Join(src, name), filepath.Join(dst, name)); err != nil {
			return nil, Errorf("copy change set %s: %v", cs, err)
		}
	}
	if _, ok := ww.changeSet(cs); !ok {
		return nil, Errorf("change set %s is not in the worktree's log at %s after copying it", cs, root)
	}
	return ww, nil
}

// copyNewFile copies src to dst unless dst exists, through a temp file and
// a rename.
func copyNewFile(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		return nil
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(dst), "."+filepath.Base(dst)+".*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	_, err = f.Write(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Chmod(tmp, 0o644)
	}
	if err == nil {
		err = os.Rename(tmp, dst)
	}
	return err
}

// recordAgentActor records the agent's actor node in the worktree's model,
// in a landed record change set written by the launcher, and returns its
// handle: <agent>-<first 8 hex of the session ID>, made unique.
func recordAgentActor(ww *workspace, agent, sid string, cs, launcher, forID id.ID) (string, error) {
	m := ww.model(id.ID{})
	handle := uniqueHandle(handleSlug(agent+"-"+sid[:8]), func(h string) bool {
		_, taken := m.HandleTaken(h)
		return taken
	})
	node, err := ww.ids.New()
	if err != nil {
		return "", Errorf("%s", err.Error())
	}
	opID, err := ww.ids.New()
	if err != nil {
		return "", Errorf("%s", err.Error())
	}
	fields, err := textFields(map[string]string{
		"name":    handle,
		"agent":   agent,
		"session": sid,
		"for":     launcher.String(),
		"role":    "runs change set " + cs.String(),
	})
	if err != nil {
		return "", err
	}
	c := op.Op{ID: opID, Type: op.TypeNodeCreate, Node: node, Kind: op.KindActor, Handle: handle, Fields: fields}
	if _, _, err := ww.writeRecord("Record agent actor "+handle+" for session "+sid, []op.Op{c}, launcher, forID); err != nil {
		return "", err
	}
	return handle, nil
}

// rexWindowArgs is the argv after "rex" that opens the session's tab in
// worktree, running `<self> session run <sid>`.
func rexWindowArgs(tab session.Tab, worktree, self, sid string) []string {
	args := []string{"window", "new"}
	if tab.Session != "" {
		args = append(args, "-s", tab.Session)
	}
	return append(args, tab.Title, "-c", worktree, "--json", "--", self, "session", "run", sid)
}

// rexCloseArgs is the argv after "rex" that closes the session's tab.
func rexCloseArgs(tab session.Tab) []string {
	args := []string{"window", "close"}
	if tab.Session != "" {
		args = append(args, "-s", tab.Session)
	}
	return append(args, tab.Window)
}

// closeSessionTab closes rec's tab in Rex, running rex in dir. It reports
// false, with no error, when the record has no tab or Rex no longer has
// it: a missing tab is already closed.
func closeSessionTab(ctx *Ctx, rec session.Record, dir string) (bool, error) {
	if rec.Tab.Window == "" {
		return false, nil
	}
	if ctx.Exec == nil {
		return false, Errorf("closing a Rex tab needs rex, and no programs can be run here")
	}
	if _, err := ctx.Exec("rex", rexCloseArgs(rec.Tab), dir, nil); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, Errorf("close tab %s: %v", rec.Tab.Window, err)
	}
	return true, nil
}

// parseRexWindow reads the window and first block from the JSON of
// `rex window new --json` into tab.
func parseRexWindow(tab session.Tab, out []byte) (session.Tab, error) {
	var win struct {
		WindowID string   `json:"window_id"`
		BlockIDs []string `json:"block_ids"`
	}
	if err := json.Unmarshal(out, &win); err != nil {
		return tab, fmt.Errorf("rex window new: unreadable JSON: %v", err)
	}
	if win.WindowID == "" || len(win.BlockIDs) == 0 || win.BlockIDs[0] == "" {
		return tab, fmt.Errorf("rex window new: no window_id or block_ids in %s", strings.TrimSpace(string(out)))
	}
	tab.Window, tab.Block = win.WindowID, win.BlockIDs[0]
	return tab, nil
}

// sessionStore opens the session store of ctx's git repository.
func sessionStore(ctx *Ctx) (*session.Store, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	return session.NewStore(session.Dir(common)), nil
}

// runSessionLs lists the recorded sessions, reading only their files.
func runSessionLs(ctx *Ctx, _ any) (any, error) {
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	recs, err := st.List()
	if err != nil {
		return nil, Errorf("read sessions: %v", err)
	}
	resp := SessionLsResponse{Sessions: []SessionSummary{}}
	for _, r := range recs {
		rep, src := session.Resolve(r)
		resp.Sessions = append(resp.Sessions, SessionSummary{
			ID: r.ID, Agent: r.Agent, Change: r.Change, Title: r.Title, Step: r.Step, Tab: sessionTab(r.Tab),
			Worktree: r.Worktree, Branch: r.Branch, PermissionMode: r.PermissionMode,
			Status: rep.State, Source: string(src), StartedAt: rfc3339(r.StartedAt),
		})
	}
	return resp, nil
}

// runSessionShow reports one session, reading only its file.
func runSessionShow(ctx *Ctx, req any) (any, error) {
	r := req.(SessionShowRequest)
	ref := strings.TrimSpace(r.ID)
	if ref == "" {
		return nil, Usagef("session show needs a session ID").WithNext(sessionLsNext)
	}
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := readSession(st, ref)
	if err != nil {
		return nil, err
	}
	return sessionInfo(rec), nil
}

// readSession reads the session ref names, a full ID or a unique prefix,
// mapping the store's errors to not_found and usage.
func readSession(st *session.Store, ref string) (session.Record, error) {
	rec, err := st.Read(ref)
	var nf *session.NotFoundError
	var amb *session.AmbiguousError
	switch {
	case errors.As(err, &nf):
		return rec, NotFoundf("no session matches %q", ref).WithNext(sessionLsNext)
	case errors.As(err, &amb):
		return rec, Usagef("%v", amb).WithNext(sessionLsNext)
	case err != nil:
		return rec, Errorf("read session %s: %v", ref, err)
	}
	return rec, nil
}

// sessionInfo is the response form of a record.
func sessionInfo(r session.Record) SessionInfo {
	rep, src := session.Resolve(r)
	return SessionInfo{
		ID: r.ID, Agent: r.Agent, Change: r.Change, Title: r.Title, Step: r.Step, Run: r.Run,
		Actor: r.Actor, For: r.For, PermissionMode: r.PermissionMode, Prompt: r.Prompt, Worktree: r.Worktree,
		Branch: r.Branch, Tab: sessionTab(r.Tab), StartedAt: rfc3339(r.StartedAt),
		Reports: SessionStatus{
			Hook: sessionReportPtr(r.Status.Hook),
			Rex:  sessionReportPtr(r.Status.Rex),
			Exit: sessionReportPtr(r.Status.Exit),
		},
		Status: sessionReport(rep),
		Source: string(src),
	}
}

func sessionTab(t session.Tab) SessionTab {
	return SessionTab{Session: t.Session, Window: t.Window, Block: t.Block, Title: t.Title}
}

func sessionReport(r session.Report) SessionReport {
	return SessionReport{State: r.State, Kind: r.Kind, Detail: r.Detail, At: rfc3339(r.At), Code: r.Code}
}

func sessionReportPtr(r *session.Report) *SessionReport {
	if r == nil {
		return nil
	}
	s := sessionReport(*r)
	return &s
}

// rfc3339 formats t in UTC to the second.
func rfc3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// Human reports the session launched, or the plan of a dry run.
func (r SessionStartResponse) Human(w io.Writer) {
	s := r.Session
	worktree := "created"
	if r.WorktreeReused {
		worktree = "reused"
	}
	if r.DryRun {
		fmt.Fprintf(w, "Would start session %s for change set %s: %s\n\n", s.ID, s.Change, s.Title)
		if !r.WorktreeReused {
			worktree = "to create"
		}
	} else {
		fmt.Fprintf(w, "Started session %s for change set %s: %s\n\n", s.ID, s.Change, s.Title)
	}
	fmt.Fprintf(w, "- Worktree: %s on %s (%s)\n", s.Worktree, s.Branch, worktree)
	fmt.Fprintf(w, "- Tab: %s\n", tabLine(s.Tab))
	fmt.Fprintf(w, "- Agent: %s, permission mode %s\n", s.Agent, s.PermissionMode)
	if r.DryRun {
		fmt.Fprintf(w, "- Command: %s\n", strings.Join(r.RexArgv, " "))
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written or run.\n")
		return
	}
	fmt.Fprintf(w, "- Actor: %s\n", s.Actor)
	fmt.Fprintf(w, "- Status: %s\n", s.Status.State)
	fmt.Fprintf(w, "\nnext: rhiza session show %s\n", s.ID)
}

// Human lists the sessions, one per line.
func (r SessionLsResponse) Human(w io.Writer) {
	if len(r.Sessions) == 0 {
		fmt.Fprintf(w, "No sessions yet.\n\nnext: rhiza session start --permission-mode <mode>\n")
		return
	}
	for _, s := range r.Sessions {
		fmt.Fprintf(w, "- %s %s (%s) %s: %s, %s, %s\n", s.ID, s.Status, s.Source, s.Change, s.Tab.Title, s.PermissionMode, s.Worktree)
	}
}

// Human shows the session and every report.
func (r SessionInfo) Human(w io.Writer) {
	fmt.Fprintf(w, "# Session %s\n\n", r.ID)
	fmt.Fprintf(w, "- Change set: %s: %s\n", r.Change, r.Title)
	fmt.Fprintf(w, "- Step: %s\n", r.Step)
	if r.Run != "" {
		fmt.Fprintf(w, "- Run: %s\n", r.Run)
	}
	fmt.Fprintf(w, "- Agent: %s, permission mode %s\n", r.Agent, r.PermissionMode)
	fmt.Fprintf(w, "- Actor: %s, for %s\n", r.Actor, r.For)
	fmt.Fprintf(w, "- Worktree: %s on %s\n", r.Worktree, r.Branch)
	fmt.Fprintf(w, "- Tab: %s\n", tabLine(r.Tab))
	fmt.Fprintf(w, "- Started: %s\n", r.StartedAt)
	fmt.Fprintf(w, "- Status: %s (from %s)\n", reportLine(r.Status), r.Source)
	fmt.Fprintf(w, "\n## Reports\n\n")
	for _, s := range []struct {
		name string
		rep  *SessionReport
	}{{"hook", r.Reports.Hook}, {"rex", r.Reports.Rex}, {"exit", r.Reports.Exit}} {
		if s.rep == nil {
			fmt.Fprintf(w, "- %s: none\n", s.name)
			continue
		}
		fmt.Fprintf(w, "- %s: %s\n", s.name, reportLine(*s.rep))
	}
}

// tabLine is a tab's title and, once open, its window and block.
func tabLine(t SessionTab) string {
	if t.Window == "" {
		return t.Title + " (not open)"
	}
	return fmt.Sprintf("%s (window %s, block %s)", t.Title, t.Window, t.Block)
}

// reportLine is a report as "<state>[ <kind>][ code <n>][: <detail>] at <at>".
func reportLine(r SessionReport) string {
	s := r.State
	if r.Kind != "" {
		s += " " + r.Kind
	}
	if r.Code != nil {
		s += fmt.Sprintf(" code %d", *r.Code)
	}
	if r.Detail != "" {
		s += ": " + r.Detail
	}
	return s + " at " + r.At
}
