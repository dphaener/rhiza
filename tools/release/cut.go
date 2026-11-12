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
// It pushes once, as the last step. A major release waits for the owner's
// sign-off (publishing its draft GitHub release); publishing makes a
// lightweight tag, which Cut then replaces with the annotated one.
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
	case !AllPass(r.Gates):
		printGates(w, r.Commit, r.Gates)
		fmt.Fprintf(w, "%s is blocked\n", v)
		return Blocked, nil
	}
	local, err := localTag(git, v)
	if err != nil {
		return "", err
	}
	remote, err := remoteTag(git, v)
	if err != nil {
		return "", err
	}
	// Publishing a major's draft release (its sign-off) makes a lightweight
	// tag at the candidate commit. That tag is not a release cut by the tool:
	// only an annotated tag is.
	ours := func(t tagRef) bool { return t.Exists && (t.Annotated || r.Next.Kind != KindMajor) }
	switch {
	case ours(local):
		fmt.Fprintf(w, "%s is already tagged\n", v)
		return Existing, nil
	case ours(remote):
		fmt.Fprintf(w, "%s is already on the remote\n", v)
		return Existing, nil
	}
	replace := false
	if r.Next.Kind == KindMajor {
		signed, err := majorSignOff(w, env, v, r.Commit, r.Notes)
		if err != nil {
			return "", err
		}
		if !signed {
			return SignOff, nil
		}
		for _, t := range []tagRef{local, remote} {
			if t.Exists && t.Commit != r.Commit {
				return "", fmt.Errorf("%s is tagged at %s, not at the candidate %s", v, t.Commit, r.Commit)
			}
		}
		replace = local.Exists || remote.Exists
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
	tag, push := []string{"tag", "-a", v, r.Commit, "-F", f.Name()}, []string{"push", "origin", "refs/tags/" + v}
	if replace {
		// Replace the signed-off release's lightweight tag with the
		// annotated one, at the same commit.
		tag = []string{"tag", "-a", "-f", v, r.Commit, "-F", f.Name()}
		push = []string{"push", "--force", "origin", "refs/tags/" + v}
	}
	if _, err := git(tag...); err != nil {
		return "", err
	}
	if _, err := git(push...); err != nil {
		git("tag", "-d", v) // retry on the next run
		return "", err
	}
	fmt.Fprintf(w, "tagged and pushed %s at %s\n", v, r.Commit)
	return Tagged, nil
}

// tagRef is a tag as one repo holds it: whether it exists, whether it is
// annotated (a tag object) or lightweight, and the commit it points at.
type tagRef struct {
	Exists, Annotated bool
	Commit            string
}

// localTag reads tag v in the local repo.
func localTag(git GitFunc, v string) (tagRef, error) {
	out, err := git("tag", "--list", v)
	if err != nil || out == "" {
		return tagRef{}, err
	}
	kind, err := git("cat-file", "-t", "refs/tags/"+v)
	if err != nil {
		return tagRef{}, err
	}
	commit, err := git("rev-parse", "refs/tags/"+v+"^{commit}")
	if err != nil {
		return tagRef{}, err
	}
	return tagRef{Exists: true, Annotated: kind == "tag", Commit: commit}, nil
}

// remoteTag reads tag v on origin: an annotated tag lists a second, peeled
// line (refs/tags/v^{}) naming its commit.
func remoteTag(git GitFunc, v string) (tagRef, error) {
	out, err := git("ls-remote", "--tags", "origin", "refs/tags/"+v, "refs/tags/"+v+"^{}")
	if err != nil || out == "" {
		return tagRef{}, err
	}
	t := tagRef{Exists: true}
	for _, line := range strings.Split(out, "\n") {
		sha, ref, _ := strings.Cut(line, "\t")
		switch ref {
		case "refs/tags/" + v + "^{}":
			t.Annotated, t.Commit = true, sha
		case "refs/tags/" + v:
			if !t.Annotated {
				t.Commit = sha
			}
		}
	}
	return t, nil
}

// Release states of a version's GitHub release.
const (
	ReleaseNone      = ""          // no release for the version
	ReleaseDraft     = "draft"     // prepared, not yet published by the owner
	ReleasePublished = "published" // the owner published it: the sign-off
)

// majorSignOff reports whether the owner has signed off major release v. The
// sign-off is publishing the version's GitHub release. With no release yet
// it prepares a draft (the candidate commit and the notes) for the owner to
// publish; it never creates one twice, and never tags.
func majorSignOff(w io.Writer, env CheckEnv, v, commit, notes string) (bool, error) {
	if env.Release == nil || env.Draft == nil {
		return false, fmt.Errorf("no GitHub release access to sign off %s", v)
	}
	state, err := env.Release(v)
	if err != nil {
		return false, err
	}
	switch state {
	case ReleasePublished:
		fmt.Fprintf(w, "%s is signed off by the owner\n", v)
		return true, nil
	case ReleaseDraft:
		fmt.Fprintf(w, "%s is a major release: its draft release waits for the owner to publish it\n", v)
		return false, nil
	}
	if err := env.Draft(v, commit, notes); err != nil {
		return false, err
	}
	fmt.Fprintf(w, "%s is a major release: drafted a GitHub release at %s, it waits for the owner to publish it\n", v, commit)
	return false, nil
}

// ghRelease is the state of tag's GitHub release, through gh.
func ghRelease(tag string) (string, error) {
	out, err := exec.Command("gh", "release", "list", "--limit", "100", "--json", "tagName,isDraft").Output()
	if err != nil {
		return "", fmt.Errorf("gh release list: %w", err)
	}
	var rs []struct {
		TagName string `json:"tagName"`
		IsDraft bool   `json:"isDraft"`
	}
	if err := json.Unmarshal(out, &rs); err != nil {
		return "", err
	}
	for _, r := range rs {
		if r.TagName == tag {
			if r.IsDraft {
				return ReleaseDraft, nil
			}
			return ReleasePublished, nil
		}
	}
	return ReleaseNone, nil
}

// ghDraft creates a draft GitHub release for tag at commit with the notes.
func ghDraft(tag, commit, notes string) error {
	out, err := exec.Command("gh", "release", "create", tag, "--draft", "--target", commit,
		"--title", tag, "--notes", notes).CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh release create: %w: %s", err, out)
	}
	return nil
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
