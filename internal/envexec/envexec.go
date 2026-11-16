// Package envexec starts programs for a caller that runs for an
// environment and working directory other than the process's (run() in a
// test script): a bare program name is found on the PATH of the given
// environment, never the process's, a program path is joined onto the
// given working directory, never the process's, and a program found
// through a '.', empty or relative PATH entry is refused with exec.ErrDot,
// as os/exec refuses it.
// os/exec's Command and LookPath read the process's PATH and resolve
// against the process's working directory, so in-process callers would
// start the host's programs instead of a script's stand-ins.
package envexec

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dphaener/rhiza/internal/userdir"
)

// LookPath finds file as exec.LookPath does, but on the PATH of environ (a
// list in os.Environ() form, the last PATH entry winning) and against wd:
// a file holding a path separator is joined onto wd when relative and
// returned without a PATH search. A match through a relative (or empty,
// meaning the current folder) PATH entry is refused, as exec.LookPath
// refuses it: an *exec.Error wrapping exec.ErrDot, with the match (joined
// onto wd) as the path. The result is absolute when wd is. A miss is an
// *exec.Error wrapping exec.ErrNotFound.
func LookPath(file string, environ []string, wd string) (string, error) {
	getenv := userdir.Getenv(environ)
	if hasSeparator(file) {
		p := file
		if !filepath.IsAbs(p) {
			p = filepath.Join(wd, p)
		}
		found, err := findExecutable(p, getenv)
		if err != nil {
			return "", &exec.Error{Name: file, Err: err}
		}
		return found, nil
	}
	for _, dir := range filepath.SplitList(getenv("PATH")) {
		if dir == "" {
			dir = "."
		}
		rel := !filepath.IsAbs(dir)
		if rel {
			dir = filepath.Join(wd, dir)
		}
		if found, err := findExecutable(filepath.Join(dir, file), getenv); err == nil {
			if rel {
				return found, &exec.Error{Name: file, Err: exec.ErrDot}
			}
			return found, nil
		}
	}
	return "", &exec.Error{Name: file, Err: exec.ErrNotFound}
}

// findExecutable returns p when it is an executable regular file (on
// Windows, p or p with one of PATHEXT's extensions when p has none).
func findExecutable(p string, getenv func(string) string) (string, error) {
	if runtime.GOOS != "windows" {
		return p, executable(p)
	}
	if filepath.Ext(p) != "" {
		if executable(p) == nil {
			return p, nil
		}
	}
	exts := getenv("PATHEXT")
	if exts == "" {
		exts = ".com;.exe;.bat;.cmd"
	}
	for _, e := range filepath.SplitList(strings.ToLower(exts)) {
		if e != "" && executable(p+e) == nil {
			return p + e, nil
		}
	}
	return "", exec.ErrNotFound
}

func executable(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		return err
	}
	m := fi.Mode()
	if m.IsDir() {
		return fs.ErrPermission
	}
	if runtime.GOOS != "windows" && m&0o111 == 0 {
		return fs.ErrPermission
	}
	return nil
}

// Command is exec.CommandContext (exec.Command when c is nil) with name
// found by LookPath on environ's PATH against wd, and Env set to environ.
// Args[0] stays name, as a shell passes it. When name is not found, or is
// found only through a relative PATH entry (exec.ErrDot), the command's
// Err holds the error, so Start, Run and Output return it.
func Command(c context.Context, environ []string, wd, name string, args ...string) *exec.Cmd {
	var path string
	var err error
	switch {
	case hasSeparator(name):
		// As os/exec does, no search: Start reports a missing program.
		path = name
		if !filepath.IsAbs(path) {
			path = filepath.Join(wd, path)
		}
	default:
		if path, err = LookPath(name, environ, wd); err != nil {
			// Never let os/exec search the process's PATH for it: Err
			// below stops Start before this path is used.
			path = filepath.Join(string(filepath.Separator), "nonexistent", name)
		}
	}
	var cmd *exec.Cmd
	if c != nil {
		cmd = exec.CommandContext(c, path, args...)
	} else {
		cmd = exec.Command(path, args...)
	}
	cmd.Args[0] = name
	cmd.Env = environ
	if err != nil {
		cmd.Err = err
	}
	return cmd
}

func hasSeparator(file string) bool {
	return strings.ContainsRune(file, '/') || strings.ContainsRune(file, filepath.Separator)
}
