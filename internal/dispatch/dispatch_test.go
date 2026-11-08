package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultsWhenNoFile(t *testing.T) {
	r, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range Kinds {
		if r.ForRequest(k) != "plan" || r.ForReady(k) != "ship-change" {
			t.Errorf("%s: got %q/%q", k, r.ForRequest(k), r.ForReady(k))
		}
	}
}

func TestReadyBugOnly(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".rhiza"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, File), []byte("[ready]\nbug = \"bugfix\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	d := Defaults()
	for _, k := range Kinds {
		want := d.ForReady(k)
		if k == KindBug {
			want = "bugfix"
		}
		if r.ForReady(k) != want || r.ForRequest(k) != "plan" {
			t.Errorf("%s: got %q/%q", k, r.ForRequest(k), r.ForReady(k))
		}
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct{ in, want string }{
		{"[other]\nbug = \"x\"\n", "[other]"},
		{"[request]\nchore = \"x\"\n", "[request] unknown key \"chore\""},
		{"[ready]\nmilestone = \"x\"\n", "[ready] unknown key \"milestone\": want one of feature, bug"},
		{"[ready]\nbug = 3\n", "[ready] bug"},
		{"[ready]\nbug = \"a/b\"\n", "[ready] bug"},
		{"[request]\nfeature = \"x.toml\"\n", "[request] feature"},
		{"[request]\nfeature = \"\"\n", "[request] feature"},
	}
	for _, c := range cases {
		_, err := Parse([]byte(c.in))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: got %v, want %q", c.in, err, c.want)
		}
	}
}

func TestSchedulesAbsent(t *testing.T) {
	for _, in := range []string{"", "[request]\nbug = \"plan\"\n"} {
		r, err := Parse([]byte(in))
		if err != nil || len(r.Schedules) != 0 {
			t.Errorf("%q: %v %v", in, r.Schedules, err)
		}
	}
}

func TestSchedulesRead(t *testing.T) {
	r, err := Parse([]byte(`[request]
bug = "plan"

[[schedules]]
name = "a"
every = "15m"
workflow = "w"

[[schedules]]
name = "b"
cron = "0 3 * * *"
command = "make x"
catch_up = "skip"
overlap = "queue"
`))
	if err != nil {
		t.Fatal(err)
	}
	want := []Schedule{
		{Name: "a", Every: "15m", Workflow: "w", CatchUp: "once", Overlap: "skip"},
		{Name: "b", Cron: "0 3 * * *", Command: "make x", CatchUp: "skip", Overlap: "queue"},
	}
	if len(r.Schedules) != 2 || r.Schedules[0] != want[0] || r.Schedules[1] != want[1] {
		t.Errorf("got %+v", r.Schedules)
	}
	if r.Request["bug"] != "plan" {
		t.Errorf("request lost")
	}
}

func TestSchedulesErrors(t *testing.T) {
	ok := "name = \"s\"\nevery = \"5m\"\nworkflow = \"w\"\n"
	cases := []struct{ name, in, want string }{
		{"unknown key", "[[schedules]]\n" + ok + "color = \"x\"\n", `schedule "s": unknown key "color"`},
		{"duplicate", "[[schedules]]\n" + ok + "[[schedules]]\n" + ok, `schedule "s": name: duplicate`},
		{"every and cron", "[[schedules]]\n" + ok + "cron = \"* * * * *\"\n", `schedule "s": every and cron`},
		{"neither when", "[[schedules]]\nname = \"s\"\nworkflow = \"w\"\n", `schedule "s": every, cron`},
		{"workflow and command", "[[schedules]]\n" + ok + "command = \"x\"\n", `schedule "s": workflow and command`},
		{"neither what", "[[schedules]]\nname = \"s\"\nevery = \"5m\"\n", `schedule "s": workflow, command`},
		{"bad cron", "[[schedules]]\nname = \"s\"\ncron = \"nope\"\nworkflow = \"w\"\n", `schedule "s": cron`},
		{"every 30s", "[[schedules]]\nname = \"s\"\nevery = \"30s\"\nworkflow = \"w\"\n", `schedule "s": every`},
		{"catch_up all", "[[schedules]]\n" + ok + "catch_up = \"all\"\n", `schedule "s": catch_up`},
		{"overlap parallel", "[[schedules]]\n" + ok + "overlap = \"parallel\"\n", `schedule "s": overlap`},
		{"unknown table", "[other]\nx = 1\n", `unknown table [other]`},
	}
	for _, c := range cases {
		_, err := Parse([]byte(c.in))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want %q", c.name, err, c.want)
		}
	}
}
