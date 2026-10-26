package api

// This file holds run.amend, which amends a running agent step: it
// replaces the step's brief or plan file (a file that already exists in
// the run's worktree) with new content, journals the amendment and its
// sender's actor in the run's journal, and types a one-line pointer to the
// file at the step's session prompt, as session.message does. The engine
// that holds the run's lock never reads an amended entry back to decide
// anything: the entry is an O_APPEND line beside its own.

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

func init() {
	register(Operation{
		Name:     "run.amend",
		Summary:  "Amend a running agent step: replace its brief or plan file, journal the amendment and type a pointer to it at its prompt",
		Write:    true,
		Request:  reflect.TypeFor[RunAmendRequest](),
		Response: reflect.TypeFor[RunAmendResponse](),
		Run:      runRunAmend,
	})
}

const runAmendNext = "rhiza describe run.amend"

// RunAmendRequest names the run, the step, the file and the new content.
type RunAmendRequest struct {
	Run     string `json:"run,omitempty" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters; default the latest run of the current change set, else the latest run"`
	Step    string `json:"step,omitempty" desc:"The agent step to amend; default the run's current step"`
	File    string `json:"file" desc:"The step's brief or plan file, a path relative to the run's worktree; it must already exist"`
	From    string `json:"from" desc:"A file holding the new content, or - for stdin; the content must not be blank"`
	Timeout string `json:"timeout,omitempty" desc:"How long to wait for the step's session to be at its prompt, a Go duration such as 30s or 10m; default 10m"`
}

// RunAmendResponse is the amendment made, or in a dry run the one that
// would be.
type RunAmendResponse struct {
	Run     string `json:"run" desc:"The run's full ID"`
	Step    string `json:"step" desc:"The amended step"`
	Session string `json:"session" desc:"The step's session, which was sent the pointer"`
	File    string `json:"file" desc:"The amended file, relative to the run's worktree"`
	Path    string `json:"path" desc:"The amended file, an absolute path"`
	Actor   string `json:"actor" desc:"The sender's actor handle"`
	Pointer string `json:"pointer" desc:"The one line typed at the session's prompt"`
	Waited  string `json:"waited" desc:"How long it waited for the prompt, a Go duration"`
	DryRun  bool   `json:"dry_run,omitempty" desc:"True if nothing was written, journaled or sent"`
}

// runRunAmend checks the request, the run, the step's session and the
// file before writing anything; then replaces the file, journals the
// amendment and types the pointer once the session is at its prompt. A
// dry run does all but write, journal and send.
func runRunAmend(ctx *Ctx, req any) (any, error) {
	r := req.(RunAmendRequest)
	rel, err := amendRel(r.File)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(r.From) == "" {
		return nil, Usagef("run amend needs --from: a file holding the new content, or - for stdin").WithNext(runAmendNext)
	}
	timeout, err := messageTimeout(r.Timeout)
	if err != nil {
		return nil, err
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	var j run.Journal
	if ref := strings.TrimSpace(r.Run); ref != "" {
		if j, err = readRun(rs, ref); err != nil {
			return nil, err
		}
	} else if j, err = latestRun(ctx, rs); err != nil {
		return nil, err
	}
	st := j.State()
	if st.Ended {
		return nil, Conflictf("run %s is finished: it %s", j.ID, st.Result).WithNext("rhiza run status " + j.ID)
	}
	step := strings.TrimSpace(r.Step)
	if step == "" {
		step = st.Current
	}
	sid := st.Session[step]
	if step == "" || sid == "" {
		return nil, Conflictf("step %s of run %s has no running session", step, j.ID).WithNext("rhiza run status " + j.ID)
	}
	sessions := session.NewStore(session.Dir(common))
	rec, err := readSession(sessions, sid)
	if err != nil {
		return nil, err
	}
	if err := messageable(rec); err != nil {
		return nil, err
	}
	worktree := st.Worktree
	if worktree == "" {
		worktree = rec.Worktree
	}
	path, mode, err := amendTarget(worktree, rel)
	if err != nil {
		return nil, err
	}
	content, err := readAmendment(ctx, r.From)
	if err != nil {
		return nil, err
	}
	if !ctx.DryRun && ctx.Exec == nil {
		return nil, Errorf("cannot amend a step: this caller runs no programs")
	}
	from, err := messageSender(ctx)
	if err != nil {
		return nil, err
	}
	resp := RunAmendResponse{
		Run: j.ID, Step: step, Session: sid, File: rel, Path: path, Actor: from.Handle,
		Pointer: amendPointer(from.Handle, rel), Waited: roundWait(0),
	}
	if ctx.DryRun {
		resp.DryRun = true
		return resp, nil
	}
	if err := replaceFile(path, content, mode); err != nil {
		return nil, err
	}
	e := run.Entry{Type: run.TypeAmended, At: ctx.Now(), Step: step, Session: sid, Actor: from.Handle, File: rel}
	if err := rs.Append(j.ID, e); err != nil {
		return nil, Errorf("journal the amendment of %s in run %s: %v", rel, j.ID, err).WithNext("rhiza run status " + j.ID)
	}
	waited, err := typePointer(ctx, sessions, sid, resp.Pointer, timeout)
	if err != nil {
		return nil, withNote(err, "the amendment is in "+path)
	}
	resp.Waited = roundWait(waited)
	return resp, nil
}

// amendRel is the request's file cleaned and slash-separated: a path
// relative to the worktree that stays inside it and outside .git.
func amendRel(file string) (string, error) {
	f := strings.TrimSpace(file)
	if f == "" {
		return "", Usagef("run amend needs --file: the step's brief or plan file, relative to the run's worktree").WithNext(runAmendNext)
	}
	if filepath.IsAbs(f) || strings.HasPrefix(f, "/") {
		return "", Usagef("file %q is absolute: give a path relative to the run's worktree", f).WithNext(runAmendNext)
	}
	rel := filepath.ToSlash(filepath.Clean(filepath.FromSlash(f)))
	if rel == "." || rel == ".." || strings.HasPrefix(rel, "../") {
		return "", Usagef("file %q is outside the run's worktree", f).WithNext(runAmendNext)
	}
	first, _, _ := strings.Cut(rel, "/")
	if first == ".git" {
		return "", Usagef("file %q is under .git: amend edits a step's brief or plan, not git's files", f).WithNext(runAmendNext)
	}
	return rel, nil
}

// amendTarget is the absolute path of rel in worktree and its mode. The
// file must already exist as a regular file (amend edits, it does not
// create), and its folder, symlinks resolved, must be inside the worktree.
func amendTarget(worktree, rel string) (string, fs.FileMode, error) {
	if worktree == "" {
		return "", 0, Errorf("the run has no worktree to amend a file in")
	}
	path := filepath.Join(worktree, filepath.FromSlash(rel))
	fi, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", 0, NotFoundf("no file %s in the run's worktree %s: amend edits a file, it does not create one", rel, worktree).
			WithNext(runAmendNext)
	}
	if err != nil {
		return "", 0, Errorf("file %s: %v", path, err)
	}
	if !fi.Mode().IsRegular() {
		return "", 0, Usagef("file %s is not a regular file", rel).WithNext(runAmendNext)
	}
	root, err := filepath.EvalSymlinks(worktree)
	if err != nil {
		return "", 0, Errorf("worktree %s: %v", worktree, err)
	}
	dir, err := filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil {
		return "", 0, Errorf("folder of %s: %v", path, err)
	}
	if in, err := filepath.Rel(root, dir); err != nil || in == ".." || strings.HasPrefix(in, ".."+string(filepath.Separator)) {
		return "", 0, Usagef("file %q is outside the run's worktree", rel).WithNext(runAmendNext)
	}
	return path, fi.Mode().Perm(), nil
}

// readAmendment reads the new content from the file from, or stdin for
// "-". Blank content is refused.
func readAmendment(ctx *Ctx, from string) ([]byte, error) {
	from = strings.TrimSpace(from)
	var data []byte
	var err error
	if from == "-" {
		if ctx.Stdin == nil {
			return nil, Usagef("run amend --from -: no stdin").WithNext(runAmendNext)
		}
		if data, err = io.ReadAll(ctx.Stdin); err != nil {
			return nil, Usagef("run amend --from -: cannot read stdin: %v", err).WithNext(runAmendNext)
		}
	} else if data, err = os.ReadFile(from); err != nil {
		return nil, Usagef("run amend: cannot read %s: %v", from, errors.Unwrap(err)).WithNext(runAmendNext)
	}
	if strings.TrimSpace(string(data)) == "" {
		return nil, Usagef("the new content is blank: amend replaces the file, so give it something").WithNext(runAmendNext)
	}
	return data, nil
}

// replaceFile writes content to path atomically: a temporary file in the
// same folder, given mode, renamed over path.
func replaceFile(path string, content []byte, mode fs.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".amend-*")
	if err != nil {
		return Errorf("amend %s: %v", path, err)
	}
	name := tmp.Name()
	_, err = tmp.Write(content)
	if err == nil {
		err = tmp.Chmod(mode)
	}
	if err == nil {
		err = tmp.Sync()
	}
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(name, path)
	}
	if err != nil {
		os.Remove(name)
		return Errorf("amend %s: %v", path, err)
	}
	return nil
}

// amendPointer is the line typed at the step's prompt, naming the amended
// file relative to the worktree. It is one line: every control character
// becomes a space.
func amendPointer(handle, rel string) string {
	return oneLine("rhiza: " + handle + " amended " + rel + ": read it again and carry on from it.")
}

// Human reports the amendment, or the plan of a dry run.
func (r RunAmendResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "Would amend %s of step %s in run %s\n\n", r.File, r.Step, r.Run)
	} else {
		fmt.Fprintf(w, "Amended %s of step %s in run %s\n\n", r.File, r.Step, r.Run)
	}
	fmt.Fprintf(w, "- File: %s\n", r.Path)
	fmt.Fprintf(w, "- Session: %s\n", r.Session)
	fmt.Fprintf(w, "- Actor: %s\n", r.Actor)
	fmt.Fprintf(w, "- Pointer: %s\n", r.Pointer)
	if r.DryRun {
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written, journaled or sent.\n")
		return
	}
	fmt.Fprintf(w, "- Waited: %s\n", r.Waited)
	fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.Run)
}
