package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/release"
)

// channelBaseEnv points the channels at a fixture: <base>/<channel>/manifest.json.
// It exists for tests; a real run reads the published manifests.
const channelBaseEnv = "RHIZA_TEST_CHANNEL_BASE_URL"

// versionEnv overrides the running version, for tests.
const versionEnv = "RHIZA_TEST_VERSION"

// channelVersion is the JSON of `rhiza version --channel`.
type channelVersion struct {
	Channel string `json:"channel"`
	Current string `json:"current"`
	Newest  string `json:"newest"`
	Newer   bool   `json:"newer"`
}

// versionChannel runs `rhiza version --channel <name>`: it reads the channel's
// manifest and says whether its version is newer than this build. It writes
// nothing.
func versionChannel(g globals, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	channel, err := channelArg(args)
	if err != nil {
		return writeError(stdout, stderr, err, g.json)
	}
	r := &release.Resolver{Client: &http.Client{Timeout: 30 * time.Second}}
	if base := getenv(channelBaseEnv); base != "" {
		base = strings.TrimRight(base, "/")
		r.URLs = map[string]string{
			release.ChannelRelease: base + "/" + release.ChannelRelease + "/" + release.ManifestFile,
			release.ChannelLatest:  base + "/" + release.ChannelLatest + "/" + release.ManifestFile,
		}
	}
	newest, err := r.Newest(context.Background(), channel)
	if err != nil {
		return writeError(stdout, stderr, api.Errorf("%v", err).WithNext("rhiza version"), g.json)
	}
	current := versionString()
	if v := getenv(versionEnv); v != "" {
		current = v
	}
	out := channelVersion{Channel: channel, Current: current, Newest: newest, Newer: versionNewer(newest, current)}
	if g.json {
		b, err := canon.Marshal(out)
		if err != nil {
			return writeError(stdout, stderr, err, true)
		}
		stdout.Write(b)
		return 0
	}
	status := "not newer than this build"
	if out.Newer {
		status = "newer than this build"
	}
	fmt.Fprintf(stdout, "%s %s: %s (running %s)\n", channel, newest, status, current)
	return 0
}

// channelArg returns the value of the one --channel flag.
func channelArg(args []string) (string, error) {
	channel, seen := "", false
	for i := 0; i < len(args); i++ {
		a := args[i]
		name, val, hasVal := strings.Cut(strings.TrimPrefix(a, "--"), "=")
		if !strings.HasPrefix(a, "--") || name != "channel" {
			return "", api.Usagef("unexpected argument %q", a).WithNext("rhiza version --channel release|latest")
		}
		if !hasVal {
			if i+1 >= len(args) {
				return "", api.Usagef("flag --channel needs a value").WithNext("rhiza version --channel release|latest")
			}
			i++
			val = args[i]
		}
		channel, seen = val, true
	}
	if !seen {
		return "", api.Usagef("flag --channel needs a value").WithNext("rhiza version --channel release|latest")
	}
	if _, ok := release.ChannelURLs[channel]; !ok {
		return "", api.Usagef("unknown channel %q (want %s or %s)", channel, release.ChannelRelease, release.ChannelLatest).WithNext("rhiza version --channel release|latest")
	}
	return channel, nil
}

// versionNewer reports whether newest is a later version than current. Two
// vMAJOR.MINOR.PATCH versions compare numerically (a pre-release is older than
// its release); otherwise, such as for a local "dev" build, any different
// version counts as newer.
func versionNewer(newest, current string) bool {
	n, nok := parseSemver(newest)
	c, cok := parseSemver(current)
	if !nok || !cok {
		return newest != current
	}
	for i := 0; i < 3; i++ {
		if n.num[i] != c.num[i] {
			return n.num[i] > c.num[i]
		}
	}
	switch {
	case n.pre == c.pre:
		return false
	case n.pre == "":
		return true
	case c.pre == "":
		return false
	}
	return n.pre > c.pre
}

type semver struct {
	num [3]int
	pre string
}

func parseSemver(v string) (semver, bool) {
	var s semver
	v = strings.TrimPrefix(v, "v")
	v, _, _ = strings.Cut(v, "+")
	v, s.pre, _ = strings.Cut(v, "-")
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return s, false
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return s, false
		}
		s.num[i] = n
	}
	return s, true
}
