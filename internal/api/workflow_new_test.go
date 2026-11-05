package api

import (
	"bytes"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

func TestWithNameDiffersOnlyInName(t *testing.T) {
	for _, n := range workflow.BuiltinNames() {
		src, _ := workflow.Builtin(n)
		got := withName(src, "forked")
		if n == "forked" || bytes.Equal(got, src) {
			t.Fatalf("%s: unchanged", n)
		}
		back := withName(got, n)
		if !bytes.Equal(back, src) {
			t.Errorf("%s: renaming back differs from the built-in", n)
		}
	}
}

func TestWithNameAddsMissing(t *testing.T) {
	got := string(withName([]byte("# c\n[step.a]\nkind = \"x\"\n"), "n"))
	want := "name = \"n\"\n\n# c\n[step.a]\nkind = \"x\"\n"
	if got != want {
		t.Errorf("got %q", got)
	}
}
