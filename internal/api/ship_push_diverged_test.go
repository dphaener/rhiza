package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// shipPushDivergedFixture opens the PR, then has another clone push a
// commit to origin's branch writing file with content, and commits the
// local branch's own new commit on a different file.
func shipPushDivergedFixture(t *testing.T, file, content string) (*Ctx, *openPRGH, string, func(...string) string) {
	t.Helper()
	ctx, g, origin, git := openPRFixture(t, true)
	if _, err := callOp(t, ctx, "ship.open-pr", `{}`); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(t.TempDir(), "other")
	if _, err := g.exec("git", []string{"clone", "-q", "-b", "cs/ship-it", origin, other}, "", nil); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(other, file), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, a := range [][]string{{"add", file}, {"commit", "-q", "-m", "earlier push"}, {"push", "-q", "origin", "cs/ship-it"}} {
		if _, err := g.exec("git", a, other, nil); err != nil {
			t.Fatal(err)
		}
	}
	return ctx, g, origin, git
}

func TestShipPushDivergedMergesOrigin(t *testing.T) {
	ctx, g, origin, git := shipPushDivergedFixture(t, "theirs.txt", "theirs\n")
	if err := os.WriteFile(filepath.Join(ctx.Root, "mine.txt"), []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("add", "mine.txt")
	git("commit", "-q", "-m", "mine")
	if _, err := callOp(t, ctx, "ship.open-pr", `{}`); err != nil {
		t.Fatal(err)
	}
	head := git("rev-parse", "HEAD")
	out, err := g.exec("git", []string{"--git-dir", origin, "rev-parse", "cs/ship-it"}, "", nil)
	if err != nil || strings.TrimSpace(string(out)) != head {
		t.Fatalf("origin at %q, want %s (%v)", out, head, err)
	}
	if log := git("log", "--format=%s"); !strings.Contains(log, "earlier push") || !strings.Contains(log, "mine") {
		t.Errorf("history lacks a side:\n%s", log)
	}
	if git("rev-list", "--merges", "--count", "HEAD") == "0" {
		t.Error("no merge commit: history was rebased")
	}
}

func TestShipPushDivergedConflictStops(t *testing.T) {
	ctx, g, origin, git := shipPushDivergedFixture(t, "clash.txt", "theirs\n")
	if err := os.WriteFile(filepath.Join(ctx.Root, "clash.txt"), []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("add", "clash.txt")
	git("commit", "-q", "-m", "mine")
	before := git("rev-parse", "HEAD")
	_, err := callOp(t, ctx, "ship.open-pr", `{}`)
	if err == nil || !strings.Contains(err.Error(), "earlier push") {
		t.Fatalf("err = %v, want one naming origin's commit", err)
	}
	if _, serr := ctx.Git("rev-parse", "--verify", "--quiet", "MERGE_HEAD"); serr == nil {
		t.Error("merge left in progress")
	}
	if got := git("rev-parse", "HEAD"); got != before {
		t.Errorf("HEAD moved to %s", got)
	}
	out, _ := g.exec("git", []string{"--git-dir", origin, "log", "--format=%s", "cs/ship-it"}, "", nil)
	if strings.Contains(string(out), "mine") {
		t.Error("conflicting branch was pushed")
	}
}
