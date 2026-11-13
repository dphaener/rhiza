package workflow

import (
	"fmt"
	"path"
	"strings"
)

// checkGlob refuses an empty or malformed risky path glob.
func checkGlob(g string) error {
	if strings.TrimSpace(g) == "" {
		return fmt.Errorf("is empty; want a path glob, like \".github/**\"")
	}
	for _, seg := range strings.Split(g, "/") {
		if seg == "**" {
			continue
		}
		if _, err := path.Match(seg, ""); err != nil {
			return fmt.Errorf("%q is not a valid glob: %v", g, err)
		}
	}
	return nil
}

// Risky reports whether file, a slash-separated path from the repo root,
// matches one of the role's RiskyPaths.
func (r Role) Risky(file string) bool {
	for _, g := range r.RiskyPaths {
		if MatchPath(g, file) {
			return true
		}
	}
	return false
}

// MatchPath matches a slash-separated path against a glob: `*` and `?`
// stay within one name, `**` as a whole name matches any number of
// directories (none included), and a glob without a slash matches the
// file's name at any depth.
func MatchPath(glob, file string) bool {
	if !strings.Contains(glob, "/") {
		ok, _ := path.Match(glob, path.Base(file))
		return ok
	}
	return matchSegs(strings.Split(glob, "/"), strings.Split(file, "/"))
}

func matchSegs(g, f []string) bool {
	for len(g) > 0 {
		if g[0] == "**" {
			for len(g) > 0 && g[0] == "**" {
				g = g[1:]
			}
			if len(g) == 0 {
				return true
			}
			for i := range f {
				if matchSegs(g, f[i:]) {
					return true
				}
			}
			return false
		}
		if len(f) == 0 {
			return false
		}
		if ok, _ := path.Match(g[0], f[0]); !ok {
			return false
		}
		g, f = g[1:], f[1:]
	}
	return len(f) == 0
}
