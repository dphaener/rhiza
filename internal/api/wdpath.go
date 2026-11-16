package api

import (
	"fmt"
	"path/filepath"
)

// absIn is p made absolute against wd, the caller's working directory,
// never the process's: p cleaned when it is absolute, else joined onto wd.
// A relative p with wd unknown (empty) is an error, not a guess.
func absIn(wd, p string) (string, error) {
	if filepath.IsAbs(p) {
		return filepath.Clean(p), nil
	}
	if wd == "" {
		return "", fmt.Errorf("%s is relative and the working directory is unknown", p)
	}
	return filepath.Join(wd, p), nil
}

// abs is p, a path a request names (a file argument, a folder, a
// workflow file), made absolute against the Ctx's working directory with
// absIn.
func (c *Ctx) abs(p string) (string, error) { return absIn(c.Wd, p) }
