package main

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// semverLess orders vX.Y.Z[-pre][+build] by the semver rules (build ignored).
func semverLess(a, b string) bool {
	split := func(v string) ([3]int, []string) {
		v, _, _ = strings.Cut(strings.TrimPrefix(v, "v"), "+")
		core, pre, _ := strings.Cut(v, "-")
		var n [3]int
		for i, p := range strings.SplitN(core, ".", 3) {
			n[i], _ = strconv.Atoi(p)
		}
		if pre == "" {
			return n, nil
		}
		return n, strings.Split(pre, ".")
	}
	an, ap := split(a)
	bn, bp := split(b)
	if an != bn {
		for i := range an {
			if an[i] != bn[i] {
				return an[i] < bn[i]
			}
		}
	}
	if len(ap) == 0 || len(bp) == 0 {
		return len(ap) > 0 && len(bp) == 0
	}
	for i := 0; i < len(ap) && i < len(bp); i++ {
		x, xe := strconv.Atoi(ap[i])
		y, ye := strconv.Atoi(bp[i])
		switch {
		case xe == nil && ye == nil && x != y:
			return x < y
		case xe == nil && ye != nil:
			return true
		case xe != nil && ye == nil:
			return false
		case xe != nil && ye != nil && ap[i] != bp[i]:
			return ap[i] < bp[i]
		}
	}
	return len(ap) < len(bp)
}

func TestLatestVersion(t *testing.T) {
	for _, c := range []struct {
		tag     string
		n       int
		sha, as string
	}{
		{"", 12, "abc1234", "v0.0.1-main.12+abc1234"},
		{"v0.1.0", 3, "abc1234", "v0.1.1-main.3+abc1234"},
		{"v1.2.9", 40, "abc1234", "v1.2.10-main.40+abc1234"},
		{"v0.1.0", 0, "abc1234", "v0.1.0"},
	} {
		got, err := latestVersion(c.tag, c.n, c.sha)
		if err != nil || got != c.as {
			t.Errorf("latestVersion(%q,%d) = %q, %v; want %q", c.tag, c.n, got, err, c.as)
		}
	}
	if _, err := latestVersion("release-1", 1, "x"); err == nil {
		t.Error("bad tag accepted")
	}
}

func TestLatestVersionOrdering(t *testing.T) {
	tag := "v0.1.0"
	prev, _ := latestVersion(tag, 3, "aaaaaaa")
	next, _ := latestVersion(tag, 4, "bbbbbbb")
	ten, _ := latestVersion(tag, 10, "ccccccc")
	if !semverLess(tag, prev) {
		t.Errorf("%s should sort after %s", prev, tag)
	}
	if !semverLess(prev, next) || !semverLess(next, ten) {
		t.Errorf("builds out of order: %s %s %s", prev, next, ten)
	}
	if !semverLess(ten, "v0.1.1") {
		t.Errorf("%s should sort before v0.1.1", ten)
	}
}

func TestGitLatestVersion(t *testing.T) {
	dir := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		c := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "tag.gpgsign=false"}, args...)...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "-q")
	git("commit", "-q", "--allow-empty", "-m", "one")
	git("commit", "-q", "--allow-empty", "-m", "two")
	v, err := gitLatestVersion(dir)
	if err != nil || !strings.HasPrefix(v, "v0.0.1-main.2+") {
		t.Fatalf("no tag: %q, %v", v, err)
	}
	git("tag", "v0.3.0")
	if v, _ := gitLatestVersion(dir); v != "v0.3.0" {
		t.Fatalf("tag at HEAD: %q", v)
	}
	git("commit", "-q", "--allow-empty", "-m", "three")
	git("commit", "-q", "--allow-empty", "-m", "four")
	if v, _ := gitLatestVersion(dir); !strings.HasPrefix(v, "v0.3.1-main.2+") {
		t.Fatalf("tag 2 back: %q", v)
	}
}
