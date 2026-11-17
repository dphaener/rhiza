package guide

import (
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// kindEntries returns the Step kinds section of body as its top-level
// entries, "- `<kind>`: ...", each with its indented continuation lines,
// keyed by kind.
func kindEntries(body string) map[string]string {
	out := map[string]string{}
	kind := ""
	for _, l := range strings.Split(sectionOf(body, "## Step kinds"), "\n") {
		if rest, ok := strings.CutPrefix(l, "- `"); ok {
			if name, _, ok := strings.Cut(rest, "`"); ok {
				kind = name
				out[kind] = l
				continue
			}
		}
		if kind != "" && strings.HasPrefix(l, "  ") {
			out[kind] += "\n" + l
			continue
		}
		kind = ""
	}
	return out
}

// keysMissing returns each key, kind, sub-table and outcome the given
// tables document that body does not name where it belongs: a file key in
// File keys, a shared key in Shared step keys, a role key in Role keys, and
// each kind as its own entry in Step kinds naming its keys, sub-tables and
// outcomes.
func keysMissing(body string, file, shared, roles []workflow.KeyDoc, kinds []workflow.KindDoc) []string {
	var out []string
	named := func(section string, keys []workflow.KeyDoc, what string) {
		text := sectionOf(body, section)
		for _, k := range keys {
			if !strings.Contains(text, "`"+k.Key+"`") {
				out = append(out, what+" "+k.Key)
			}
		}
	}
	named("## File keys", file, "file key")
	named("## Shared step keys", shared, "shared step key")
	named("## Role keys", roles, "role key")
	entries := kindEntries(body)
	for _, d := range kinds {
		e, ok := entries[d.Kind]
		if !ok {
			out = append(out, "kind "+d.Kind)
			continue
		}
		for _, k := range d.Keys {
			if !strings.Contains(e, "`"+k.Key+"`") {
				out = append(out, "kind "+d.Kind+" key "+k.Key)
			}
		}
		for _, tb := range d.Tables {
			if !strings.Contains(e, tb.Key) {
				out = append(out, "kind "+d.Kind+" table "+tb.Key)
			}
		}
		for _, o := range d.Outcomes {
			if !strings.Contains(e, "`"+o+"`") {
				out = append(out, "kind "+d.Kind+" outcome "+o)
			}
		}
	}
	return out
}

func TestWorkflowKeys(t *testing.T) {
	wk, ok := Lookup("workflow-keys")
	if !ok || wk.Skill {
		t.Fatalf("workflow-keys: found %v, skill %v; want a topic", ok, wk.Skill)
	}
	if n := Tokens(wk.Body); n > TopicCap {
		t.Errorf("workflow-keys costs %d tokens, cap %d", n, TopicCap)
	}
	if !slices.ContainsFunc(List(), func(tx Text) bool { return tx.Name == "workflow-keys" }) {
		t.Errorf("workflow-keys is not listed")
	}
	for _, m := range keysMissing(wk.Body, workflow.FileKeys, workflow.SharedKeys, workflow.RoleKeys(), workflow.KindDocs()) {
		t.Errorf("workflow-keys names no %s", m)
	}
	checkPhrases(t, "workflow-keys", map[string][]string{
		"defaults": {
			"the file's base name", "the first step in the file", "unset 3", "means once",
			"the watchdog's default", "(unset: the run's)", "the agent's default",
			"the last stays", "10s, or 30s with `until`", "`idle`, the default",
			"`auto`, the default", "`rhiza.autonomy.<stage>`", "`cleanup-<n>`",
			"a built-in arbiter",
		},
		"conductor step": {"`[[step.<name>.check]]` is a conductor step"},
		"workflow new": {
			"rhiza workflow new <name> [--from <workflow>]", "With `--from`",
			"a copy of\nthe resolved workflow with its `name` set", "a skeleton",
			"agent step `build`", "command step `test`", "placeholder that fails on\npurpose",
			"`max_visits = 3`", "`on.fail = \"handback:build\"`",
			"refuses an existing file", "validates the result as\n`rhiza workflow check` does",
		},
	})
	if strings.Contains(wk.Body+wk.Description, "Claude") {
		t.Errorf("workflow-keys names an agent product")
	}
}

// The conduct topic's Read first and the workflows topic point to
// workflow-keys.
func TestConductReadsWorkflowKeys(t *testing.T) {
	c, _ := Lookup("conduct")
	if !strings.Contains(sectionOf(c.Body, "## Read first"), "`rhiza guide workflow-keys`") {
		t.Errorf("conduct: Read first does not name rhiza guide workflow-keys")
	}
	if n := Tokens(c.Body); n > TopicCap {
		t.Errorf("conduct costs %d tokens, cap %d", n, TopicCap)
	}
}

// A key or kind spec.go documents that the topic does not name is reported
// by name.
func TestWorkflowKeysCatchesDrift(t *testing.T) {
	wk, _ := Lookup("workflow-keys")
	file := append(slices.Clone(workflow.FileKeys), workflow.KeyDoc{Key: "timezone"})
	shared := append(slices.Clone(workflow.SharedKeys), workflow.KeyDoc{Key: "retries"})
	roles := append(workflow.RoleKeys(), workflow.KeyDoc{Key: "effort"})
	kinds := append(workflow.KindDocs(), workflow.KindDoc{Kind: "lint"})
	got := keysMissing(wk.Body, file, shared, roles, kinds)
	want := []string{"file key timezone", "shared step key retries", "role key effort", "kind lint"}
	if !slices.Equal(got, want) {
		t.Errorf("keysMissing with added keys = %q, want %q", got, want)
	}
	noGate := strings.Replace(wk.Body, "- `gate`:", "- `gat`:", 1)
	if got := keysMissing(noGate, workflow.FileKeys, workflow.SharedKeys, workflow.RoleKeys(), workflow.KindDocs()); !slices.Equal(got, []string{"kind gate"}) {
		t.Errorf("keysMissing without gate = %q, want [kind gate]", got)
	}
	noStall := strings.Replace(wk.Body, "- `stall`:", "- stall:", 1)
	if got := keysMissing(noStall, workflow.FileKeys, workflow.SharedKeys, workflow.RoleKeys(), workflow.KindDocs()); !slices.Equal(got, []string{"shared step key stall"}) {
		t.Errorf("keysMissing without stall = %q, want [shared step key stall]", got)
	}
}
