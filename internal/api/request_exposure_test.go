package api

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	_ "github.com/dphaener/rhiza/internal/adapter/claudecode"
)

// exposureRepo is a temporary repo whose .rhiza/workflows holds files, so
// they replace the built-ins of the same name; no personal workflows.
func exposureRepo(t *testing.T, files map[string]string) *Ctx {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	root := t.TempDir()
	dir := filepath.Join(root, filepath.FromSlash(workflowsDir))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(dir, name+".toml"), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return &Ctx{Root: root}
}

func exposureSettings(plan, accept, ship, mode string) []Setting {
	return []Setting{
		{Key: "rhiza.autonomy.plan", Value: plan},
		{Key: "rhiza.autonomy.accept", Value: accept},
		{Key: "rhiza.autonomy.ship", Value: ship},
		{Key: "rhiza.daemon.permissionMode", Value: mode},
	}
}

const (
	exposurePlan = `
name = "plan"
[step.interview]
kind = "agent"
prompt = "plan {change}"
until = "true"
until_when = "done"
`
	exposureShip = `
name = "ship-change"
[step.build]
kind = "workflow"
workflow = "build"
on.pass = "review"
[step.review]
kind = "workflow"
workflow = "review"
`
	exposureBuild = `
name = "build"
[step.build]
kind = "leaves"
[step.build.builder]
prompt = "build {leaf}"
%s
[step.build.verifier]
prompt = "verify {leaf}"
`
	exposureReview = `
name = "review"
[step.review]
kind = "checks"
[[step.review.check]]
name = "qa"
prompt = "check {check}"
[step.review.critic]
prompt = "critic {check}"
%s
[step.review.fixer]
prompt = "fix {cycle}"
[step.review.verifier]
prompt = "verify {leaf}"
`
)

func exposureNames(t *testing.T, got []ExposedWorkflow) map[string]ExposedWorkflow {
	t.Helper()
	m := map[string]ExposedWorkflow{}
	for _, e := range got {
		m[e.Workflow] = e
	}
	return m
}

func TestRequestExposedWorkflows(t *testing.T) {
	const bypass = "bypassPermissions"

	t.Run("plan auto under a bypass mode flags plan only", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		got := requestExposedWorkflows(ctx, exposureSettings("auto", "owner", "owner", bypass))
		if len(got) != 1 || got[0].Workflow != "plan" || got[0].Trigger != TriggerRequest {
			t.Fatalf("got %+v, want only plan from a request", got)
		}
		if len(got[0].Steps) == 0 || got[0].Steps[0] != "interview" {
			t.Errorf("steps %v, want the interview first", got[0].Steps)
		}
	})

	t.Run("plan, accept and ship auto flag plan and ship-change", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		got := exposureNames(t, requestExposedWorkflows(ctx, exposureSettings("auto", "auto", "auto", bypass)))
		if len(got) != 2 || got["plan"].Trigger != TriggerRequest || got["ship-change"].Trigger != TriggerReady {
			t.Fatalf("got %+v, want plan (request) and ship-change (ready)", got)
		}
	})

	t.Run("accept owner flags plan only", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		got := requestExposedWorkflows(ctx, exposureSettings("auto", "owner", "auto", bypass))
		if len(got) != 1 || got[0].Workflow != "plan" {
			t.Fatalf("got %+v, want only plan", got)
		}
	})

	t.Run("plan owner still flags plan, not ship-change", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		got := requestExposedWorkflows(ctx, exposureSettings("owner", "auto", "auto", bypass))
		if len(got) != 1 || got[0].Workflow != "plan" {
			t.Fatalf("got %+v, want only plan: an interactive plan run is no gate", got)
		}
	})

	t.Run("no bypass mode flags nothing", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		if got := requestExposedWorkflows(ctx, exposureSettings("auto", "auto", "auto", "acceptEdits")); len(got) != 0 {
			t.Fatalf("got %+v, want none", got)
		}
		if got := requestExposedWorkflows(ctx, exposureSettings("auto", "auto", "auto", "")); len(got) != 0 {
			t.Fatalf("got %+v, want none with no mode", got)
		}
	})

	t.Run("a step's own mode wins over the run's", func(t *testing.T) {
		src := exposurePlan + "permission_mode = \"acceptEdits\"\n" + "\n[step.arbiter]\nkind = \"arbiter\"\npermission_mode = \"acceptEdits\"\n"
		ctx := exposureRepo(t, map[string]string{"plan": src})
		if got := requestExposedWorkflows(ctx, exposureSettings("auto", "owner", "owner", bypass)); len(got) != 0 {
			t.Fatalf("acceptEdits step under a bypass daemon mode: got %+v, want none", got)
		}
		src = exposurePlan + "permission_mode = \"bypassPermissions\"\n"
		ctx = exposureRepo(t, map[string]string{"plan": src})
		got := requestExposedWorkflows(ctx, exposureSettings("auto", "owner", "owner", ""))
		// The arbiter has no mode of its own and the run has none.
		want := []ExposedWorkflow{{Workflow: "plan", Trigger: TriggerRequest, Steps: []string{"interview"}}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	})

	t.Run("arbiter steps are counted", func(t *testing.T) {
		src := "name = \"plan\"\n" + `
[step.interview]
kind = "agent"
prompt = "plan {change}"
permission_mode = "acceptEdits"
until = "true"
until_when = "done"

[step.arbiter]
kind = "arbiter"
permission_mode = "bypassPermissions"
`
		ctx := exposureRepo(t, map[string]string{"plan": src})
		got := requestExposedWorkflows(ctx, exposureSettings("auto", "owner", "owner", ""))
		want := []ExposedWorkflow{{Workflow: "plan", Trigger: TriggerRequest, Steps: []string{"arbiter"}}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	})

	t.Run("ship-change is flagged through an included builder alone", func(t *testing.T) {
		ctx := exposureRepo(t, map[string]string{
			"plan":        exposurePlan,
			"ship-change": exposureShip,
			"build":       sprintfExposure(exposureBuild, "permission_mode = \"bypassPermissions\""),
			"review":      sprintfExposure(exposureReview, ""),
		})
		got := exposureNames(t, requestExposedWorkflows(ctx, exposureSettings("auto", "auto", "auto", "")))
		if _, ok := got["plan"]; ok {
			t.Errorf("plan has no bypass mode: got %+v", got["plan"])
		}
		want := ExposedWorkflow{Workflow: "ship-change", Trigger: TriggerReady, Steps: []string{"build/build builder"}}
		if !reflect.DeepEqual(got["ship-change"], want) {
			t.Fatalf("got %+v, want %+v", got["ship-change"], want)
		}
	})

	t.Run("ship-change is flagged through an included critic alone", func(t *testing.T) {
		ctx := exposureRepo(t, map[string]string{
			"plan":        exposurePlan,
			"ship-change": exposureShip,
			"build":       sprintfExposure(exposureBuild, ""),
			"review":      sprintfExposure(exposureReview, "permission_mode = \"bypassPermissions\""),
		})
		got := exposureNames(t, requestExposedWorkflows(ctx, exposureSettings("auto", "auto", "auto", "")))
		want := ExposedWorkflow{Workflow: "ship-change", Trigger: TriggerReady, Steps: []string{"review/review critic"}}
		if !reflect.DeepEqual(got["ship-change"], want) {
			t.Fatalf("got %+v, want %+v", got["ship-change"], want)
		}
	})

	t.Run("every bypassing session is named, with the roles of a leaves and checks step", func(t *testing.T) {
		ctx := exposureRepo(t, nil)
		got := exposureNames(t, requestExposedWorkflows(ctx, exposureSettings("auto", "auto", "auto", bypass)))
		steps := got["ship-change"].Steps
		for _, w := range []string{"build/build builder", "build/build verifier", "review/review check qa", "review/review critic", "review/review fixer", "review/review verifier"} {
			found := false
			for _, s := range steps {
				found = found || s == w
			}
			if !found {
				t.Errorf("ship-change steps %v lack %q", steps, w)
			}
		}
	})
}

func sprintfExposure(format, arg string) string {
	return fmt.Sprintf(format, arg)
}
