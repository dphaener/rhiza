package main

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/tui/requestapp"
)

type testRequest struct {
	Kind   string   `json:"kind" arg:"0"`
	Handle string   `json:"handle,omitempty" arg:"1"`
	Parent string   `json:"parent,omitempty"`
	Fields []string `json:"fields,omitempty" flag:"field"`
	Force  bool     `json:"force,omitempty"`
	Depth  int      `json:"depth,omitempty"`
	Hidden string   `json:"hidden,omitempty" flag:"-"`
}

type testResponse struct{}

var testOp = api.Operation{
	Name:     "test.op",
	Request:  reflect.TypeFor[testRequest](),
	Response: reflect.TypeFor[testResponse](),
	Run:      func(*api.Ctx, any) (any, error) { return testResponse{}, nil },
}

type restRequest struct {
	Files []string `json:"files" arg:"0"`
}

var restOp = api.Operation{
	Name:     "rest.op",
	Request:  reflect.TypeFor[restRequest](),
	Response: reflect.TypeFor[testResponse](),
	Run:      func(*api.Ctx, any) (any, error) { return testResponse{}, nil },
}

func TestRequestFromArgs(t *testing.T) {
	cases := []struct {
		name string
		op   api.Operation
		args []string
		want any
	}{
		{"positional", testOp, []string{"capability"}, testRequest{Kind: "capability"}},
		{"two positionals", testOp, []string{"capability", "store"}, testRequest{Kind: "capability", Handle: "store"}},
		{"flag space", testOp, []string{"c", "--parent", "p"}, testRequest{Kind: "c", Parent: "p"}},
		{"flag equals", testOp, []string{"--parent=p=q", "c"}, testRequest{Kind: "c", Parent: "p=q"}},
		{"repeated", testOp, []string{"c", "--field", "a=1", "--field=b=2"}, testRequest{Kind: "c", Fields: []string{"a=1", "b=2"}}},
		{"bool", testOp, []string{"c", "--force"}, testRequest{Kind: "c", Force: true}},
		{"bool false", testOp, []string{"c", "--force=false"}, testRequest{Kind: "c"}},
		{"int", testOp, []string{"c", "--depth", "3"}, testRequest{Kind: "c", Depth: 3}},
		{"dash value", testOp, []string{"c", "--parent", "-x"}, testRequest{Kind: "c", Parent: "-x"}},
		{"double dash", testOp, []string{"--", "--weird"}, testRequest{Kind: "--weird"}},
		{"stdin dash", testOp, []string{"-"}, testRequest{Kind: "-"}},
		{"rest", restOp, []string{"a", "b", "c"}, restRequest{Files: []string{"a", "b", "c"}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := requestFromArgs(c.op, c.args)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got %#v, want %#v", got, c.want)
			}
		})
	}
}

func TestRequestsAppArgsKeepDashLeadingText(t *testing.T) {
	o, ok := api.Lookup("request")
	if !ok {
		t.Fatal("no request op")
	}
	got, err := requestFromArgs(o, requestapp.RequestArgs("bug", "--json output is wrong"))
	if err != nil {
		t.Fatal(err)
	}
	want := api.RequestRequest{Text: "--json output is wrong", Kind: "bug"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestRequestFromArgsErrors(t *testing.T) {
	cases := []struct {
		name string
		op   api.Operation
		args []string
		want string
	}{
		{"unknown flag", testOp, []string{"c", "--bogus"}, "unknown flag --bogus"},
		{"hidden flag", testOp, []string{"c", "--hidden", "x"}, "unknown flag --hidden"},
		{"short flag", testOp, []string{"c", "-x"}, "unknown flag -x"},
		{"missing positional", testOp, nil, "missing <kind>"},
		{"too many", testOp, []string{"a", "b", "c"}, "too many arguments"},
		{"no value", testOp, []string{"c", "--parent"}, "--parent needs a value"},
		{"twice", testOp, []string{"c", "--parent", "a", "--parent", "b"}, "--parent given more than once"},
		{"bad bool", testOp, []string{"c", "--force=maybe"}, "--force"},
		{"bad int", testOp, []string{"c", "--depth", "x"}, "--depth"},
		{"missing rest", restOp, nil, "missing <files>"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := requestFromArgs(c.op, c.args)
			if api.ExitCode(err) != 2 {
				t.Fatalf("err = %v, exit %d, want 2", err, api.ExitCode(err))
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("err = %q, want it to contain %q", err, c.want)
			}
		})
	}
}

func TestParseGlobals(t *testing.T) {
	g, rest, err := parseGlobals([]string{"--json", "change", "--change", "C1", "new", "--actor=A", "x", "--dry-run", "--for", "F", "--record", "-h", "--", "--json"})
	if err != nil {
		t.Fatal(err)
	}
	want := globals{json: true, dryRun: true, record: true, help: true, change: "C1", actor: "A", forRef: "F", changeSet: true, actorSet: true, forSet: true}
	if g != want {
		t.Errorf("globals = %+v, want %+v", g, want)
	}
	if !slices.Equal(rest, []string{"change", "new", "x", "--", "--json"}) {
		t.Errorf("rest = %q", rest)
	}
	if _, _, err := parseGlobals([]string{"x", "--change"}); api.ExitCode(err) != 2 {
		t.Errorf("--change without value: err = %v", err)
	}
}

func TestMatchCommand(t *testing.T) {
	o, n, ok := matchCommand([]string{"describe", "describe", "x"})
	if !ok || o.Name != "describe" || n != 1 {
		t.Errorf("matchCommand = %s, %d, %v", o.Name, n, ok)
	}
	if _, _, ok := matchCommand([]string{"nope"}); ok {
		t.Error("matched nope")
	}
}

func TestCtxFromEnv(t *testing.T) {
	// newCtx finds the repo from the working directory; outside a temp dir it
	// would find this repo's own .rhiza and offset the test clock and IDs.
	t.Chdir(t.TempDir())
	env := map[string]string{
		"RHIZA_CHANGE":    "envchange",
		"RHIZA_ACTOR":     "envactor",
		"RHIZA_TEST_NOW":  "2026-10-24T12:00:00Z",
		"RHIZA_TEST_SEED": "7",
	}
	getenv := func(k string) string { return env[k] }
	ctx, err := newCtx(globals{change: "flagchange", changeSet: true, dryRun: true}, nil, getenv)
	if err != nil {
		t.Fatal(err)
	}
	if ctx.ChangeRef != "flagchange" || ctx.ActorRef != "envactor" || ctx.ForRef != "" || !ctx.DryRun {
		t.Errorf("ctx refs = %q %q %q dry %v", ctx.ChangeRef, ctx.ActorRef, ctx.ForRef, ctx.DryRun)
	}
	if got := ctx.Now().Format("2006-01-02T15:04:05Z07:00"); got != "2026-10-24T12:00:00Z" {
		t.Errorf("now = %s", got)
	}
	a, _ := ctx.IDs.New()
	ctx2, _ := newCtx(globals{}, nil, getenv)
	b, _ := ctx2.IDs.New()
	if a != b {
		t.Errorf("seeded IDs differ: %s %s", a, b)
	}
	env["RHIZA_TEST_NOW"] = "yesterday"
	if _, err := newCtx(globals{}, nil, getenv); api.ExitCode(err) != 2 {
		t.Errorf("bad RHIZA_TEST_NOW: err = %v", err)
	}
}

func TestFindRoot(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "a", "b")
	if err := os.MkdirAll(filepath.Join(dir, ".rhiza"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := findRoot(sub); got != dir {
		t.Errorf("findRoot(sub) = %s, want %s", got, dir)
	}
	old := t.TempDir()
	if err := os.Mkdir(filepath.Join(old, ".genesis"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := findRoot(old); got != old {
		t.Errorf("findRoot(old-style root) = %s, want %s", got, old)
	}
	other := t.TempDir()
	if got := findRoot(other); got != other {
		t.Errorf("findRoot(no repo) = %s, want %s", got, other)
	}
}

// In test mode each run offsets the clock and seed by the number of entry
// files, so successive runs make different IDs, and the same log the same.
func TestCtxTestIDsFollowTheLog(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	env := map[string]string{"RHIZA_TEST_NOW": "2026-10-24T12:00:00Z", "RHIZA_TEST_SEED": "1"}
	getenv := func(k string) string { return env[k] }
	first := func() (id.ID, string) {
		t.Helper()
		ctx, err := newCtx(globals{}, nil, getenv)
		if err != nil {
			t.Fatal(err)
		}
		i, err := ctx.IDs.New()
		if err != nil {
			t.Fatal(err)
		}
		return i, ctx.Now().Format(time.RFC3339)
	}
	a0, now0 := first()
	b0, _ := first()
	if a0 != b0 {
		t.Errorf("same log, different IDs: %s %s", a0, b0)
	}
	cs := filepath.Join(dir, ".rhiza", "log", "CS")
	if err := os.MkdirAll(cs, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"E1.json", ".tmp.json"} {
		if err := os.WriteFile(filepath.Join(cs, name), []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	a1, now1 := first()
	b1, _ := first()
	if a1 == a0 {
		t.Errorf("one more entry file, same first ID %s", a1)
	}
	if a1 != b1 {
		t.Errorf("same log, different IDs: %s %s", a1, b1)
	}
	if now0 != "2026-10-24T12:00:00Z" || now1 != "2026-10-24T12:00:01Z" {
		t.Errorf("clocks = %s, %s", now0, now1)
	}
	// repo.json counts as one more unit, so the first write after init
	// does not reuse init's IDs.
	if err := os.WriteFile(filepath.Join(dir, ".rhiza", "repo.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	a2, now2 := first()
	if a2 == a1 || a2 == a0 {
		t.Errorf("repo.json added, same first ID %s", a2)
	}
	if now2 != "2026-10-24T12:00:02Z" {
		t.Errorf("clock with repo.json = %s", now2)
	}
}
