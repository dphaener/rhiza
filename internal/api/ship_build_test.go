package api

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// onDemandFixture is a ship land fixture whose workflow names a head build
// (SHIP_HEAD_RHIZA) and its command (SHIP_HEAD_BUILD, when build is set).
// The fake command writes the binary and logs "build"; the binary's render
// logs "render". The binary is missing at first.
func onDemandFixture(t *testing.T, build string, buildErr error) (f *sessionFixture, bin string, log *[]string) {
	t.Helper()
	f, _, _ = shipLandFixture(t)
	bin = filepath.Join(t.TempDir(), "rhiza")
	log = new([]string)
	f.ctx.Getenv = func(k string) string {
		switch k {
		case envShipRhiza:
			return bin
		case envShipBuild:
			return build
		}
		return ""
	}
	f.ctx.Exec = func(name string, args []string, dir string, _ []string) ([]byte, error) {
		if dir != f.ctx.Root {
			t.Errorf("%s ran in %s, not %s", name, dir, f.ctx.Root)
		}
		switch {
		case name == "sh" && len(args) == 2 && args[0] == "-c" && args[1] == build:
			*log = append(*log, "build")
			if buildErr != nil {
				return []byte("boom output"), buildErr
			}
			return nil, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755)
		case name == bin && len(args) > 0 && args[0] == "render":
			*log = append(*log, "render")
			return []byte(`{"written":[".rhiza/rendered/a.md"],"removed":[],"unchanged":1}`), nil
		}
		return nil, errors.New(name + " " + strings.Join(args, " ") + ": unexpected call")
	}
	return f, bin, log
}

func TestShipHeadBuildOnDemand(t *testing.T) {
	const cmd = "make head"
	render := func(t *testing.T, f *sessionFixture) error {
		t.Helper()
		tg, err := findShipTarget(f.ctx)
		if err != nil {
			t.Fatal(err)
		}
		_, err = shipRender(f.ctx, tg)
		return err
	}
	want := func(t *testing.T, log *[]string, w ...string) {
		t.Helper()
		if got := strings.Join(*log, ","); got != strings.Join(w, ",") {
			t.Errorf("calls = %s, want %s", got, strings.Join(w, ","))
		}
	}

	t.Run("missing binary is built once, then rendered with", func(t *testing.T) {
		f, bin, log := onDemandFixture(t, cmd, nil)
		if err := render(t, f); err != nil {
			t.Fatal(err)
		}
		want(t, log, "build", "render")
		b, err := os.ReadFile(bin + shipBuiltAtSuffix)
		if err != nil || strings.TrimSpace(string(b)) != "beef000111" {
			t.Errorf("built-at = %q, %v", b, err)
		}
		// Built at the current HEAD: not built again.
		if err := render(t, f); err != nil {
			t.Fatal(err)
		}
		want(t, log, "build", "render", "render")
	})

	t.Run("built at an earlier commit is built again", func(t *testing.T) {
		f, bin, log := onDemandFixture(t, cmd, nil)
		if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(bin+shipBuiltAtSuffix, []byte("0ld0000\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := render(t, f); err != nil {
			t.Fatal(err)
		}
		want(t, log, "build", "render")
		b, _ := os.ReadFile(bin + shipBuiltAtSuffix)
		if strings.TrimSpace(string(b)) != "beef000111" {
			t.Errorf("built-at = %q", b)
		}
	})

	t.Run("built at the current HEAD is not built again", func(t *testing.T) {
		f, bin, log := onDemandFixture(t, cmd, nil)
		if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(bin+shipBuiltAtSuffix, []byte("beef000111\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := render(t, f); err != nil {
			t.Fatal(err)
		}
		want(t, log, "render")
	})

	t.Run("no build command: a missing binary fails", func(t *testing.T) {
		f, _, log := onDemandFixture(t, "", nil)
		err := render(t, f)
		if err == nil || !strings.Contains(err.Error(), "not a built rhiza") {
			t.Fatalf("err = %v", err)
		}
		want(t, log)
	})

	t.Run("a failing build renders nothing", func(t *testing.T) {
		f, _, log := onDemandFixture(t, cmd, errors.New("exit 2"))
		err := render(t, f)
		if err == nil || !strings.Contains(err.Error(), "boom output") || !strings.Contains(err.Error(), "nothing was rendered or pushed") {
			t.Fatalf("err = %v", err)
		}
		want(t, log, "build")
	})

	t.Run("ship land builds on demand, lands, and fails clean on a bad build", func(t *testing.T) {
		f, _, log := onDemandFixture(t, cmd, errors.New("exit 2"))
		_, err := callOp(t, f.ctx, "ship.land", `{}`)
		if err == nil {
			t.Fatal("land with a failing build succeeded")
		}
		want(t, log, "build")
		for _, c := range f.git.calls {
			if strings.HasPrefix(c, "commit") || strings.HasPrefix(c, "push") {
				t.Errorf("git %q after a failed build", c)
			}
		}
	})
}

// headFixture is a ship land fixture whose environment names the head
// build bin (and build, when set) and whose programs are logged: the build
// command writes bin, and bin records its argv and answers land and render.
func headFixture(t *testing.T, build string, writeBin bool) (f *sessionFixture, bin string, calls *[]string) {
	t.Helper()
	f, _, _ = shipLandFixture(t)
	bin = filepath.Join(t.TempDir(), "rhiza-head")
	if writeBin {
		if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	calls = new([]string)
	f.ctx.Getenv = func(k string) string {
		switch k {
		case envShipRhiza:
			return bin
		case envShipBuild:
			return build
		}
		return ""
	}
	f.ctx.Self = filepath.Join(t.TempDir(), "rhiza-installed")
	f.ctx.Exec = func(name string, args []string, dir string, _ []string) ([]byte, error) {
		if dir != f.ctx.Root {
			t.Errorf("%s ran in %s, not %s", name, dir, f.ctx.Root)
		}
		switch {
		case name == "sh" && len(args) == 2 && args[1] == build:
			*calls = append(*calls, "build")
			return nil, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755)
		case name == bin && len(args) > 0 && args[0] == "render":
			*calls = append(*calls, name+" "+strings.Join(args, " "))
			return []byte(`{"written":[".rhiza/rendered/a.md"],"removed":[],"unchanged":1}`), nil
		case name == bin && len(args) > 0 && args[0] == "land":
			*calls = append(*calls, name+" "+strings.Join(args, " "))
			return []byte(`{"change_set":"` + currentOf(t, f.ctx).String() + `","from":"accepted","to":"landed"}`), nil
		}
		return nil, errors.New(name + " " + strings.Join(args, " ") + ": unexpected call")
	}
	return f, bin, calls
}

func TestLandRenderHeadBuild(t *testing.T) {
	t.Run("named binary: render and land exec it and write nothing in-process", func(t *testing.T) {
		f, bin, calls := headFixture(t, "", true)
		if _, err := callOp(t, f.ctx, "render", `{}`); err != nil {
			t.Fatal(err)
		}
		if _, err := callOp(t, f.ctx, "land", `{}`); err != nil {
			t.Fatal(err)
		}
		got := strings.Join(*calls, "|")
		if want := bin + " render --json|" + bin + " land --json"; got != want {
			t.Errorf("calls = %s, want %s", got, want)
		}
		if _, err := os.Stat(filepath.Join(f.ctx.Root, renderedDir)); err == nil {
			t.Error("render wrote in-process")
		}
		// The head binary did the landing: nothing was landed in-process.
		if cs := currentOf(t, f.ctx); cs.IsZero() {
			t.Error("the change set was landed in-process")
		}
	})

	t.Run("unset: in-process", func(t *testing.T) {
		f, _, calls := headFixture(t, "", true)
		f.ctx.Getenv = func(string) string { return "" }
		if _, err := callOp(t, f.ctx, "render", `{}`); err != nil {
			t.Fatal(err)
		}
		if _, err := callOp(t, f.ctx, "land", `{}`); err != nil {
			t.Fatal(err)
		}
		if len(*calls) != 0 {
			t.Errorf("calls = %v", *calls)
		}
	})

	t.Run("the running binary is the build: in-process", func(t *testing.T) {
		f, bin, calls := headFixture(t, "", true)
		f.ctx.Self = bin
		if _, err := callOp(t, f.ctx, "render", `{}`); err != nil {
			t.Fatal(err)
		}
		if _, err := callOp(t, f.ctx, "land", `{}`); err != nil {
			t.Fatal(err)
		}
		if len(*calls) != 0 {
			t.Errorf("calls = %v", *calls)
		}
	})

	t.Run("missing with a build command: built first", func(t *testing.T) {
		f, bin, calls := headFixture(t, "make head", false)
		if _, err := callOp(t, f.ctx, "render", `{}`); err != nil {
			t.Fatal(err)
		}
		if got := strings.Join(*calls, "|"); got != "build|"+bin+" render --json" {
			t.Errorf("calls = %s", got)
		}
	})

	t.Run("missing with no build command: both fail, nothing landed or rendered", func(t *testing.T) {
		f, _, calls := headFixture(t, "", false)
		for _, op := range []string{"render", "land"} {
			_, err := callOp(t, f.ctx, op, `{}`)
			if err == nil || !strings.Contains(err.Error(), "not a built rhiza") {
				t.Errorf("%s err = %v", op, err)
			}
		}
		if len(*calls) != 0 {
			t.Errorf("calls = %v", *calls)
		}
		if _, err := os.Stat(filepath.Join(f.ctx.Root, renderedDir)); err == nil {
			t.Error("render wrote in-process")
		}
	})
}
