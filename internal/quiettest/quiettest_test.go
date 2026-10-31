package quiettest

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"
)

// TestTestProcesses reads ps output as Linux (procps) and darwin print it
// for the same go test run: this process is pid 300 under the go command,
// pid 200.
func TestTestProcesses(t *testing.T) {
	const ps = `    1     0 /sbin/init splash
  100     1 zsh
  200   100 go test ./...
  300   200 /tmp/go-build1/b001/api.test -test.paniconexit0 -test.timeout=10m0s
  301   200 /tmp/go-build1/b002/slice.test -test.paniconexit0 -test.timeout=10m0s
  302   200 /usr/local/go/pkg/tool/linux_amd64/compile -o /tmp/go-build1/b003/_pkg_.a -trimpath /tmp/go-build1/b003=> -p index
  303   200 /usr/local/go/pkg/tool/linux_amd64/vet -atomic -bool -buildtags
  304   200 /usr/local/go/pkg/tool/linux_amd64/link -o /tmp/go-build1/b004/gen.test
  305   301 /tmp/go-build1/b002/slice.test -test.run=^TestChild$
  306   100 tail -5
  307   200 git status
  308   200 /private/var/folders/x/T/go-build2/b005/changeset.test -test.paniconexit0
bad line
`
	got := testProcesses(ps, 300, 200, nil)
	want := []string{"slice.test", "compile", "vet", "link", "changeset.test"}
	if !slices.Equal(got, want) {
		t.Errorf("testProcesses = %q, want %q", got, want)
	}
	// 308 is another waiter: blocked on the lock or timing, not running.
	want = []string{"slice.test", "compile", "vet", "link"}
	if got := testProcesses(ps, 300, 200, map[int]bool{308: true}); !slices.Equal(got, want) {
		t.Errorf("testProcesses beside waiter 308 = %q, want %q", got, want)
	}
	if got := testProcesses(ps, 300, 999, nil); got != nil {
		t.Errorf("testProcesses under another parent = %q, want none", got)
	}
	if got := testProcesses("", 300, 200, nil); got != nil {
		t.Errorf("testProcesses of nothing = %q, want none", got)
	}
}

func TestQuietString(t *testing.T) {
	cases := []struct {
		q    Quiet
		want string
	}{
		{Quiet{waited: 1750 * time.Millisecond}, "alone after 1.8s"},
		{Quiet{waited: 41340 * time.Millisecond, last: []string{"slice.test"}}, "alone after 41.3s (last beside slice.test)"},
		{Quiet{waited: 3 * time.Minute, last: []string{"index.test", "slice.test"}, running: []string{"slice.test"}}, "beside slice.test after 3m0s"},
		{Quiet{err: fmt.Errorf("ps: exec: no ps")}, "other packages unknown (ps: exec: no ps)"},
	}
	for _, c := range cases {
		if got := c.q.String(); got != c.want {
			t.Errorf("%+v.String() = %q, want %q", c.q, got, c.want)
		}
	}
}

// TestTakeTurnOneAtATime checks that a second waiter blocks on the lock
// until the first one's test ends, and that each leaves its mark only
// while it waits or times.
func TestTakeTurnOneAtATime(t *testing.T) {
	dir := t.TempDir()
	oldDir, oldLock := waiterDir, waiterLock
	waiterDir, waiterLock = filepath.Join(dir, "quiet"), filepath.Join(dir, "quiet.lock")
	t.Cleanup(func() { waiterDir, waiterLock = oldDir, oldLock })

	var firstEnd, secondGot time.Time
	started, release := make(chan struct{}), make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		t.Run("first", func(t *testing.T) {
			if err := takeTurn(t); err != nil {
				t.Fatal(err)
			}
			if !waiters()[os.Getpid()] {
				t.Error("first waiter not marked")
			}
			close(started)
			<-release
			firstEnd = time.Now()
		})
	}()
	<-started
	go func() {
		defer wg.Done()
		t.Run("second", func(t *testing.T) {
			if err := takeTurn(t); err != nil {
				t.Fatal(err)
			}
			secondGot = time.Now()
		})
	}()
	time.Sleep(200 * time.Millisecond)
	close(release)
	wg.Wait()
	if !secondGot.After(firstEnd) {
		t.Errorf("second waiter took its turn at %v, before the first ended at %v", secondGot, firstEnd)
	}
	if len(waiters()) != 0 {
		t.Errorf("waiters after both ended = %v, want none", waiters())
	}
}
