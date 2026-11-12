package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeEntry(t *testing.T, dir, cs, name, ops string) {
	t.Helper()
	d := filepath.Join(dir, cs)
	os.MkdirAll(d, 0o755)
	body := `{"body":{"ops":[` + ops + `]},"change_set":"` + cs + `"}`
	if err := os.WriteFile(filepath.Join(d, name+".json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func note(node, text string) string {
	return `{"op":"node.create","kind":"note","node":"` + node + `","fields":{"statement":"` + text + `"}}`
}

func discusses(link, from, to string) string {
	return `{"op":"link.create","type":"discusses","link":"` + link + `","from":"` + from + `","to":"` + to + `"}`
}

func TestReadMarks(t *testing.T) {
	dir := t.TempDir()
	writeEntry(t, dir, "CS1", "a", note("N1", "breaking: drops the --old flag")+","+discusses("L1", "N1", "CS1"))
	writeEntry(t, dir, "CS2", "b", note("N2", "Severe: data loss on land")+","+discusses("L2", "N2", "CS2")+","+
		note("N3", "just a remark")+","+discusses("L3", "N3", "CS2")+","+
		note("N4", "breaking changes ahead")+","+discusses("L4", "N4", "CS2"))
	// A note on a node (not a change set), a retired note and a note with a retired link.
	writeEntry(t, dir, "CS3", "c", note("N5", "breaking: on a node")+","+discusses("L5", "N5", "NODE1")+","+
		note("N6", "severe: retired")+","+discusses("L6", "N6", "CS3")+","+
		note("N7", "severe: link retired")+","+discusses("L7", "N7", "CS3"))
	writeEntry(t, dir, "CS3", "d", `{"op":"node.retire","node":"N6"},{"op":"link.retire","link":"L7"}`)

	got, err := ReadMarks(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]Marks{
		"CS1": {Breaking: []string{"drops the --old flag"}},
		"CS2": {Severe: []string{"data loss on land"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
	if !got["CS1"].IsBreaking() || got["CS1"].IsSevere() || !got["CS2"].IsSevere() {
		t.Fatal("IsBreaking/IsSevere wrong")
	}
}

func TestReadMarksEmptyLog(t *testing.T) {
	got, err := ReadMarks(t.TempDir())
	if err != nil || len(got) != 0 {
		t.Fatalf("%v %v", got, err)
	}
}
