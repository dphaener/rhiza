package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GitFunc runs git in the repo being released and returns its trimmed
// output. Cut tags and pushes only through it, so a test can fake it.
type GitFunc func(args ...string) (string, error)

// realGit runs git in dir.
func realGit(dir string) GitFunc {
	return func(args ...string) (string, error) {
		c := exec.Command("git", args...)
		c.Dir = dir
		out, err := c.CombinedOutput()
		s := strings.TrimSpace(string(out))
		if err != nil {
			return s, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, s)
		}
		return s, nil
	}
}

// Outcome is what Cut did.
type Outcome string

const (
	Tagged   Outcome = "tagged"   // the tag was made and pushed
	Nothing  Outcome = "none"     // nothing to release
	Blocked  Outcome = "blocked"  // a gate fails
	SignOff  Outcome = "sign-off" // a major release waits for the owner
	Existing Outcome = "exists"   // the version is already tagged
)

// Cut releases the next version when it is a minor or a patch and every gate
// passes: it pushes an annotated tag on the candidate commit whose message
// is the version then the notes. It never tags a major release, a version
// already tagged here or on the remote, or a release with a failing gate.
// It pushes once, as the last step.
func Cut(dir string, w io.Writer, env CheckEnv, git GitFunc) (Outcome, error) {
	if _, err := git("fetch", "--quiet", "--tags", "origin"); err != nil {
		return "", err
	}
	r, err := evaluate(dir, env)
	if err != nil {
		return "", err
	}
	v := r.Next.Version
	switch {
	case r.Next.Kind == KindNone:
		fmt.Fprintf(w, "none: no release since %s\n", orNone(r.Tag))
		return Nothing, nil
	case r.Next.Kind == KindMajor:
		fmt.Fprintf(w, "%s is a major release: it waits for the owner's sign-off\n", v)
		return SignOff, nil
	case !AllPass(r.Gates):
		printGates(w, r.Commit, r.Gates)
		fmt.Fprintf(w, "%s is blocked\n", v)
		return Blocked, nil
	}
	if out, err := git("tag", "--list", v); err != nil {
		return "", err
	} else if out != "" {
		fmt.Fprintf(w, "%s is already tagged\n", v)
		return Existing, nil
	}
	if out, err := git("ls-remote", "--tags", "origin", "refs/tags/"+v); err != nil {
		return "", err
	} else if out != "" {
		fmt.Fprintf(w, "%s is already on the remote\n", v)
		return Existing, nil
	}
	f, err := os.CreateTemp("", "release-notes-*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(v + "\n\n" + r.Notes); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	if _, err := git("tag", "-a", v, r.Commit, "-F", f.Name()); err != nil {
		return "", err
	}
	if _, err := git("push", "origin", "refs/tags/"+v); err != nil {
		git("tag", "-d", v) // retry on the next run
		return "", err
	}
	fmt.Fprintf(w, "tagged and pushed %s at %s\n", v, r.Commit)
	return Tagged, nil
}

// DaemonRun is the root daemon's running commit and when it started.
type DaemonRun struct {
	Commit  string
	Started time.Time
}

// daemonStatus reads `rhiza daemon status --json` for the commit and start
// time of the running root daemon.
func daemonStatus() (DaemonRun, error) {
	out, err := exec.Command("rhiza", "daemon", "status", "--json").Output()
	if err != nil {
		return DaemonRun{}, fmt.Errorf("rhiza daemon status: %w", err)
	}
	return parseDaemonStatus(out)
}

func parseDaemonStatus(b []byte) (DaemonRun, error) {
	var s struct {
		Running bool   `json:"running"`
		Commit  string `json:"commit"`
		Started string `json:"started"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return DaemonRun{}, err
	}
	if !s.Running || s.Commit == "" {
		return DaemonRun{}, fmt.Errorf("the daemon is not running a known commit")
	}
	at, err := time.Parse(time.RFC3339, s.Started)
	if err != nil {
		return DaemonRun{}, fmt.Errorf("daemon started %q: %w", s.Started, err)
	}
	return DaemonRun{Commit: s.Commit, Started: at}, nil
}

// RecordDaemonRun appends {"commit","at"} for the daemon's running commit to
// the history file at path, only when the history lacks that commit, and
// reports whether it did. The earliest line of a commit counts, so a daemon
// restarted on the same commit never resets its soak.
func RecordDaemonRun(path string, status func() (DaemonRun, error)) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("no run history file: pass -runs or set RHIZA_RELEASE_RUNS")
	}
	h, err := ReadRunHistory(path)
	if err != nil {
		return false, err
	}
	d, err := status()
	if err != nil {
		return false, err
	}
	if _, ok := h[d.Commit]; ok {
		return false, nil
	}
	line, err := json.Marshal(struct {
		Commit string `json:"commit"`
		At     string `json:"at"`
	}{d.Commit, d.Started.UTC().Format(time.RFC3339)})
	if err != nil {
		return false, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return false, err
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		f.Close()
		return false, err
	}
	return true, f.Close()
}

// parseRecordRun reads record-run's flags; the history defaults to
// $RHIZA_RELEASE_RUNS.
func parseRecordRun(args []string) (string, error) {
	fs := flag.NewFlagSet("record-run", flag.ContinueOnError)
	runs := fs.String("runs", os.Getenv("RHIZA_RELEASE_RUNS"), "daemon run history file")
	err := fs.Parse(args)
	return *runs, err
}
