package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNextVersion(t *testing.T) {
	br := []string{"drops --old"}
	for _, c := range []struct {
		name, tag string
		in        Inputs
		want      string
		kind      Kind
	}{
		{"breaking from v1 is major", "v1.2.3", Inputs{Breaking: br, Milestone: []string{"m"}}, "v2.0.0", KindMajor},
		{"breaking before v1 is minor", "v0.4.2", Inputs{Breaking: br}, "v0.5.0", KindMinor},
		{"milestone is minor", "v1.2.3", Inputs{Milestone: []string{"m"}, Fixes: []string{"f"}}, "v1.3.0", KindMinor},
		{"fixes are a patch", "v1.2.3", Inputs{Fixes: []string{"f"}, Landed: 1}, "v1.2.4", KindPatch},
		{"features alone are none", "v1.2.3", Inputs{Landed: 2}, "", KindNone},
		{"nothing since the tag", "v1.2.3", Inputs{}, "", KindNone},
		{"no tag, breaking is minor", "", Inputs{Breaking: br}, "v0.1.0", KindMinor},
	} {
		t.Run(c.name, func(t *testing.T) {
			n, err := NextVersion(c.tag, c.in)
			if err != nil {
				t.Fatal(err)
			}
			if n.Version != c.want || n.Kind != c.kind {
				t.Fatalf("got %s %s, want %s %s", n.Version, n.Kind, c.want, c.kind)
			}
		})
	}
	if _, err := NextVersion("1.2", Inputs{}); err == nil {
		t.Fatal("bad tag accepted")
	}
}

func writeAt(t *testing.T, dir, cs, name, typ, at, body string) {
	t.Helper()
	d := filepath.Join(dir, cs)
	os.MkdirAll(d, 0o755)
	s := `{"at":"` + at + `","body":` + body + `,"change_set":"` + cs + `","type":"` + typ + `"}`
	if err := os.WriteFile(filepath.Join(d, name+".json"), []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

// fixture: CS1 (feature, landed before the tag), CS2 (bug, after), CS3 and
// CS4 serve milestone m2, CS3 landed before the tag, CS4 after; CS5 is
// marked breaking and landed after; CS6 serves m3 and has not landed.
func fixture(t *testing.T) (string, time.Time) {
	dir := t.TempDir()
	before, after := "2026-11-01T00:00:00Z", "2026-11-10T00:00:00Z"
	writeAt(t, dir, "CS1", "a", "open", before, `{"kind":"intention"}`)
	writeAt(t, dir, "CS1", "b", "land", before, `{}`)
	writeAt(t, dir, "CS2", "a", "open", after, `{"kind":"intention","request_kind":"bug"}`)
	writeAt(t, dir, "CS2", "b", "land", after, `{}`)
	writeAt(t, dir, "CS3", "a", "open", before, `{"kind":"intention"}`)
	writeAt(t, dir, "CS3", "b", "ops", before, `{"ops":[{"op":"node.create","kind":"intent","node":"I2","handle":"m2"},
		{"op":"node.create","kind":"intent","node":"I3","handle":"m3"},
		{"op":"link.create","type":"serves","from":"CS3","to":"I2"},
		{"op":"link.create","type":"serves","from":"CS4","to":"I2"},
		{"op":"link.create","type":"serves","from":"CS6","to":"I3"}]}`)
	writeAt(t, dir, "CS3", "c", "land", before, `{}`)
	writeAt(t, dir, "CS4", "a", "open", after, `{"kind":"intention"}`)
	writeAt(t, dir, "CS4", "b", "land", after, `{}`)
	writeAt(t, dir, "CS5", "a", "open", after, `{"kind":"intention"}`)
	writeAt(t, dir, "CS5", "b", "land", after, `{}`)
	writeAt(t, dir, "CS6", "a", "open", after, `{"kind":"intention"}`)
	return dir, time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)
}

func TestNextVersionFromLogs(t *testing.T) {
	dir, tag := fixture(t)
	h, err := ReadHistory(dir)
	if err != nil {
		t.Fatal(err)
	}
	in := h.Inputs(tag, nil)
	if in.Landed != 3 || len(in.Fixes) != 1 || in.Fixes[0] != "CS2" || len(in.Milestone) != 1 || in.Milestone[0] != "m2" {
		t.Fatalf("inputs %+v", in)
	}
	n, _ := NextVersion("v0.1.0", in)
	if n.Version != "v0.2.0" || n.Kind != KindMinor {
		t.Fatalf("%+v", n)
	}
	// A breaking mark on a landed change set: minor before v1, major after.
	marks := map[string]Marks{"CS5": {Breaking: []string{"drops a flag"}}, "CS6": {Breaking: []string{"unlanded"}}}
	in = h.Inputs(tag, marks)
	if len(in.Breaking) != 1 {
		t.Fatalf("breaking %v", in.Breaking)
	}
	if n, _ := NextVersion("v1.0.0", in); n.Version != "v2.0.0" {
		t.Fatalf("%+v", n)
	}
	if n, _ := NextVersion("v0.1.0", in); n.Version != "v0.2.0" {
		t.Fatalf("%+v", n)
	}
	// No change since the tag.
	in = h.Inputs(time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC), marks)
	if n, _ := NextVersion("v1.0.0", in); n.Kind != KindNone || n.Version != "" {
		t.Fatalf("%+v", n)
	}
}
