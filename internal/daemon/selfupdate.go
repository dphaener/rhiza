package daemon

import (
	"fmt"
	"strings"
	"time"
)

// SelfUpdateEvery is the least time between two fetches of origin/main.
const SelfUpdateEvery = 5 * time.Minute

// SelfUpdate brings the daemon's own build up to date with main: when
// origin/main has moved past the main checkout it fast-forwards the checkout
// and runs make update there, whose rebuilt binary the daemon's re-exec then
// takes over. Git and Make are injectable; the real ones run in the main
// checkout.
type SelfUpdate struct {
	// Git runs git with args in the main checkout.
	Git func(args ...string) (string, error)
	// Make runs make with args in the main checkout.
	Make func(args ...string) error
	// Fetch fetches origin/main, bounded by a timeout; nil fetches through
	// Git.
	Fetch func() error
	// Log reports what was done or skipped.
	Log func(format string, args ...any)
	// Now is the clock; nil fetches on every Run.
	Now func() time.Time
	// Every is the least time between two fetches.
	Every time.Duration

	// fetched is when origin/main was last fetched, or tried to be.
	fetched time.Time

	// tried is the origin/main commit an update was last attempted for, so a
	// failed one is not retried until main moves again.
	tried string
}

func (u *SelfUpdate) logf(format string, args ...any) {
	if u.Log != nil {
		u.Log(format, args...)
	}
}

func (u *SelfUpdate) rev(ref string) (string, error) {
	out, err := u.Git("rev-parse", "--verify", ref)
	return strings.TrimSpace(out), err
}

// Run does one self-update step and reports whether make update ran. Call
// it at the end of a poll with no start pending. A checkout that is dirty,
// off main or diverged, and a failing fetch, fast-forward or make update,
// are logged and skipped; once origin/main was tried, it is not tried again
// until it moves. A Run within Every of the last fetch does nothing.
func (u *SelfUpdate) Run() bool {
	if u == nil || u.Git == nil || u.Make == nil {
		return false
	}
	if u.Now != nil {
		now := u.Now()
		if !u.fetched.IsZero() && now.Sub(u.fetched) < u.Every {
			return false
		}
		u.fetched = now
	}
	if err := u.fetch(); err != nil {
		u.logf("self-update: fetch: %v", err)
		return false
	}
	remote, err := u.rev("refs/remotes/origin/main")
	if err != nil {
		u.logf("self-update: origin/main: %v", err)
		return false
	}
	head, err := u.rev("HEAD")
	if err != nil {
		u.logf("self-update: HEAD: %v", err)
		return false
	}
	// Not moved past the checkout: equal, or the checkout is ahead.
	if remote == head {
		return false
	}
	if _, err := u.Git("merge-base", "--is-ancestor", remote, head); err == nil {
		return false
	}
	if remote == u.tried {
		return false
	}
	u.tried = remote
	if why := u.blocked(); why != "" {
		u.logf("self-update: skipped until main moves again: %s", why)
		return false
	}
	if _, err := u.Git("merge", "--ff-only", remote); err != nil {
		u.logf("self-update: skipped until main moves again: fast-forward: %v", err)
		return false
	}
	if err := u.Make("update"); err != nil {
		u.logf("self-update: make update failed, skipped until main moves again: %v", err)
		return true
	}
	u.logf("self-update: updated to %s", short(remote))
	return true
}

func (u *SelfUpdate) fetch() error {
	if u.Fetch != nil {
		return u.Fetch()
	}
	_, err := u.Git("fetch", "origin", "main")
	return err
}

// blocked says why the checkout cannot be fast-forwarded, or "".
func (u *SelfUpdate) blocked() string {
	b, err := u.Git("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "the main checkout is not on a branch"
	}
	if b = strings.TrimSpace(b); b != "main" {
		return fmt.Sprintf("the main checkout is on %s, not main", b)
	}
	st, err := u.Git("status", "--porcelain", "--untracked-files=no")
	if err != nil {
		return err.Error()
	}
	if strings.TrimSpace(st) != "" {
		return "the main checkout has uncommitted changes"
	}
	return ""
}

func short(sha string) string {
	if len(sha) > 9 {
		return sha[:9]
	}
	return sha
}
