// Package userdir places the user's cache, home and temp folders by the
// rules of os.UserCacheDir, os.UserHomeDir and os.TempDir, but reads the
// environment through a getenv func instead of the process's, so a caller
// running in-process for another environment (a test script) places them
// where that environment says. A nil getenv is os.Getenv.
package userdir

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func lookup(getenv func(string) string) func(string) string {
	if getenv == nil {
		return os.Getenv
	}
	return getenv
}

// Home is os.UserHomeDir over getenv: $HOME on Unix, %USERPROFILE% on
// Windows, $home on Plan 9.
func Home(getenv func(string) string) (string, error) {
	getenv = lookup(getenv)
	env, name := "HOME", "$HOME"
	switch runtime.GOOS {
	case "windows":
		env, name = "USERPROFILE", "%userprofile%"
	case "plan9":
		env, name = "home", "$home"
	}
	if v := getenv(env); v != "" {
		return v, nil
	}
	switch runtime.GOOS {
	case "android":
		return "/sdcard", nil
	case "ios":
		return "/", nil
	}
	return "", errors.New(name + " is not defined")
}

// Cache is os.UserCacheDir over getenv: $XDG_CACHE_HOME (absolute) or
// $HOME/.cache on Unix, $HOME/Library/Caches on Darwin, %LocalAppData% on
// Windows, $home/lib/cache on Plan 9.
func Cache(getenv func(string) string) (string, error) {
	getenv = lookup(getenv)
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = getenv("LocalAppData")
		if dir == "" {
			return "", errors.New("%LocalAppData% is not defined")
		}
	case "darwin", "ios":
		dir = getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}
		dir = filepath.Join(dir, "Library", "Caches")
	case "plan9":
		dir = getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}
		dir = filepath.Join(dir, "lib", "cache")
	default:
		dir = getenv("XDG_CACHE_HOME")
		if dir == "" {
			dir = getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_CACHE_HOME nor $HOME are defined")
			}
			dir = filepath.Join(dir, ".cache")
		} else if !filepath.IsAbs(dir) {
			return "", errors.New("path in $XDG_CACHE_HOME is relative")
		}
	}
	return dir, nil
}

// Temp is os.TempDir over getenv: $TMPDIR, else /tmp, on Unix; %TMP%,
// %TEMP% or %USERPROFILE% on Windows; /tmp on Plan 9.
func Temp(getenv func(string) string) string {
	getenv = lookup(getenv)
	switch runtime.GOOS {
	case "windows":
		for _, k := range []string{"TMP", "TEMP", "USERPROFILE"} {
			if v := getenv(k); v != "" {
				return v
			}
		}
		return os.TempDir()
	case "plan9":
		return "/tmp"
	}
	if v := getenv("TMPDIR"); v != "" {
		return v
	}
	if runtime.GOOS == "android" {
		return "/data/local/tmp"
	}
	return "/tmp"
}

// Getenv is a getenv func over environ, a list in os.Environ() form: the
// value of the last entry for a key wins, as os/exec reads a Cmd's Env.
// Entries without "=" are skipped.
func Getenv(environ []string) func(string) string {
	m := make(map[string]string, len(environ))
	for _, kv := range environ {
		if k, v, ok := cutEnv(kv); ok {
			m[k] = v
		}
	}
	return func(key string) string { return m[key] }
}

// cutEnv splits kv at its first "=" after the first byte (Windows keeps
// "=C:=C:\\" style entries).
func cutEnv(kv string) (string, string, bool) {
	for i := 1; i < len(kv); i++ {
		if kv[i] == '=' {
			return kv[:i], kv[i+1:], true
		}
	}
	return "", "", false
}
