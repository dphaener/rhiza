package daemon

import (
	"fmt"
	"os"

	"github.com/dphaener/rhiza/internal/rootcfg"
	"path/filepath"
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
	// Dir is the main checkout's folder: where untracked log copies that
	// origin/main also holds are dropped before the fast-forward. Empty
	// drops nothing.
	Dir string
	// Now is the clock; nil fetches on every Run.
	Now func() time.Time
	// Every is the least time between two fetches.
	Every time.Duration

	// NoneLogged is set by the caller once it has logged that no Rhiza
	// checkout is registered, so that is said once, not every poll.
	NoneLogged bool

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
	u.dropIncomingLogCopies(remote)
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

// dropIncomingLogCopies removes the untracked .rhiza/log files in the main
// checkout that remote holds byte for byte, so the fast-forward may write
// them as tracked files. Files that differ stay. Best effort.
func (u *SelfUpdate) dropIncomingLogCopies(remote string) {
	if u.Dir == "" {
		return
	}
	out, err := u.Git("ls-files", "--others", "--exclude-standard", "-z", "--", ".rhiza/log")
	if err != nil {
		return
	}
	for _, f := range strings.Split(out, "\x00") {
		if f == "" {
			continue
		}
		ours, err := u.Git("hash-object", "--", f)
		if err != nil {
			continue
		}
		theirs, err := u.Git("rev-parse", "--verify", "--quiet", remote+":"+f)
		if err != nil || strings.TrimSpace(ours) != strings.TrimSpace(theirs) {
			continue
		}
		_ = os.Remove(filepath.Join(u.Dir, f))
	}
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

// RhizaModule is the Go module path of Rhiza's own source.
const RhizaModule = "github.com/dphaener/rhiza"

// RhizaCheckout finds the registered project that is Rhiza's own source
// checkout: the first project, by name, whose go.mod declares RhizaModule.
// It reports false when none is, and the daemon then updates nothing.
func RhizaCheckout(root string, reg rootcfg.Config) (string, bool) {
	for _, name := range reg.Names() {
		dir := reg.Dir(root, name)
		b, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			f := strings.Fields(line)
			if len(f) == 2 && f[0] == "module" && strings.Trim(f[1], `"`) == RhizaModule {
				return dir, true
			}
		}
	}
	return "", false
}
