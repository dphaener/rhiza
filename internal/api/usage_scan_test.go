package api

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/adapter/claudecode"
)

// TestUsageScanOncePerPoll: the usage governor walks the transcript tree at
// most once a minute and rereads only files that changed, however many
// projects poll.
func TestUsageScanOncePerPoll(t *testing.T) {
	holdFakeRecords = nil
	config := t.TempDir()
	getenv := func(k string) string {
		if k == "CLAUDE_CONFIG_DIR" {
			return config
		}
		return ""
	}
	base := time.Now().UTC().Truncate(time.Second)
	line := func(id string, at time.Time) string {
		return fmt.Sprintf(`{"type":"assistant","timestamp":%q,"message":{"id":%q,"usage":{"input_tokens":10,"output_tokens":20,"cache_creation_input_tokens":5}}}`+"\n", at.Format(time.RFC3339), id)
	}
	write := func(name string, appendTo bool, text string) {
		path := filepath.Join(config, "projects", "-p", name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		if appendTo {
			flag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
		}
		f, err := os.OpenFile(path, flag, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString(text); err != nil {
			t.Fatal(err)
		}
		f.Close()
		// Distinct mtimes make a change visible however fast the writes.
		m := base.Add(time.Duration(len(text)) * time.Millisecond)
		if appendTo {
			m = time.Now().Add(time.Hour)
		}
		os.Chtimes(path, m, m)
	}
	const files = 120
	for i := range files {
		write(fmt.Sprintf("s%03d.jsonl", i), false, line(fmt.Sprintf("m%d", i), base.Add(-time.Duration(i)*time.Minute)))
	}

	now := base
	f := usageFix(t)
	f.ctx.Getenv = getenv
	f.ctx.Now = func() time.Time { return now }
	from, to := base.Add(-24*time.Hour), base.Add(time.Hour)
	// A root poll reads the records for each of two projects.
	poll := func() (walks, reads int64, n int) {
		w0, r0 := claudecode.UsageScanStats()
		for range 2 {
			n = len(usageRecordsFor(f.ctx, from, to))
		}
		w1, r1 := claudecode.UsageScanStats()
		return w1 - w0, r1 - r0, n
	}
	check := func(what string, wantWalks, wantReads int64, wantN int) {
		t.Helper()
		w, r, n := poll()
		if w != wantWalks || r != wantReads || n != wantN {
			t.Fatalf("%s: walks %d reads %d records %d, want %d %d %d", what, w, r, n, wantWalks, wantReads, wantN)
		}
	}
	check("first poll", 1, files, files)
	now = now.Add(10 * time.Second)
	check("nothing changed", 0, 0, files)
	write("s005.jsonl", true, line("grown", base))
	now = now.Add(10 * time.Second)
	check("one file grew", 0, 1, files+1)
	write("new.jsonl", false, line("fresh", base))
	now = now.Add(10 * time.Second)
	check("new file waits for the walk", 0, 0, files+1)
	now = now.Add(time.Minute)
	check("next full walk", 1, 1, files+2)

	want, err := claudecode.Adapter{}.UsageRecords(getenv, from, to)
	if err != nil {
		t.Fatal(err)
	}
	var got []adapter.UsageRecord
	got = usageRecordsFor(f.ctx, from, to)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("records differ from a full read: %d vs %d", len(got), len(want))
	}
}
