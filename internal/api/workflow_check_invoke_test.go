package api

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// invokeWorkflow has a leaves step whose builder starts with a Rhiza skill
// and an agent step whose invoke starts with a skill of the user's own.
const invokeWorkflow = `name = "invoke"

[step.build]
kind = "leaves"
on.pass = "tidy"

[step.build.builder]
prompt = "/rhiza-builder Build {leaf}"

[step.build.verifier]
prompt = "Verify {leaf}"

[step.tidy]
kind = "agent"
prompt = "/my-tidy Tidy {change}"
on.pass = "done"
`

func TestWorkflowCheckUnresolvedInvoke(t *testing.T) {
	t.Parallel()
	setup := func(t *testing.T) (*Ctx, string, string) {
		t.Helper()
		root, home := t.TempDir(), t.TempDir()
		ctx := testCtx(root)
		ctx.Home = home
		ctx.Getenv = func(k string) string {
			if k == "HOME" {
				return home
			}
			return ""
		}
		path := filepath.Join(root, "invoke.toml")
		if err := os.WriteFile(path, []byte(invokeWorkflow), 0o644); err != nil {
			t.Fatal(err)
		}
		return ctx, root, path
	}
	check := func(t *testing.T, ctx *Ctx, path string) WorkflowCheckResponse {
		t.Helper()
		got, err := runWorkflowCheck(ctx, WorkflowCheckRequest{Workflow: path})
		if err != nil {
			t.Fatal(err)
		}
		r := got.(WorkflowCheckResponse)
		if !r.Valid || r.ExitCode() != 0 {
			t.Errorf("valid %v exit %d problems %v: a warning changed the verdict", r.Valid, r.ExitCode(), r.Problems)
		}
		return r
	}
	skill := func(t *testing.T, dir, name string) {
		t.Helper()
		d := filepath.Join(dir, ".claude", "skills", name)
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, "SKILL.md"), []byte("# "+name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("Missing", func(t *testing.T) {
		ctx, _, path := setup(t)
		r := check(t, ctx, path)
		want := []string{
			"step build: the builder's invoke starts with /rhiza-builder, which claude has not installed; run rhiza guide install --agent claude",
			"step tidy: the agent step's invoke starts with /my-tidy, which claude has not installed",
		}
		if !slices.Equal(r.Warnings, want) {
			t.Errorf("warnings\n%q\nwant\n%q", r.Warnings, want)
		}
	})
	t.Run("Installed", func(t *testing.T) {
		ctx, root, path := setup(t)
		skill(t, ctx.Home, "rhiza-builder") // the agent's config folder
		skill(t, root, "my-tidy")           // the repo's
		if r := check(t, ctx, path); len(r.Warnings) != 0 {
			t.Errorf("installed skills warned: %q", r.Warnings)
		}
	})
	t.Run("Rhiza skill", func(t *testing.T) {
		files, err := pluginFiles("test")
		if err != nil {
			t.Fatal(err)
		}
		for name, want := range map[string]bool{"rhiza-builder": true, "rhiza:rhiza-verifier": true, "rhiza": true, "my-tidy": false} {
			if got := isRhizaSkill(files, name); got != want {
				t.Errorf("isRhizaSkill(%q) = %v, want %v", name, got, want)
			}
		}
	})
}
