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
	"fmt"
	"os"
	"strings"
)

// prePushCheck is shipPrePushCheck; tests whose fake git has no real
// checkout to check replace it.
var prePushCheck = shipPrePushCheck

// shipPrePushCheck checks the commit HEAD of the checkout at dir, which is
// about to be pushed. It returns nil when the check passes, else an error
// naming the head and carrying the failing lines.
func shipPrePushCheck(ctx *Ctx, dir string) error {
	head, err := gitAt(ctx, dir, "rev-parse", "HEAD")
	if err != nil {
		return Errorf("pre-push check: cannot read HEAD in %s: %v", dir, err)
	}
	head = strings.TrimSpace(head)
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
		return Errorf("the pre-push check failed on %s, so nothing was pushed:\n%s", shortSHA(head), strings.Join(lines, "\n"))
	}
	return nil
}

// prePushCheckLines runs rhiza check, then the prePushCheck command, in the
// clean checkout tmp of dir's HEAD, and returns the failing lines: none when
// both pass. The command runs only when rhiza check passed.
func prePushCheckLines(ctx *Ctx, dir, tmp string) ([]string, error) {
	c := *ctx
	c.Root, c.DryRun = tmp, false
	c.Git = func(args ...string) (string, error) { return gitAt(ctx, tmp, args...) }
	r, err := runCheck(&c, CheckRequest{SkipReference: true, BranchNewerFormat: true})
	if err != nil {
		return nil, Errorf("pre-push check: rhiza check could not run: %v", err)
	}
	if rep := r.(CheckResponse); !rep.OK {
		var lines []string
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
	lines := []string{fmt.Sprintf("%s (%s) failed: %v", PrePushCheckKey, cmd, err)}
	tail := strings.TrimRight(string(out), "\n")
	if len(tail) > shipGateTail {
		tail = "...\n" + tail[len(tail)-shipGateTail:]
	}
	if tail != "" {
		lines = append(lines, tail)
	}
	return lines, nil
}
