package guide

import (
	"strings"
	"testing"
)

func sharedFiles(t *testing.T) map[string]string {
	t.Helper()
	ents, err := files.ReadDir("text/shared")
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, e := range ents {
		b, err := files.ReadFile("text/shared/" + e.Name())
		if err != nil {
			t.Fatal(err)
		}
		out[e.Name()] = string(b)
	}
	return out
}

func TestSharedBlockExpandsInPlace(t *testing.T) {
	shared := sharedFiles(t)
	if len(shared) == 0 {
		t.Fatal("no shared blocks")
	}
	for name, body := range shared {
		block := strings.TrimSuffix(name, ".md")
		got, err := Expand("a\n<!-- shared: "+block+" -->\nb\n", loadShared)
		if err != nil {
			t.Fatal(err)
		}
		if want := "a\n" + strings.TrimRight(body, "\n") + "\nb\n"; got != want {
			t.Errorf("%s: expanded %q, want %q", block, got, want)
		}
	}
	for _, tx := range List() {
		if strings.Contains(tx.Body, "<!-- shared:") {
			t.Errorf("%s: a shared marker is left", tx.Name)
		}
	}
	b, _ := Lookup("build")
	if !strings.Contains(b.Body, strings.TrimRight(shared["command-budget.md"], "\n")) {
		t.Errorf("build lacks the command-budget block")
	}
}

func TestSharedNestedBlockDemotesHeadings(t *testing.T) {
	block := func(string) (string, bool) { return "## Title\nbody\n### Sub\n", true }
	got, err := Expand("## Section\n<!-- shared: x nested -->\nb\n", block)
	if err != nil {
		t.Fatal(err)
	}
	if want := "## Section\n### Title\nbody\n#### Sub\nb\n"; got != want {
		t.Errorf("nested expand = %q, want %q", got, want)
	}
}

func TestSharedUnknownBlockFails(t *testing.T) {
	_, err := Expand("x\n<!-- shared: no-such-block -->\n", loadShared)
	if err == nil || !strings.Contains(err.Error(), "no-such-block") {
		t.Fatalf("Expand = %v, want an error naming no-such-block", err)
	}
}

func TestSharedTextsWithinCap(t *testing.T) {
	for _, tx := range List() {
		if err := checkCap(tx); err != nil {
			t.Error(err)
		}
	}
}

func TestSharedNeverList(t *testing.T) {
	n := strings.ToLower(sharedFiles(t)["never-list.md"])
	for _, want := range []string{"sleep or poll", "stubs, todos or placeholders", "new dependency", "rewrite history", "chain `session done`"} {
		if !strings.Contains(n, want) {
			t.Errorf("never-list lacks %q", want)
		}
	}
}

func TestSharedPhrasesInExactlyOneFile(t *testing.T) {
	shared := sharedFiles(t)
	for _, p := range []string{
		"## Command budget", "Never sleep and never poll", "about four\nminutes",
		"runs once in the background and is\nwaited for once",
		"Builders and verifiers run only the tests covering what the change touches",
		"scoped commands", "never the repo's full\nsuite", "the PR's CI runs it before merge",
		"While CI is unstable", "not explained by your change", "do not push or rerun CI to see",
		"Reproduce the\nfailing test locally", "push only after it and the\ntouched packages pass locally",
		"injecting a delay at the suspected race",
		"raise a question", "linked `blocks` to the leaf",
		"discovered-from draft change set, never built in this run",
		"rhiza change new", "discovered-from <this>", "rhiza change use <this>",
		"rhiza session done", "as its own command", "after its recording\ncommand succeeded", "`;` or `&&`",
	} {
		n := 0
		for _, b := range shared {
			if strings.Contains(b, p) {
				n++
			}
		}
		if n != 1 {
			t.Errorf("%q is in %d shared files, want 1", p, n)
		}
	}
}
