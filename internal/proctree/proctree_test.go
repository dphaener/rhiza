package proctree

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

// table is fake ps output: two sessions' trees, each nested three deep
// below its root, beside unrelated processes, a session whose ID only
// shares a prefix with sessA, a Rex tab passing a session run line on, a
// header line and malformed lines.
const table = `    PID    PPID   RSS COMMAND
      1       0  10836 /usr/lib/systemd/systemd --switched-root
    100       1   2000 /usr/bin/rex daemon
    200     100   1000 /home/u/bin/rhiza session run sessA
    201     200  50000 claude --model opus
    202     201   3000 /bin/zsh -c go test ./internal/x
    203     202   7000 /tmp/go-build/x.test -test.v
    204     200    500 /bin/sh -c tail -f log
    300     100   1500 /home/u/bin/rhiza-head session run sessB --json
    301     300  40000 claude --model sonnet
    302     301   2500 /bin/zsh -c make gate
    303     302   9000 go test ./...
    400     100   1200 /home/u/bin/rhiza session run sessA2
    401     400  60000 claude
    500     100    800 rex tab new -- /home/u/bin/rhiza session run sessA
    600       1    900 /usr/sbin/sshd -D
garbage line
    700     abc   100 broken ppid
    701     1
    702     1    xx  bad rss
`

func pids(t Tree) []int {
	var out []int
	for _, p := range t.Processes {
		out = append(out, p.PID)
	}
	return out
}

func TestSessionTree(t *testing.T) {
	procs := Parse(table)
	for _, tc := range []struct {
		session string
		root    int
		pids    []int
		rss     int64
	}{
		{"sessA", 200, []int{200, 201, 202, 203, 204}, (1000 + 50000 + 3000 + 7000 + 500) * 1024},
		{"sessB", 300, []int{300, 301, 302, 303}, (1500 + 40000 + 2500 + 9000) * 1024},
		{"sessA2", 400, []int{400, 401}, (1200 + 60000) * 1024},
	} {
		tr, err := SessionTree(procs, tc.session)
		if err != nil {
			t.Fatalf("%s: %v", tc.session, err)
		}
		if tr.Session != tc.session || tr.Root != tc.root {
			t.Errorf("%s: session %q root %d, want root %d", tc.session, tr.Session, tr.Root, tc.root)
		}
		if got := pids(tr); !slices.Equal(got, tc.pids) {
			t.Errorf("%s: pids %v, want %v", tc.session, got, tc.pids)
		}
		if tr.Count != len(tc.pids) {
			t.Errorf("%s: count %d, want %d", tc.session, tr.Count, len(tc.pids))
		}
		if tr.RSS != tc.rss {
			t.Errorf("%s: rss %d, want %d", tc.session, tr.RSS, tc.rss)
		}
	}
}

func TestSessionTreeKeepsEachProcess(t *testing.T) {
	tr, err := SessionTree(Parse(table), "sessA")
	if err != nil {
		t.Fatal(err)
	}
	want := Process{PID: 202, PPID: 201, RSS: 3000 * 1024, Args: "/bin/zsh -c go test ./internal/x"}
	if tr.Processes[2] != want {
		t.Errorf("process %+v, want %+v", tr.Processes[2], want)
	}
}

func TestSessionTreePrefixIsNotRoot(t *testing.T) {
	// sess only prefixes sessA and sessA2: neither is its root.
	if _, err := SessionTree(Parse(table), "sess"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err %v, want ErrNotFound", err)
	}
	if IsSessionRoot("/home/u/bin/rhiza session run sessA2", "sessA") {
		t.Error("sessA2 matched sessA")
	}
}

func TestSessionTreeNotFound(t *testing.T) {
	tr, err := SessionTree(Parse(table), "sessZ")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err %v, want ErrNotFound", err)
	}
	if !strings.Contains(err.Error(), "sessZ") {
		t.Errorf("err %q does not name the session", err)
	}
	if tr.Count != 0 || tr.Processes != nil {
		t.Errorf("tree %+v, want none", tr)
	}
	if _, err := SessionTree(nil, "sessA"); !errors.Is(err, ErrNotFound) {
		t.Errorf("empty table: err %v, want ErrNotFound", err)
	}
}

func TestParseSkipsHeaderAndMalformed(t *testing.T) {
	procs := Parse(table)
	if len(procs) != 15 {
		t.Fatalf("parsed %d processes, want 15: %+v", len(procs), procs)
	}
	for _, p := range procs {
		if p.PID >= 700 || strings.Contains(p.Args, "COMMAND") {
			t.Errorf("kept %+v", p)
		}
	}
	if procs[0] != (Process{PID: 1, PPID: 0, RSS: 10836 * 1024, Args: "/usr/lib/systemd/systemd --switched-root"}) {
		t.Errorf("first %+v", procs[0])
	}
}

func TestReadPSFails(t *testing.T) {
	old := runPS
	t.Cleanup(func() { runPS = old })
	runPS = func() ([]byte, error) { return nil, errors.New("exit status 1") }
	if _, err := Read("sessA"); err == nil || !strings.HasPrefix(err.Error(), "ps: ") {
		t.Fatalf("err %v, want one naming ps", err)
	}
	if _, err := ReadTable(); err == nil || !strings.Contains(err.Error(), "ps") {
		t.Fatalf("err %v, want one naming ps", err)
	}
}

func TestReadUsesTable(t *testing.T) {
	old := runPS
	t.Cleanup(func() { runPS = old })
	runPS = func() ([]byte, error) { return []byte(table), nil }
	tr, err := Read("sessB")
	if err != nil || tr.Count != 4 {
		t.Fatalf("tree %+v, err %v", tr, err)
	}
}
