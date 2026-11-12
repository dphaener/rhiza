package api

// This file holds the pre-push check every push a run makes goes through,
// whatever CI's state: shipPush (ship open-pr and the land push), ship
// rebase's push of a merged head and the watchdog's merge of origin/main
// into a run's branch. It checks exactly the committed head being pushed,
// in a temporary clean checkout of that commit (a detached git worktree,
// removed afterwards), so uncommitted or untracked files in the run's
// worktree neither pass nor fail it: rhiza check (model, log canonicity,
// stale render), then the project's rhiza.ship.prePushCheck command. A
// failure pushes nothing.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
)

// ShipPrePushCheck is a pre-push check that failed: the committed head it
// checked and its failing lines. Nothing was pushed.
type ShipPrePushCheck struct {
	Head  string   `json:"head" desc:"The committed head the check ran on, which was not pushed"`
	Lines []string `json:"lines" desc:"The failing lines: rhiza check's failing findings, or the rhiza.ship.prePushCheck command's error and output tail"`
}

// Human writes the failure and its lines.
func (c ShipPrePushCheck) Human(w io.Writer) {
	fmt.Fprintf(w, "The pre-push check failed on %s, so nothing was pushed:\n", shortSHA(c.Head))
	for _, l := range c.Lines {
		fmt.Fprintln(w, l)
	}
}

// prePushFailure is the error of a pre-push check that failed: its message
// carries the failing lines, and check holds them for a response.
type prePushFailure struct {
	check ShipPrePushCheck
	err   *Error
}

func (f *prePushFailure) Error() string { return f.err.Error() }
func (f *prePushFailure) Unwrap() error { return f.err }

// prePushFailed is the failed check in err, if err is a pre-push check's
// failure.
func prePushFailed(err error) (*ShipPrePushCheck, bool) {
	var f *prePushFailure
	if errors.As(err, &f) {
		c := f.check
		return &c, true
	}
	return nil, false
}

// prePushCheck is shipPrePushCheck; tests whose fake git has no real
// checkout to check replace it.
var prePushCheck = shipPrePushCheck

// shipPrePushCheck checks the commit HEAD of the checkout at dir, which is
// about to be pushed. It returns nil when the check passes, else an error
// naming the head and carrying the failing lines.
func shipPrePushCheck(ctx *Ctx, dir string) error {
	head, err := prePushCheckHead(ctx, dir)
	if err != nil {
		return err
	}
	err = prePushCheckAt(ctx, dir, head)
	journalPrePush(ctx, dir, head, err)
	return err
}

// prePushCheckHead is the commit HEAD of the checkout at dir.
func prePushCheckHead(ctx *Ctx, dir string) (string, error) {
	head, err := gitAt(ctx, dir, "rev-parse", "HEAD")
	if err != nil {
		return "", Errorf("pre-push check: cannot read HEAD in %s: %v", dir, err)
	}
	return strings.TrimSpace(head), nil
}

// prePushCheckAt runs the check on head, the commit of dir's checkout.
func prePushCheckAt(ctx *Ctx, dir, head string) error {
	tmp, err := os.MkdirTemp("", "rhiza-prepush-*")
	if err != nil {
		return Errorf("pre-push check: %v", err)
	}
	defer os.RemoveAll(tmp)
	if _, err := gitAt(ctx, dir, "worktree", "add", "--detach", "-f", tmp, head); err != nil {
		return Errorf("pre-push check: could not check out %s: %v", shortSHA(head), err)
	}
	defer func() {
		gitAt(ctx, dir, "worktree", "remove", "--force", tmp)
		gitAt(ctx, dir, "worktree", "prune")
	}()

	lines, err := prePushCheckLines(ctx, dir, tmp)
	if err != nil {
		return err
	}
	if len(lines) > 0 {
		return &prePushFailure{
			check: ShipPrePushCheck{Head: head, Lines: lines},
			err:   Errorf("the pre-push check failed on %s, so nothing was pushed:\n%s", shortSHA(head), strings.Join(lines, "\n")),
		}
	}
	return nil
}

// journalPrePush journals the pre-push check's result on head in the run
// whose checkout dir is (RHIZA_RUN's, else the unfinished run with that
// worktree): passed or failed, the head, the command and, when failed, the
// end of the output. Best effort: a check that errored before it could
// judge the head, a dry run or a missing run journals nothing.
func journalPrePush(ctx *Ctx, dir, head string, checkErr error) {
	if ctx.DryRun {
		return
	}
	e := run.Entry{Type: run.TypePrePush, At: ctx.Now(), Commit: head, Outcome: run.CheckPassed, Detail: "rhiza check"}
	if out, err := gitAt(ctx, dir, "config", "--get", PrePushCheckKey); err == nil {
		if cmd := strings.TrimSpace(out); cmd != "" {
			e.Detail = "rhiza check; " + cmd
		}
	}
	if checkErr != nil {
		c, ok := prePushFailed(checkErr)
		if !ok {
			return
		}
		e.Outcome = run.CheckFailed
		e.Text = strings.Join(c.Lines, "\n")
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return
	}
	rs := run.NewStore(run.Dir(common))
	id := strings.TrimSpace(ctxGetenv(ctx)(envRun))
	if id == "" {
		js, err := rs.List()
		if err != nil {
			return
		}
		abs, _ := filepath.Abs(dir)
		for _, j := range js {
			st := j.State()
			if w, _ := filepath.Abs(st.Worktree); !st.Ended && st.Worktree != "" && w == abs {
				id, e.Step = j.ID, st.Current
			}
		}
	} else if j, err := rs.Read(id); err == nil {
		id, e.Step = j.ID, j.State().Current
	} else {
		return
	}
	if id != "" {
		rs.Append(id, e)
	}
}

// prePushCheckLines runs rhiza check, then the prePushCheck command, in the
// clean checkout tmp of dir's HEAD, and returns the failing lines: none when
// both pass. The command runs only when rhiza check passed.
func prePushCheckLines(ctx *Ctx, dir, tmp string) ([]string, error) {
	rep, lines, err := prePushRhizaCheck(ctx, tmp)
	if err != nil {
		return nil, err
	}
	if len(lines) > 0 {
		return lines, nil
	}
	if !rep.OK {
		for _, f := range rep.Findings {
			if f.Severity != "fail" {
				continue
			}
			l := "rhiza check: " + f.Rule + ": "
			if f.Path != "" {
				l += f.Path + ": "
			}
			lines = append(lines, l+f.Message)
		}
		if len(lines) == 0 {
			lines = []string{"rhiza check failed"}
		}
		return lines, nil
	}
	cmd := ""
	if out, err := gitAt(ctx, dir, "config", "--get", PrePushCheckKey); err == nil {
		cmd = strings.TrimSpace(out)
	}
	if cmd == "" {
		return nil, nil
	}
	if ctx.Exec == nil {
		return nil, Errorf("pre-push check: %s is set, and no programs can be run here", PrePushCheckKey)
	}
	out, err := ctx.Exec("sh", []string{"-c", cmd}, tmp, nil)
	if err == nil {
		return nil, nil
	}
	return prePushTail([]string{fmt.Sprintf("%s (%s) failed: %v", PrePushCheckKey, cmd, err)}, out), nil
}

// prePushRhizaCheck runs rhiza check in the clean checkout tmp. In Rhiza's
// own repo it runs a build of that head, as ship land renders with one: a
// branch that changes rendering is checked against its own rendering, not
// this binary's. Elsewhere it runs in-process. Lines are set when the head
// build could not give a report (a failed build, say), a failure too.
func prePushRhizaCheck(ctx *Ctx, tmp string) (CheckResponse, []string, error) {
	req := CheckRequest{SkipReference: true, BranchNewerFormat: true}
	if ok, err := isRhizaRepo(tmp); err != nil {
		return CheckResponse{}, nil, err
	} else if ok && ctx.Exec != nil {
		return prePushHeadCheck(ctx, tmp)
	}
	c := *ctx
	c.Root, c.DryRun = tmp, false
	c.Git = func(args ...string) (string, error) { return gitAt(ctx, tmp, args...) }
	r, err := runCheck(&c, req)
	if err != nil {
		return CheckResponse{}, nil, Errorf("pre-push check: rhiza check could not run: %v", err)
	}
	return r.(CheckResponse), nil, nil
}

// prePushHeadCheck builds rhiza from the checkout tmp, outside it, and
// runs that build's check --json there.
func prePushHeadCheck(ctx *Ctx, tmp string) (CheckResponse, []string, error) {
	var rep CheckResponse
	bdir, err := os.MkdirTemp("", "rhiza-prepush-bin-*")
	if err != nil {
		return rep, nil, Errorf("pre-push check: %v", err)
	}
	defer os.RemoveAll(bdir)
	bin := filepath.Join(bdir, "rhiza")
	if out, err := ctx.Exec("go", []string{"build", "-o", bin, "./cmd/rhiza"}, tmp, nil); err != nil {
		return rep, prePushTail([]string{fmt.Sprintf("rhiza check: could not build rhiza from the head: %v", err)}, out), nil
	}
	out, err := ctx.Exec(bin, []string{"check", "--json", "--skip-reference", "--branch-newer-format"}, tmp, nil)
	if jerr := json.Unmarshal(bytes.TrimSpace(out), &rep); jerr != nil || rep.Findings == nil && !rep.OK {
		if err == nil {
			err = jerr
		}
		return rep, prePushTail([]string{fmt.Sprintf("rhiza check with the head's build failed: %s", headFailure(out, err))}, out), nil
	}
	return rep, nil, nil
}

// prePushTail appends the tail of out to lines.
func prePushTail(lines []string, out []byte) []string {
	tail := strings.TrimRight(string(out), "\n")
	if len(tail) > shipGateTail {
		tail = "...\n" + tail[len(tail)-shipGateTail:]
	}
	if tail != "" {
		lines = append(lines, tail)
	}
	return lines
}
