package claudecode

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
)

func usageJSON(ts, id string, sidechain bool, in, out, create, read int) string {
	return fmt.Sprintf(`{"type":"assistant","isSidechain":%t,"timestamp":%q,"message":{"id":%q,"usage":{"input_tokens":%d,"output_tokens":%d,"cache_creation_input_tokens":%d,"cache_read_input_tokens":%d}}}`,
		sidechain, ts, id, in, out, create, read)
}

func writeLines(t *testing.T, path string, lines ...string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	s := ""
	for _, l := range lines {
		s += l + "\n"
	}
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestUsageRecords(t *testing.T) {
	base := t.TempDir()
	config := filepath.Join(base, "config")
	// A file outside the config dir that must never be read.
	writeLines(t, filepath.Join(base, "outside", "projects", "p", "x.jsonl"),
		usageJSON("2026-11-03T10:30:00Z", "outside", false, 999, 0, 0, 0))
	proj := filepath.Join(config, "projects")
	writeLines(t, filepath.Join(proj, "-a", "s1.jsonl"),
		usageJSON("2026-11-03T09:59:59Z", "before", false, 1000, 0, 0, 0),
		usageJSON("2026-11-03T10:10:00Z", "m1", false, 10, 20, 30, 5000),
		usageJSON("2026-11-03T10:10:00Z", "m1", false, 10, 20, 30, 5000),
		`not json`,
		usageJSON("2026-11-03T12:00:00Z", "after", false, 1000, 0, 0, 0),
		usageJSON("2026-11-03T10:40:00Z", "m3", false, 1, 1, 1, 1))
	writeLines(t, filepath.Join(proj, "-a", "s1", "subagents", "agent-1.jsonl"),
		usageJSON("2026-11-03T10:20:00Z", "h1", true, 100, 200, 0, 9000))
	writeLines(t, filepath.Join(proj, "-b", "s2.jsonl"),
		usageJSON("2026-11-03T10:05:00Z", "m2", false, 7, 0, 0, 0),
		`{"type":"user","timestamp":"2026-11-03T10:06:00Z"}`)
	// A link inside projects to a transcript outside the config dir is not
	// followed.
	if err := os.Symlink(filepath.Join(base, "outside", "projects", "p", "x.jsonl"), filepath.Join(proj, "-b", "link.jsonl")); err != nil {
		t.Fatal(err)
	}

	getenv := func(k string) string {
		if k == "CLAUDE_CONFIG_DIR" {
			return config
		}
		return ""
	}
	from := time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)
	to := time.Date(2026, 11, 3, 11, 0, 0, 0, time.UTC)
	got, err := Adapter{}.UsageRecords(getenv, from, to)
	if err != nil {
		t.Fatal(err)
	}
	at := func(m, s int) time.Time { return time.Date(2026, 11, 3, 10, m, s, 0, time.UTC) }
	want := []adapter.UsageRecord{
		{At: at(5, 0), Tokens: 7},
		{At: at(10, 0), Tokens: 60},
		{At: at(20, 0), Tokens: 300},
		{At: at(40, 0), Tokens: 3},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestUsageRecordsNoConfig(t *testing.T) {
	got, err := Adapter{}.UsageRecords(func(string) string { return "" }, time.Time{}, time.Now())
	if err != nil || len(got) != 0 {
		t.Fatalf("got %v, %v", got, err)
	}
}

var _ adapter.UsageReader = Adapter{}

// A file last written before the span is not read, and a changed file is
// read again rather than served from the cache.
func TestUsageRecordsSkipsOldAndRereadsChanged(t *testing.T) {
	config := t.TempDir()
	getenv := func(k string) string {
		if k == "CLAUDE_CONFIG_DIR" {
			return config
		}
		return ""
	}
	from := time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)
	to := time.Date(2026, 11, 3, 11, 0, 0, 0, time.UTC)
	old := filepath.Join(config, "projects", "-a", "old.jsonl")
	writeLines(t, old, usageJSON("2026-11-03T10:30:00Z", "o1", false, 5, 0, 0, 0))
	if err := os.Chtimes(old, from.Add(-time.Hour), from.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	live := filepath.Join(config, "projects", "-a", "live.jsonl")
	writeLines(t, live, usageJSON("2026-11-03T10:10:00Z", "l1", false, 1, 0, 0, 0))
	got, err := Adapter{}.UsageRecords(getenv, from, to)
	if err != nil {
		t.Fatal(err)
	}
	if want := []adapter.UsageRecord{{At: from.Add(10 * time.Minute), Tokens: 1}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	writeLines(t, live,
		usageJSON("2026-11-03T10:10:00Z", "l1", false, 1, 0, 0, 0),
		usageJSON("2026-11-03T10:20:00Z", "l2", false, 2, 0, 0, 0))
	got, err = Adapter{}.UsageRecords(getenv, from, to)
	if err != nil {
		t.Fatal(err)
	}
	want := []adapter.UsageRecord{{At: from.Add(10 * time.Minute), Tokens: 1}, {At: from.Add(20 * time.Minute), Tokens: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// A transcript that grows is read on from the end of its last complete
// line, with the same records a whole read gives: an unterminated last
// line counts until it is finished, and a rewritten file is read whole.
func TestUsageRecordsReadsOnlyAppended(t *testing.T) {
	config := t.TempDir()
	getenv := func(k string) string {
		if k == "CLAUDE_CONFIG_DIR" {
			return config
		}
		return ""
	}
	from := time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)
	to := time.Date(2026, 11, 3, 11, 0, 0, 0, time.UTC)
	path := filepath.Join(config, "projects", "-a", "grow.jsonl")
	var lines []string
	for i := range 50 {
		lines = append(lines, usageJSON(fmt.Sprintf("2026-11-03T10:%02d:00Z", i), fmt.Sprintf("g%d", i), false, 1, 0, 0, 0))
	}
	write := func(s string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	whole := func() []adapter.UsageRecord {
		t.Helper()
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var out []adapter.UsageRecord
		for _, r := range parseUsage(data) {
			out = append(out, r.UsageRecord)
		}
		return out
	}
	check := func(maxRead int64) {
		t.Helper()
		before := usageBytesRead.Load()
		got, err := Adapter{}.UsageRecords(getenv, from, to)
		if err != nil {
			t.Fatal(err)
		}
		if want := whole(); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %d records %v, want %d %v", len(got), got, len(want), want)
		}
		if read := usageBytesRead.Load() - before; maxRead >= 0 && read > maxRead {
			t.Fatalf("read %d bytes, want at most %d", read, maxRead)
		}
	}
	body := ""
	for _, l := range lines[:40] {
		body += l + "\n"
	}
	write(body)
	check(-1)
	// An appended line and an unterminated one: only the new bytes and the
	// mark before them are read.
	grown := body + lines[40] + "\n" + lines[41][:20]
	write(grown)
	check(int64(len(grown)-len(body)) + usageMark)
	// The unterminated line is finished.
	done := body + lines[40] + "\n" + lines[41] + "\n"
	write(done)
	check(int64(len(done)-len(body)-len(lines[40])-1) + usageMark)
	// A rewrite that keeps growing but changes what was read is read whole.
	rewritten := ""
	for _, l := range lines[1:45] {
		rewritten += l + "\n"
	}
	write(rewritten)
	check(-1)
	// A file cut shorter is read whole.
	write(lines[2] + "\n")
	check(-1)
}
