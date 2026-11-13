package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestQueueWaitRoundTrip(t *testing.T) {
	dir := t.TempDir()
	at := time.Date(2026, 11, 13, 10, 0, 0, 0, time.UTC)
	want := []QueueWait{
		{At: at, Change: "C1", Title: "One", Reason: "no-slot", Detail: "2 of 2 busy"},
		{At: at.Add(time.Minute), Change: "C1", Title: "One", Reason: QueueWaitStarted},
	}
	for _, w := range want {
		if err := RecordQueueWaitIn(dir, w); err != nil {
			t.Fatal(err)
		}
	}
	got := ReadQueueWaitsIn(dir)
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestQueueWaitSkipsMalformed(t *testing.T) {
	dir := t.TempDir()
	a := QueueWait{At: time.Unix(1, 0).UTC(), Change: "A", Title: "a", Reason: "no-slot"}
	b := QueueWait{At: time.Unix(2, 0).UTC(), Change: "B", Title: "b", Reason: "started"}
	if err := RecordQueueWaitIn(dir, a); err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(filepath.Join(dir, QueueWaitFile), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("{not json\n")
	f.Close()
	if err := RecordQueueWaitIn(dir, b); err != nil {
		t.Fatal(err)
	}
	got := ReadQueueWaitsIn(dir)
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Fatalf("got %+v", got)
	}
}

func TestQueueWaitMissingFile(t *testing.T) {
	if got := ReadQueueWaitsIn(t.TempDir()); len(got) != 0 {
		t.Fatalf("got %+v", got)
	}
	if got := ReadQueueWaits(t.TempDir()); len(got) != 0 {
		t.Fatalf("got %+v", got)
	}
}

func TestQueueWaitConcurrentAppends(t *testing.T) {
	common := t.TempDir()
	const n = 50
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := RecordQueueWait(common, QueueWait{At: time.Unix(int64(i), 0).UTC(), Change: fmt.Sprint("C", i), Title: "t", Reason: "no-slot", Detail: "d"})
			if err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
	got := ReadQueueWaits(common)
	if len(got) != n {
		t.Fatalf("read %d lines, want %d", len(got), n)
	}
	seen := map[string]bool{}
	for _, w := range got {
		seen[w.Change] = true
	}
	if len(seen) != n {
		t.Fatalf("distinct changes %d, want %d", len(seen), n)
	}
}
