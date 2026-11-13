// Package ciworkflow holds the test that pins the shape of Rhiza's own CI
// workflow file.
package ciworkflow

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type line struct {
	indent int
	text   string
}

// block returns the lines nested under the first key line equal to key at
// the given indent within lines, or nil. Only the small YAML subset the
// workflow uses is understood: block mappings, one key per line.
func block(lines []line, indent int, key string) []line {
	for i, l := range lines {
		if l.indent != indent || !strings.HasPrefix(l.text, key+":") {
			continue
		}
		var out []line
		for _, n := range lines[i+1:] {
			if n.indent <= indent {
				break
			}
			out = append(out, n)
		}
		return out
	}
	return nil
}

// value returns the inline value of key at indent, or "" and false.
func value(lines []line, indent int, key string) (string, bool) {
	for _, l := range lines {
		if l.indent == indent && strings.HasPrefix(l.text, key+":") {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(l.text, key+":")), `"'`), true
		}
	}
	return "", false
}

func parse(t *testing.T) []line {
	t.Helper()
	data, err := os.ReadFile("../../.github/workflows/ci.yml")
	if err != nil {
		t.Fatal(err)
	}
	var out []line
	for _, raw := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		out = append(out, line{len(raw) - len(strings.TrimLeft(raw, " ")), trimmed})
	}
	return out
}

func TestCIRunsOnPullRequestsAndDispatch(t *testing.T) {
	lines := parse(t)
	on := block(lines, 0, "on")
	if on == nil {
		t.Fatal("ci.yml has no on: block")
	}
	var triggers []string
	for _, l := range on {
		if l.indent == on[0].indent {
			triggers = append(triggers, strings.TrimSuffix(strings.SplitN(l.text, ":", 2)[0], ":"))
		}
	}
	if len(triggers) != 2 || triggers[0] != "pull_request" || triggers[1] != "workflow_dispatch" {
		t.Errorf("ci.yml triggers = %v; want exactly pull_request and workflow_dispatch", triggers)
	}
	wd := block(on, on[0].indent, "workflow_dispatch")
	sha := block(block(wd, wd[0].indent, "inputs"), wd[0].indent+2, "sha")
	if sha == nil {
		t.Fatal("workflow_dispatch declares no input sha")
	}
	if v, _ := value(sha, sha[0].indent, "required"); v != "true" {
		t.Errorf("input sha required = %q; want true", v)
	}
}

func TestCICheckoutAndRunName(t *testing.T) {
	lines := parse(t)
	found := false
	for _, l := range lines {
		if l.text == "ref: ${{ inputs.sha || github.sha }}" {
			found = true
		}
	}
	if !found {
		t.Error("checkout ref is not ${{ inputs.sha || github.sha }}")
	}
	name, _ := value(lines, 0, "run-name")
	if !strings.Contains(name, "inputs.sha") {
		t.Errorf("run-name = %q; want it to include inputs.sha", name)
	}
}

func TestCICancelsStaleRuns(t *testing.T) {
	c := block(parse(t), 0, "concurrency")
	if c == nil {
		t.Fatal("ci.yml has no concurrency: block")
	}
	group, _ := value(c, c[0].indent, "group")
	if !strings.Contains(group, "github.event.pull_request.number") {
		t.Errorf("concurrency group = %q; want it keyed on github.event.pull_request.number", group)
	}
	if !strings.Contains(group, "|| inputs.sha") {
		t.Errorf("concurrency group = %q; want it falling back to inputs.sha", group)
	}
	if v, _ := value(c, c[0].indent, "cancel-in-progress"); v != "true" {
		t.Errorf("cancel-in-progress = %q; want true", v)
	}
}

// jobBlock returns the lines under the job named name.
func jobBlock(t *testing.T, name string) []line {
	t.Helper()
	jobs := block(parse(t), 0, "jobs")
	j := block(jobs, 2, name)
	if j == nil {
		t.Fatalf("ci.yml has no job %q", name)
	}
	return j
}

func has(lines []line, text string) bool {
	for _, l := range lines {
		if l.text == text {
			return true
		}
	}
	return false
}

func TestCIJobsAndMatrices(t *testing.T) {
	jobs := block(parse(t), 0, "jobs")
	var names []string
	for _, l := range jobs {
		if l.indent == 2 {
			names = append(names, strings.TrimSuffix(l.text, ":"))
		}
	}
	want := []string{"changes", "model", "lint", "scripts", "api", "rest", "budgets", "check"}
	if strings.Join(names, ",") != strings.Join(want, ",") {
		t.Fatalf("ci.yml jobs = %v; want %v", names, want)
	}
	targets := map[string]string{
		"lint": "make lint", "scripts": "make test-scripts SHARD=${{ matrix.shard }}",
		"api": "make test-api SHARD=${{ matrix.shard }}", "rest": "make test-rest",
		"budgets": "make test-budgets",
	}
	shards := map[string]string{
		"scripts": "[1/4, 2/4, 3/4, 4/4]", "api": "[1/2, 2/2]",
	}
	for name, run := range targets {
		j := jobBlock(t, name)
		if v, _ := value(j, 4, "run"); !strings.Contains(strings.Join(texts(j), "\n"), "run: "+run) {
			t.Errorf("job %s: no step running %q (%q)", name, run, v)
		}
		if !has(j, "cache: true") || !has(j, "cache-dependency-path: go.sum") {
			t.Errorf("job %s: setup-go must set cache: true and cache-dependency-path: go.sum", name)
		}
		if shard, ok := shards[name]; ok {
			s := block(j, 4, "strategy")
			if !has(s, "fail-fast: false") {
				t.Errorf("job %s: strategy.fail-fast must be false", name)
			}
			if !has(s, "shard: "+shard) {
				t.Errorf("job %s: matrix shard want %s", name, shard)
			}
		} else if block(j, 4, "strategy") != nil {
			t.Errorf("job %s: unexpected matrix", name)
		}
	}
}

func texts(ls []line) []string {
	var out []string
	for _, l := range ls {
		out = append(out, l.text)
	}
	return out
}

func TestCIChangesDecidesFromGitDiff(t *testing.T) {
	j := jobBlock(t, "changes")
	if out, _ := value(block(j, 4, "outputs"), 6, "code"); !strings.Contains(out, "steps.diff.outputs.code") {
		t.Errorf("changes output code = %q; want steps.diff.outputs.code", out)
	}
	body := strings.Join(texts(j), "\n")
	for _, want := range []string{".github/scripts/ci-changes.sh", "steps.green.outputs.sha", "code=true", `!= pull_request`} {
		if !strings.Contains(body, want) {
			t.Errorf("changes job lacks %q", want)
		}
	}
}

func TestCIGoJobsRunOnlyOnCodeChanges(t *testing.T) {
	for _, name := range []string{"lint", "scripts", "api", "rest", "budgets"} {
		j := jobBlock(t, name)
		if v, _ := value(j, 4, "needs"); v != "changes" {
			t.Errorf("job %s needs = %q; want changes", name, v)
		}
		if v, _ := value(j, 4, "if"); v != "needs.changes.outputs.code == 'true" {
			t.Errorf("job %s if = %q; want code == 'true'", name, v)
		}
	}
}

func TestCIModelJobRunsOnlyOnModelChanges(t *testing.T) {
	j := jobBlock(t, "model")
	if v, _ := value(j, 4, "needs"); v != "changes" {
		t.Errorf("model needs = %q; want changes", v)
	}
	if v, _ := value(j, 4, "if"); v != "needs.changes.outputs.code == 'false" {
		t.Errorf("model if = %q; want code == 'false'", v)
	}
	body := strings.Join(texts(j), "\n")
	for _, run := range []string{"make self-check", "go test -trimpath ./internal/op -run TestEveryFieldInRepoLogIsKnown"} {
		if !strings.Contains(body, "run: "+run) {
			t.Errorf("model job: no step running %q", run)
		}
	}
	if !has(j, "cache: true") {
		t.Error("model job: setup-go must set cache: true")
	}
}

func TestCICheckAggregatesEveryJob(t *testing.T) {
	j := jobBlock(t, "check")
	needs, _ := value(j, 4, "needs")
	got := map[string]bool{}
	for _, n := range strings.Split(strings.Trim(needs, "[]"), ",") {
		got[strings.TrimSpace(n)] = true
	}
	for _, n := range []string{"changes", "model", "lint", "scripts", "api", "rest", "budgets"} {
		if !got[n] {
			t.Errorf("check does not need %s (needs: %q)", n, needs)
		}
	}
	if len(got) != 7 {
		t.Errorf("check needs %q; want exactly the seven other jobs", needs)
	}
	if v, _ := value(j, 4, "if"); v != "always()" {
		t.Errorf("check if = %q; want always()", v)
	}
	body := strings.Join(texts(j), "\n")
	for _, want := range []string{"needs.*.result", "needs.changes.result", `!= success`, `!= skipped`, "exit 1", "no job ran"} {
		if !strings.Contains(body, want) {
			t.Errorf("check lacks %q", want)
		}
	}
}

// gitRepo is a scratch repository with helpers.
type gitRepo struct {
	t   *testing.T
	dir string
}

func (r gitRepo) git(args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// commit writes the files and commits them, returning the new sha.
func (r gitRepo) commit(files map[string]string) string {
	r.t.Helper()
	for name, body := range files {
		p := filepath.Join(r.dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			r.t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			r.t.Fatal(err)
		}
	}
	r.git("add", "-A")
	r.git("commit", "-m", "c")
	return r.git("rev-parse", "HEAD")
}

// decide runs the changes script at HEAD and returns its output.
func (r gitRepo) decide(base, green string) string {
	r.t.Helper()
	script, err := filepath.Abs("../../.github/scripts/ci-changes.sh")
	if err != nil {
		r.t.Fatal(err)
	}
	cmd := exec.Command("bash", script, base, green)
	cmd.Dir = r.dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("ci-changes.sh: %v\n%s", err, out)
	}
	return strings.TrimSpace(string(out))
}

func newRepo(t *testing.T) (gitRepo, string) {
	t.Helper()
	r := gitRepo{t, t.TempDir()}
	r.git("init", "-b", "main")
	base := r.commit(map[string]string{"main.go": "package main\n", ".rhiza/a": "1"})
	return r, base
}

func TestChangesScript(t *testing.T) {
	t.Run("land-like commit changing only .rhiza/ after green", func(t *testing.T) {
		r, base := newRepo(t)
		r.git("checkout", "-b", "pr")
		green := r.commit(map[string]string{"main.go": "package main // pr\n"})
		r.commit(map[string]string{".rhiza/log/x.json": "{}", ".rhiza/rendered/y.md": "y"})
		if got := r.decide(base, green); got != "code=false" {
			t.Errorf("got %s; want code=false", got)
		}
	})
	t.Run("merge of main bringing only .rhiza/", func(t *testing.T) {
		r, base := newRepo(t)
		r.git("checkout", "-b", "pr")
		green := r.commit(map[string]string{"main.go": "package main // pr\n"})
		r.git("checkout", "main")
		r.commit(map[string]string{".rhiza/log/m.json": "{}"})
		r.git("checkout", "pr")
		r.git("merge", "--no-ff", "-m", "merge main", "main")
		if got := r.decide(base, green); got != "code=false" {
			t.Errorf("got %s; want code=false", got)
		}
	})
	t.Run("code change since green", func(t *testing.T) {
		r, base := newRepo(t)
		r.git("checkout", "-b", "pr")
		green := r.commit(map[string]string{"main.go": "package main // pr\n"})
		r.commit(map[string]string{"main.go": "package main // more\n", ".rhiza/b": "2"})
		if got := r.decide(base, green); got != "code=true" {
			t.Errorf("got %s; want code=true", got)
		}
	})
	t.Run("merge of main bringing code", func(t *testing.T) {
		r, base := newRepo(t)
		r.git("checkout", "-b", "pr")
		green := r.commit(map[string]string{".rhiza/pr": "1"})
		r.git("checkout", "main")
		r.commit(map[string]string{"other.go": "package main\n"})
		r.git("checkout", "pr")
		r.git("merge", "--no-ff", "-m", "merge main", "main")
		if got := r.decide(base, green); got != "code=true" {
			t.Errorf("got %s; want code=true", got)
		}
	})
	t.Run("no green commit and the PR leaves .rhiza/", func(t *testing.T) {
		r, base := newRepo(t)
		r.git("checkout", "-b", "pr")
		r.commit(map[string]string{".rhiza/b": "2"})
		r.commit(map[string]string{"main.go": "package main // pr\n"})
		if got := r.decide(base, ""); got != "code=true" {
			t.Errorf("got %s; want code=true", got)
		}
	})
	t.Run("no green commit and the PR stays in .rhiza/", func(t *testing.T) {
		r, base := newRepo(t)
		r.git("checkout", "-b", "pr")
		r.commit(map[string]string{".rhiza/b": "2"})
		if got := r.decide(base, ""); got != "code=false" {
			t.Errorf("got %s; want code=false", got)
		}
	})
}
