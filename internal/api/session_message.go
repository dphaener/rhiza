package api

// This file holds session.message, which sends a running agent session a
// message: the text goes into a file inside the session's worktree (under
// the gitignored .rhiza/local/messages/<session>/), and a one-line pointer
// to that file is typed into the session's terminal tab. The pointer is
// typed only while the session waits at its prompt (resolved status idle),
// never over an open menu or question; otherwise the command waits.
//
// typePointer is the delivery shared with run.amend. Nothing here names an
// agent, and nothing is stored in the terminal: the record and the file
// are Rhiza's.

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
)

func init() {
	register(Operation{
		Name:     "session.message",
		Summary:  "Send a running agent session a message: a file in its worktree and a one-line pointer, typed at its prompt or, while it works, handed over at its next tool call",
		Write:    true,
		Request:  reflect.TypeFor[SessionMessageRequest](),
		Response: reflect.TypeFor[SessionMessageResponse](),
		Run:      runSessionMessage,
	})
}

const sessionMessageNext = "rhiza describe session.message"

// messagePoll is how often typePointer reads the session record while it
// waits for the prompt; tests lower it.
var messagePoll = 250 * time.Millisecond

// defaultMessageTimeout is how long a message waits for the prompt unless
// the request says otherwise.
const defaultMessageTimeout = 10 * time.Minute

// messagesDir is the folder, under the Rhiza repository's local folder,
// that holds each session's messages.
const messagesDir = "messages"

// maxMessageSeq is the highest message number a session's folder takes.
const maxMessageSeq = 9999

// SessionMessageRequest names the session and the message.
type SessionMessageRequest struct {
	Session string `json:"session" arg:"0" desc:"The session: its full ID or a unique prefix of 4+ characters"`
	Text    string `json:"text" arg:"1" desc:"The message, kept verbatim in the file the session reads; must not be blank"`
	Timeout string `json:"timeout,omitempty" desc:"How long to wait for the session to be at its prompt, a Go duration such as 30s or 10m; default 10m"`
}

// SessionMessageResponse is the message sent, or in a dry run the one that
// would be.
type SessionMessageResponse struct {
	Session string `json:"session" desc:"The session's full ID"`
	File    string `json:"file" desc:"The message file, an absolute path inside the session's worktree"`
	Pointer string `json:"pointer" desc:"The one line typed at the session's prompt"`
	From    string `json:"from" desc:"The sender's actor handle"`
	Waited  string `json:"waited" desc:"How long it waited for the prompt, a Go duration"`
	DryRun  bool   `json:"dry_run,omitempty" desc:"True if nothing was written or sent"`
}

// runSessionMessage checks the request and the session, resolves the
// sender, writes the message file and types the pointer once the session
// is at its prompt. A dry run does all but write and send.
func runSessionMessage(ctx *Ctx, req any) (any, error) {
	r := req.(SessionMessageRequest)
	ref := strings.TrimSpace(r.Session)
	if ref == "" {
		return nil, Usagef("session message needs a session ID").WithNext(sessionLsNext)
	}
	if strings.TrimSpace(r.Text) == "" {
		return nil, Usagef("the message is blank: say something").WithNext(sessionMessageNext)
	}
	timeout, err := messageTimeout(r.Timeout)
	if err != nil {
		return nil, err
	}
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := readSession(st, ref)
	if err != nil {
		return nil, err
	}
	if err := messageable(rec); err != nil {
		return nil, err
	}
	if !ctx.DryRun && ctx.Exec == nil {
		return nil, Errorf("cannot send a message: this caller runs no programs")
	}
	from, err := messageSender(ctx)
	if err != nil {
		return nil, err
	}
	dir, rel, err := messageFolder(ctx, rec)
	if err != nil {
		return nil, err
	}
	content := messageContent(from, rec, rfc3339(ctx.Now()), r.Text)
	resp := SessionMessageResponse{Session: rec.ID, From: from.Handle, Waited: roundWait(0)}
	if ctx.DryRun {
		name, err := nextMessageName(dir)
		if err != nil {
			return nil, err
		}
		resp.File = filepath.Join(dir, name)
		resp.Pointer = messagePointer(from.Handle, rel+"/"+name)
		resp.DryRun = true
		return resp, nil
	}
	name, err := writeMessage(dir, content)
	if err != nil {
		return nil, err
	}
	resp.File = filepath.Join(dir, name)
	resp.Pointer = messagePointer(from.Handle, rel+"/"+name)
	if queued, err := queuePointer(st, rec.ID, resp.Pointer); err != nil {
		return nil, withNote(err, "the message is in "+resp.File)
	} else if queued {
		return resp, nil
	}
	waited, err := typePointer(ctx, st, rec.ID, resp.Pointer, timeout)
	if err != nil {
		return nil, withNote(err, "the message is in "+resp.File)
	}
	resp.Waited = roundWait(waited)
	return resp, nil
}

// queuePointer keeps line in the session's record for its next tool-call
// hook to hand over, when the session is working, and reports whether it
// did. Any other status leaves it to typePointer.
func queuePointer(st *session.Store, sid, line string) (bool, error) {
	queued := false
	_, err := st.Update(sid, func(r *session.Record) {
		if r.Status.Exited() {
			return
		}
		if rep, _ := session.Resolve(*r); rep.State == session.StateWorking {
			r.Pending, queued = append(r.Pending, line), true
		}
	})
	if err != nil {
		return false, Errorf("record the message for session %s: %v", sid, err)
	}
	return queued, nil
}

// messageTimeout parses the request's timeout, default 10m.
func messageTimeout(v string) (time.Duration, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return defaultMessageTimeout, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, Usagef("timeout %q: want a positive Go duration such as 30s or 10m", v).WithNext(sessionMessageNext)
	}
	return d, nil
}

// messageable refuses a session that can take no message: one that has
// exited, or one whose tab never opened.
func messageable(rec session.Record) error {
	if rec.Status.Exited() {
		return Conflictf("session %s has exited; it can take no message", rec.ID).WithNext("rhiza session show " + rec.ID)
	}
	if rec.Tab.Block == "" {
		return Conflictf("session %s has no open tab", rec.ID).WithNext("rhiza session show " + rec.ID)
	}
	return nil
}

// sender is who sends a message or an amendment: an actor's handle and ID.
// ID is zero in a dry run whose actor is not recorded yet.
type sender struct {
	Handle string
	ID     id.ID
}

// messageSender resolves the caller's actor, as every write does, and its
// handle. A dry run whose actor is not recorded yet names the handle the
// real run would record.
func messageSender(ctx *Ctx) (sender, error) {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return sender{}, err
	}
	current, err := w.current()
	if err != nil {
		return sender{}, err
	}
	actor, _, err := w.resolveActor(current, ctx)
	if err != nil {
		return sender{}, err
	}
	m := w.model(current)
	if actor.IsZero() {
		local, _, _ := strings.Cut(gitValue(ctx.GitConfig, "user.email"), "@")
		handle := uniqueHandle(handleSlug(local), func(h string) bool {
			_, taken := m.HandleTaken(h)
			return taken
		})
		return sender{Handle: handle}, nil
	}
	n, ok := m.Node(actor)
	if !ok {
		return sender{}, Errorf("actor %s is not in the model", actor)
	}
	return sender{Handle: n.Handle, ID: actor}, nil
}

// messageFolder is the session's message folder, absolute, and the same
// folder relative to the session's worktree, slash-separated. Root's
// folder within the git repository is the same in every worktree.
func messageFolder(ctx *Ctx, rec session.Record) (dir, rel string, err error) {
	if ctx.Git == nil {
		return "", "", Errorf("messages need git: they live in the session's worktree")
	}
	prefix, err := ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return "", "", Errorf("%v", err)
	}
	rel = strings.TrimSpace(prefix) + store.Dir + "/" + localDir + "/" + messagesDir + "/" + rec.ID
	return filepath.Join(rec.Worktree, filepath.FromSlash(rel)), rel, nil
}

// messageContent is the message file: a heading and list naming the
// sender, the session and when, then the text verbatim, ending in one
// newline.
func messageContent(from sender, rec session.Record, at, text string) string {
	who := from.Handle + " (not yet recorded)"
	if !from.ID.IsZero() {
		who = from.Handle + " (" + from.ID.String() + ")"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Message from %s\n\n", from.Handle)
	fmt.Fprintf(&b, "- From: %s\n", who)
	fmt.Fprintf(&b, "- Session: %s\n", rec.ID)
	fmt.Fprintf(&b, "- Change set: %s\n", rec.Change)
	fmt.Fprintf(&b, "- Step: %s\n", rec.Step)
	fmt.Fprintf(&b, "- At: %s\n\n", at)
	b.WriteString(strings.TrimRight(text, "\n"))
	b.WriteString("\n")
	return b.String()
}

// messagePointer is the line typed at the session's prompt, naming the
// message file relative to the session's worktree. It is one line: every
// control character becomes a space.
func messagePointer(handle, rel string) string {
	return oneLine("rhiza: message from " + handle + ": read " + rel + " and act on it.")
}

// oneLine replaces every control character in s with a space.
func oneLine(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, s)
}

// messageName is the file name of message number n.
func messageName(n int) string { return fmt.Sprintf("%04d.md", n) }

// nextMessageName is the first free message name in dir, without claiming
// it: what a dry run reports.
func nextMessageName(dir string) (string, error) {
	for n := 1; n <= maxMessageSeq; n++ {
		_, err := os.Lstat(filepath.Join(dir, messageName(n)))
		if errors.Is(err, fs.ErrNotExist) {
			return messageName(n), nil
		}
		if err != nil {
			return "", Errorf("message folder %s: %v", dir, err)
		}
	}
	return "", Conflictf("message folder %s is full: it holds %d messages", dir, maxMessageSeq)
}

// writeMessage claims the next free message name in dir with an exclusive
// create, so two senders never share a file, and writes content to it.
func writeMessage(dir, content string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", Errorf("create message folder %s: %v", dir, err)
	}
	for n := 1; n <= maxMessageSeq; n++ {
		name := messageName(n)
		path := filepath.Join(dir, name)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if errors.Is(err, fs.ErrExist) {
			continue
		}
		if err != nil {
			return "", Errorf("write message %s: %v", path, err)
		}
		_, err = f.WriteString(content)
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return "", Errorf("write message %s: %v", path, err)
		}
		return name, nil
	}
	return "", Conflictf("message folder %s is full: it holds %d messages", dir, maxMessageSeq)
}

// typePointer types line into session sid's tab and presses Enter, once
// the session waits at its prompt: it reads the record every messagePoll
// and types only while the resolved status is idle. Any other status (a
// blocked session with a menu or question open, a working or starting
// one) waits. An exit, or timeout passing, is a conflict. It returns how
// long it waited.
func typePointer(ctx *Ctx, st *session.Store, sid, line string, timeout time.Duration) (time.Duration, error) {
	if ctx.Exec == nil {
		return 0, Errorf("cannot send a message: this caller runs no programs")
	}
	clock := runClock
	start := clock.Now()
	last, told := "", ""
	for {
		rec, err := st.Read(sid)
		if err != nil {
			return clock.Now().Sub(start), Errorf("read session %s: %v", sid, err)
		}
		if rec.Status.Exited() {
			return clock.Now().Sub(start), Conflictf("session %s exited before it reached its prompt", sid).
				WithNext("rhiza session show " + sid)
		}
		if rec.Tab.Block == "" {
			return clock.Now().Sub(start), Conflictf("session %s has no open tab", sid).WithNext("rhiza session show " + sid)
		}
		waited := clock.Now().Sub(start)
		rep, _ := session.Resolve(rec)
		if rep.State == session.StateIdle {
			return waited, sendPointer(ctx, rec, line)
		}
		last = waitingStatus(rep)
		if told != last {
			told = last
			if ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: waiting for session %s: %s; %s\n", sid, last, waitingHint(rep, sid, timeout))
			}
		}
		if waited >= timeout {
			return waited, Conflictf("session %s did not reach its prompt within %s (last status: %s)", sid, timeout, last).
				WithNext("rhiza session show " + sid)
		}
		clock.Sleep(min(messagePoll, timeout-waited))
	}
}

// waitingStatus describes a resolved status that is not idle, for the wait
// line and the timeout error.
func waitingStatus(rep session.Report) string {
	if rep.State == session.StateBlocked {
		s := "blocked on a " + rep.Kind
		if rep.Kind == "" {
			s = "blocked"
		}
		if rep.Detail != "" {
			s += ": " + rep.Detail
		}
		return s
	}
	if rep.State == "" {
		return "starting"
	}
	return rep.State
}

// waitingHint says what the wait is for and how to end it.
func waitingHint(rep session.Report, sid string, timeout time.Duration) string {
	if rep.State == session.StateBlocked {
		return fmt.Sprintf("answer it with rhiza inbox open (or rhiza session show %s); waiting up to %s", sid, timeout)
	}
	return fmt.Sprintf("waiting up to %s for it to reach its prompt", timeout)
}

// sendPointer types line into the session's block, then presses Enter.
func sendPointer(ctx *Ctx, rec session.Record, line string) error {
	block := rec.Tab.Block
	if _, err := ctx.Exec("rex", []string{"send", "-b", block, line}, rec.Worktree, nil); err != nil {
		return Errorf("rex could not type into session %s's tab: %v", rec.ID, err).WithNext("rhiza session show " + rec.ID)
	}
	if _, err := ctx.Exec("rex", []string{"send-key", "-b", block, "enter"}, rec.Worktree, nil); err != nil {
		return Errorf("rex could not press enter in session %s's tab: %v", rec.ID, err).WithNext("rhiza session show " + rec.ID)
	}
	return nil
}

// withNote appends "; <note>" to err's message, keeping its code and next.
func withNote(err error, note string) error {
	e := AsError(err)
	if e == nil {
		return err
	}
	out := *e
	out.Message += "; " + note
	return &out
}

// roundWait formats a wait rounded to 10ms.
func roundWait(d time.Duration) string {
	return d.Round(10 * time.Millisecond).String()
}

// Human reports the message sent, or the plan of a dry run.
func (r SessionMessageResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "Would send session %s a message from %s\n\n", r.Session, r.From)
	} else {
		fmt.Fprintf(w, "Sent session %s a message from %s\n\n", r.Session, r.From)
	}
	fmt.Fprintf(w, "- File: %s\n", r.File)
	fmt.Fprintf(w, "- Pointer: %s\n", r.Pointer)
	if r.DryRun {
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written or sent.\n")
		return
	}
	fmt.Fprintf(w, "- Waited: %s\n", r.Waited)
	fmt.Fprintf(w, "\nnext: rhiza session show %s\n", r.Session)
}
