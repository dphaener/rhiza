package op

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestKnownFieldsSortedAndCoverEveryKind(t *testing.T) {
	for _, k := range Kinds() {
		if _, ok := knownFields[k]; !ok {
			t.Errorf("kind %s has no field list", k)
		}
		if !slices.IsSorted(KnownFields(k)) {
			t.Errorf("kind %s fields are not sorted", k)
		}
	}
	if KnownFields("nope") != nil {
		t.Error("unknown kind should have no fields")
	}
}

func TestBogusFieldUnknown(t *testing.T) {
	if FieldKnown(KindEvidence, "bogus") {
		t.Error("bogus must be unknown on evidence")
	}
	if !FieldKnown(KindEvidence, "result") {
		t.Error("result must be known on evidence")
	}
}

// collect walks decoded JSON for node.create and node.update ops.
func collect(v any, kinds map[string]Kind, used map[Kind]map[string]bool) {
	switch x := v.(type) {
	case []any:
		for _, e := range x {
			collect(e, kinds, used)
		}
	case map[string]any:
		switch x["op"] {
		case "node.create":
			k, _ := x["kind"].(string)
			if id, _ := x["node"].(string); id != "" {
				kinds[id] = Kind(k)
			}
			f, _ := x["fields"].(map[string]any)
			for name := range f {
				if used[Kind(k)] == nil {
					used[Kind(k)] = map[string]bool{}
				}
				used[Kind(k)][name] = true
			}
		case "node.update":
			id, _ := x["node"].(string)
			set, _ := x["set"].(map[string]any)
			for name := range set {
				if k, ok := kinds[id]; ok {
					if used[k] == nil {
						used[k] = map[string]bool{}
					}
					used[k][name] = true
				}
			}
		}
		for _, e := range x {
			collect(e, kinds, used)
		}
	}
}

func TestEveryFieldInRepoLogIsKnown(t *testing.T) {
	files, err := filepath.Glob("../../.rhiza/log/*/*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Skip("no log in this checkout")
	}
	kinds := map[string]Kind{}
	used := map[Kind]map[string]bool{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var v any
		if err := json.Unmarshal(data, &v); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		collect(v, kinds, used)
	}
	for k, names := range used {
		for name := range names {
			if !FieldKnown(k, name) {
				t.Errorf("%s field %q is in the log but not known", k, name)
			}
		}
	}
	if len(used) == 0 {
		t.Error("walked the log and found no node fields")
	}
}

func TestNoteFields(t *testing.T) {
	for _, f := range []string{"statement", "op"} {
		if !FieldKnown(KindNote, f) {
			t.Errorf("note field %s unknown", f)
		}
	}
	if FieldKnown(KindNote, "bogus") {
		t.Error("bogus known")
	}
}
