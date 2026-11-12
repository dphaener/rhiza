package release

import (
	"strconv"
	"strings"
)

// Newer reports whether newest is a later version than current. Two
// vMAJOR.MINOR.PATCH versions compare by semver: the numbers numerically, then
// the pre-release (older than its release) identifier by identifier, numeric
// ones numerically. A git describe suffix -N-g<sha> counts as N commits after
// the version before it. Otherwise, such as for a local "dev" build, any
// different version counts as newer.
func Newer(newest, current string) bool {
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
	if d := comparePre(n.pre, c.pre); d != 0 {
		return d > 0
	}
	return n.commits > c.commits
}

type semver struct {
	num     [3]int
	pre     []string
	commits int // commits after the version, from a git describe suffix
}

func parseSemver(v string) (semver, bool) {
	var s semver
	v = strings.TrimPrefix(v, "v")
	v, _, _ = strings.Cut(v, "+")
	v = strings.TrimSuffix(v, "-dirty")
	v, s.commits = cutDescribe(v)
	v, pre, hasPre := strings.Cut(v, "-")
	if hasPre {
		s.pre = strings.Split(pre, ".")
	}
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

// cutDescribe strips a trailing git describe suffix -N-g<hex> and returns the
// rest and N, or v and 0 when there is none.
func cutDescribe(v string) (string, int) {
	i := strings.LastIndex(v, "-g")
	if i < 0 || i+2 == len(v) || strings.Trim(v[i+2:], "0123456789abcdef") != "" {
		return v, 0
	}
	j := strings.LastIndex(v[:i], "-")
	if j < 0 || !isDigits(v[j+1:i]) {
		return v, 0
	}
	n, err := strconv.Atoi(v[j+1 : i])
	if err != nil {
		return v, 0
	}
	return v[:j], n
}

// comparePre compares two semver pre-releases, split on dots: none is greater
// than any, numeric identifiers compare numerically and sort before
// alphanumeric ones, and on an equal prefix the longer one is greater.
func comparePre(a, b []string) int {
	switch {
	case len(a) == 0 && len(b) == 0:
		return 0
	case len(a) == 0:
		return 1
	case len(b) == 0:
		return -1
	}
	for i := 0; i < len(a) && i < len(b); i++ {
		x, y := a[i], b[i]
		xn, yn := isDigits(x), isDigits(y)
		switch {
		case xn && yn:
			xi, _ := strconv.Atoi(x)
			yi, _ := strconv.Atoi(y)
			if xi != yi {
				if xi > yi {
					return 1
				}
				return -1
			}
		case xn:
			return -1
		case yn:
			return 1
		default:
			if c := strings.Compare(x, y); c != 0 {
				return c
			}
		}
	}
	switch {
	case len(a) > len(b):
		return 1
	case len(a) < len(b):
		return -1
	}
	return 0
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
