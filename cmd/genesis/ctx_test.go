package main

import (
	"testing"
)

// TestNewCtxHome checks that newCtx passes the home directory, and leaves
// it empty when there is none.
func TestNewCtxHome(t *testing.T) {
	getenv := func(string) string { return "" }

	t.Run("set", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("HOME", dir)
		ctx, err := newCtx(globals{}, nil, getenv)
		if err != nil {
			t.Fatal(err)
		}
		if ctx.Home != dir {
			t.Errorf("Home = %q, want %q", ctx.Home, dir)
		}
	})

	t.Run("unset", func(t *testing.T) {
		t.Setenv("HOME", "")
		ctx, err := newCtx(globals{}, nil, getenv)
		if err != nil {
			t.Fatal(err)
		}
		if ctx.Home != "" {
			t.Errorf("Home = %q, want empty", ctx.Home)
		}
	})
}
