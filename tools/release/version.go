package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var tagRE = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`)

// latestVersion is the latest-channel version of a build: the next patch of
// the last v* tag as a semver prerelease, vX.Y.(Z+1)-main.<commits since the
// tag>+<short sha>. With no tag the base is v0.0.0, so the first build is
// v0.0.1-main.N+sha. It sorts after the tag, before the next patch, and after
// the previous latest build (N only grows). A build at the tag itself is the
// tag.
func latestVersion(tag string, commits int, sha string) (string, error) {
	if tag == "" {
		tag = "v0.0.0"
	}
	m := tagRE.FindStringSubmatch(tag)
	if m == nil {
		return "", fmt.Errorf("tag %q is not vX.Y.Z", tag)
	}
	if commits == 0 && tag != "v0.0.0" {
		return tag, nil
	}
	patch, _ := strconv.Atoi(m[3])
	return fmt.Sprintf("v%s.%s.%d-main.%d+%s", m[1], m[2], patch+1, commits, sha), nil
}

// gitLatestVersion computes latestVersion for HEAD in dir.
func gitLatestVersion(dir string) (string, error) {
	git := func(args ...string) (string, error) {
		c := exec.Command("git", args...)
		c.Dir = dir
		b, err := c.Output()
		return strings.TrimSpace(string(b)), err
	}
	sha, err := git("rev-parse", "--short=7", "HEAD")
	if err != nil {
		return "", err
	}
	tag, err := git("describe", "--tags", "--match", "v[0-9]*.[0-9]*.[0-9]*", "--abbrev=0")
	if err != nil {
		tag = ""
	}
	rng := "HEAD"
	if tag != "" {
		rng = tag + "..HEAD"
	}
	n, err := git("rev-list", "--count", rng)
	if err != nil {
		return "", err
	}
	commits, err := strconv.Atoi(n)
	if err != nil {
		return "", err
	}
	return latestVersion(tag, commits, sha)
}
