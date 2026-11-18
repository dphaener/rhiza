package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
	"github.com/rogpeppe/go-internal/txtar"
)

// byoExamples are the Bring your own agent page's worked examples, one per
// level, under docs/examples/byo.
var byoExamples = []string{"skill.toml", "conductor.toml", "workflow.toml"}

// Every example workflow under docs/examples/workflows and docs/examples/byo
// passes workflow check.
func TestDocsExampleWorkflows(t *testing.T) {
	files, err := filepath.Glob("../../docs/examples/workflows/*.toml")
	if err != nil || len(files) < 4 {
		t.Fatalf("examples: %v %v", files, err)
	}
	for _, name := range byoExamples {
		f := filepath.Join("../../docs/examples/byo", name)
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("byo example: %v", err)
		}
		files = append(files, f)
	}
	for i, f := range files {
		files[i], _ = filepath.Abs(f)
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, out)
	}
	env := []string{"HOME=" + dir, "XDG_CONFIG_HOME=" + dir, "PATH=" + os.Getenv("PATH")}
	for _, f := range files {
		var out, errb bytes.Buffer
		if code := run([]string{"workflow", "check", f}, dir, "", strings.NewReader(""), &out, &errb, env, ownChildren{}); code != 0 || !strings.Contains(out.String(), "is valid") {
			t.Errorf("%s: exit %d\n%s%s", filepath.Base(f), code, out.String(), errb.String())
		}
	}
}

// docs/bring-your-own-agent.md links each of its examples and the primitives
// page, and the README, docs/agents.md and docs/custom-workflows.md link it.
func TestDocsBringYourOwn(t *testing.T) {
	b, err := os.ReadFile("../../docs/bring-your-own-agent.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	for _, name := range byoExamples {
		if !strings.Contains(doc, "](examples/byo/"+name+")") {
			t.Errorf("docs/bring-your-own-agent.md does not link examples/byo/%s", name)
		}
	}
	for _, s := range []string{"](primitives.md)", "rhiza skill check", "rhiza init --agent", "rhiza.project.agent", "{brief}", "{slice}", "{handback}"} {
		if !strings.Contains(doc, s) {
			t.Errorf("docs/bring-your-own-agent.md does not mention %s", s)
		}
	}
	for _, f := range []string{"../../README.md", "../../docs/agents.md", "../../docs/custom-workflows.md"} {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "bring-your-own-agent.md)") {
			t.Errorf("%s does not link docs/bring-your-own-agent.md", strings.TrimPrefix(f, "../../"))
		}
	}
}

// The byo conductor example keeps every step of the built-in ship-change
// after conduct as the built-in has it, so the page's "Everything after
// conduct is the built-in's" holds.
func TestDocsBringYourOwnConductorSteps(t *testing.T) {
	src, ok := workflow.Builtin("ship-change")
	if !ok {
		t.Fatal("no built-in ship-change")
	}
	builtin, err := workflow.Parse("ship-change", src)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile("../../docs/examples/byo/conductor.toml")
	if err != nil {
		t.Fatal(err)
	}
	byo, err := workflow.Parse("myteam-conduct", b)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(byo.Order, builtin.Order) {
		t.Errorf("conductor.toml steps %q, want the built-in's %q", byo.Order, builtin.Order)
	}
	for _, name := range builtin.Order[1:] {
		if !reflect.DeepEqual(byo.Steps[name], builtin.Steps[name]) {
			t.Errorf("conductor.toml step %s = %+v, want the built-in's %+v", name, byo.Steps[name], builtin.Steps[name])
		}
	}
}

// Every script that runs the byo workflow example under the daemon archives
// a copy of it as .rhiza/workflows/myteam-flow.toml (a script cannot read
// repo files); each copy is byte-identical to docs/examples/byo/workflow.toml.
func TestDocsExampleBYOWorkflowCopies(t *testing.T) {
	want, err := os.ReadFile("../../docs/examples/byo/workflow.toml")
	if err != nil {
		t.Fatal(err)
	}
	scripts, err := filepath.Glob("testdata/script/*.txtar")
	if err != nil {
		t.Fatal(err)
	}
	copies := 0
	for _, s := range scripts {
		a, err := txtar.ParseFile(s)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range a.Files {
			if f.Name != "myteam-flow.toml" {
				continue
			}
			copies++
			if !bytes.Equal(f.Data, want) {
				t.Errorf("%s: archived myteam-flow.toml differs from docs/examples/byo/workflow.toml: copy the docs file into the archive", filepath.Base(s))
			}
		}
	}
	if copies == 0 {
		t.Error("no script archives a copy of docs/examples/byo/workflow.toml as myteam-flow.toml")
	}
}
