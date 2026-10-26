package main

import (
	"testing"
)

// TestNewCtxHome checks that newCtx passes the home directory from HOME,
// and leaves it empty when there is none.
func TestNewCtxHome(t *testing.T) {
	env := func(home string) func(string) string {
		return func(k string) string {
			if k == "HOME" {
				return home
			}
			return ""
		}
	}

	t.Run("set", func(t *testing.T) {
		dir := t.TempDir()
		ctx, err := newCtx(globals{}, nil, env(dir))
		if err != nil {
			t.Fatal(err)
		}
		if ctx.Home != dir {
			t.Errorf("Home = %q, want %q", ctx.Home, dir)
		}
	})

	t.Run("unset", func(t *testing.T) {
		ctx, err := newCtx(globals{}, nil, env(""))
		if err != nil {
			t.Fatal(err)
		}
		if ctx.Home != "" {
			t.Errorf("Home = %q, want empty", ctx.Home)
		}
	})
}
