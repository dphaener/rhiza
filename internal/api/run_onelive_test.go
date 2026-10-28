package api

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// Many starts racing on one change set: exactly one begins, the rest are
// refused with exit 2 naming a run of it; another change set is not held
// back.
func TestConcurrentStartsOneBegins(t *testing.T) {
	src := "[step.hold]\nkind = \"wait\"\nuntil = \"false\"\ntimeout = \"1h\"\n"
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(filepath.Join(t.TempDir(), "runs"))
	now := time.Unix(1700000000, 0)
	const n = 8
	errs := make([]error, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = startRun(rs, fmt.Sprintf("RUN%d", i), run.Started{Workflow: "wf", Definition: src, Change: "CS1"}, now, def)
		}()
	}
	wg.Wait()
	began := 0
	for _, err := range errs {
		if err == nil {
			began++
			continue
		}
		wantExit(t, err, 2)
	}
	if began != 1 {
		t.Fatalf("%d starts began, want 1: %v", began, errs)
	}
	runs, err := rs.List()
	if err != nil || len(runs) != 1 {
		t.Fatalf("runs = %d, %v; want 1", len(runs), err)
	}
	if err := startRun(rs, "OTHER", run.Started{Workflow: "wf", Definition: src, Change: "CS2"}, now, def); err != nil {
		t.Errorf("start on another change set: %v", err)
	}
}
