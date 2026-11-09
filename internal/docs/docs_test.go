package docs

import (
	"bytes"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	cmds := []Command{
		{Name: "request", Summary: "Queue a request.", Usage: "rhiza request <text> [--kind <kind>] [--milestone <milestone>]",
			Params: []Param{
				{Name: "<text>", Desc: "What to plan", Notes: []string{"required"}},
				{Name: "--kind", Desc: "feature, bug or milestone", Notes: []string{"one of feature, bug, milestone"}},
				{Name: "--milestone", Desc: "A milestone"},
			},
			Output:   []Field{{Name: "id", Type: "string", Desc: "The request"}},
			SetsExit: true},
		{Name: "add", Summary: "Add", Usage: "rhiza add", Writes: true},
	}
	got := CLI(cmds)
	if !bytes.Equal(got, CLI([]Command{cmds[1], cmds[0]})) {
		t.Error("output depends on input order")
	}
	s := string(got)
	if strings.Index(s, "## rhiza add") > strings.Index(s, "## rhiza request") {
		t.Error("commands not sorted")
	}
	for _, want := range []string{"- `--kind`: feature, bug or milestone (one of feature, bug, milestone)", "- `--milestone`: A milestone", "- `4`: Not found", "(#rhiza-request)", "- `id` (string): The request", "Appends to the log.", "Exit code: set by the result"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in:\n%s", want, s)
		}
	}
}
