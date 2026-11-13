package api

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// Production-sized fixture of the idle benchmarks: a log of more than a
// thousand change sets, a few runs in flight, and the Claude transcripts the
// usage governor walks.
const (
	idleChangeSets = 1100
	idleRuns       = 4
	idleTranscript = 4000
)

// idleFix is the fixture of BenchmarkIdleRootPoll and BenchmarkIdleDriverPoll.
// Everything lives in temp dirs; nothing touches the network, Rex, an agent
// or gh.
type idleFix struct {
	home, other *daemonFoldFix
	projs       []*daemonProj
	r           DaemonRequest
	sessions    *session.Store
	step        *stepRunner
	sid         string
}

// writeTranscripts lays n Claude transcript files under config/projects, a
// few assistant lines each, the way the usage reader walks them.
func writeTranscripts(tb testing.TB, config string, n int, at time.Time) {
	tb.Helper()
	line := func(i int, ts time.Time) string {
		return fmt.Sprintf(`{"type":"assistant","timestamp":%q,"requestId":"req-%d","message":{"id":"msg-%d","model":"claude-sonnet-5-5","usage":{"input_tokens":120,"output_tokens":340,"cache_read_input_tokens":9000,"cache_creation_input_tokens":200}}}`+"\n", ts.UTC().Format(time.RFC3339), i, i)
	}
	const perDir = 100
	for i := range n {
		dir := filepath.Join(config, "projects", fmt.Sprintf("-home-rhiza-code-p%d", i/perDir))
		if i%perDir == 0 {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				tb.Fatal(err)
			}
		}
		ts := at.Add(-time.Duration(i) * time.Minute)
		path := filepath.Join(dir, fmt.Sprintf("s%05d.jsonl", i))
		if err := os.WriteFile(path, []byte(line(2*i, ts)+line(2*i+1, ts.Add(time.Second))), 0o644); err != nil {
			tb.Fatal(err)
		}
		if err := os.Chtimes(path, ts, ts); err != nil {
			tb.Fatal(err)
		}
	}
}

var ulidRe = regexp.MustCompile(`[0-9A-HJKMNP-TV-Z]{26}`)

const ulidAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// growLog returns a grow function (see newDaemonFoldFixN) that copies the
// landed "done-" change sets of the log until it holds total change sets.
// A write through the operations reads the whole log each time, so writing a
// thousand that way takes minutes. A copy keeps its entries' time prefix
// (so their order holds) and gets new IDs and handles; the actors stay.
func growLog(tb testing.TB, total int) func(string) {
	return func(logDir string) {
		tb.Helper()
		dirs, err := os.ReadDir(logDir)
		if err != nil {
			tb.Fatal(err)
		}
		type src struct {
			name  string
			files map[string]string
		}
		var srcs []src
		for _, d := range dirs {
			files := map[string]string{}
			es, _ := os.ReadDir(filepath.Join(logDir, d.Name()))
			for _, e := range es {
				data, err := os.ReadFile(filepath.Join(logDir, d.Name(), e.Name()))
				if err != nil {
					tb.Fatal(err)
				}
				files[e.Name()] = string(data)
			}
			all := ""
			for _, c := range files {
				all += c
			}
			if strings.Contains(all, `"handle": "done-`) {
				srcs = append(srcs, src{d.Name(), files})
			}
		}
		if len(srcs) == 0 {
			tb.Fatal("no change set to copy")
		}
		actorRe := regexp.MustCompile(`"(?:actor|for)": "([0-9A-Z]{26})"`)
		for n := len(dirs); n < total; n++ {
			s := srcs[n%len(srcs)]
			ids := map[string]string{}
			keep := map[string]bool{}
			for _, c := range s.files {
				for _, m := range actorRe.FindAllStringSubmatch(c, -1) {
					keep[m[1]] = true
				}
			}
			suffix := func(k int) string {
				b := []byte("00000")
				for i := 4; i >= 0; i-- {
					b[i] = ulidAlphabet[k%32]
					k /= 32
				}
				return string(b)
			}
			next := 0
			rename := func(c string) string {
				c = ulidRe.ReplaceAllStringFunc(c, func(u string) string {
					if keep[u] {
						return u
					}
					if v, ok := ids[u]; ok {
						return v
					}
					next++
					v := u[:10] + suffix(n) + u[15:25] + ulidAlphabet[next%32:next%32+1]
					ids[u] = v
					return v
				})
				return strings.ReplaceAll(c, `"handle": "done-`, fmt.Sprintf(`"handle": "c%d-done-`, n))
			}
			// The same renaming for every file: IDs first seen in one file
			// show up in the others.
			var names []string
			for name := range s.files {
				names = append(names, name)
			}
			sortStrings(names)
			out := map[string]string{}
			for _, name := range names {
				out[rename(name)] = rename(s.files[name])
			}
			dir := filepath.Join(logDir, rename(s.name))
			if err := os.MkdirAll(dir, 0o755); err != nil {
				tb.Fatal(err)
			}
			for name, c := range out {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(c), 0o644); err != nil {
					tb.Fatal(err)
				}
			}
		}
	}
}

func sortStrings(s []string) { slices.Sort(s) }

func newIdleFix(b testing.TB) *idleFix {
	b.Helper()
	f := &idleFix{home: newDaemonFoldFixN(b, 98, growLog(b, idleChangeSets)), other: newDaemonFoldFixN(b, 50, nil)}
	config := b.TempDir()
	writeTranscripts(b, config, idleTranscript, f.home.ctx.Now())
	env := func(k string) string {
		switch k {
		case "CLAUDE_CONFIG_DIR":
			return config
		case "RHIZA_TEST_NOW":
			return "1"
		case "RHIZA_TEST_FREE_BYTES":
			return "100000000000"
		}
		return ""
	}
	for i, fx := range []*daemonFoldFix{f.home, f.other} {
		fx.ctx.Getenv = env
		name := []string{"home", "other"}[i]
		p, err := newDaemonProj(fx.ctx, name, fx.ctx.Root, daemon.Holder{}, i == 0, f.r, nil)
		if err != nil {
			b.Fatal(err)
		}
		b.Cleanup(p.stop)
		// What the project could start was started and is held: an idle
		// poll starts nothing.
		p.unjournaled[fx.target] = daemonPending{at: fx.ctx.Now()}
		f.projs = append(f.projs, p)
	}
	// Runs in flight, each in the step that waits on an agent session.
	home := f.projs[0]
	src := "[step.s]\nkind = \"agent\"\nprompt = \"work\"\n"
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		b.Fatal(err)
	}
	for i := range idleRuns {
		cs := f.home.target
		if err := run.Start(home.rs, fmt.Sprintf("R%d", i), run.Started{Workflow: "wf", Definition: src, Change: cs}, f.home.ctx.Now(), def); err != nil {
			b.Fatal(err)
		}
	}
	f.sessions = session.NewStore(session.Dir(home.common))
	f.sid = "S-idle"
	if err := f.sessions.Write(session.Record{ID: f.sid, Change: f.home.target, Role: "builder"}); err != nil {
		b.Fatal(err)
	}
	f.step = &stepRunner{ctx: f.home.ctx, change: f.home.target, run: "R0", worktree: f.home.ctx.Root, sessions: f.sessions}
	return f
}

// rootPoll is one poll of the root daemon: both registered projects'
// settings, usage and demand read, the slots shared round-robin, then each
// project's own turn (daemonPoll1, notices, watchdog or nudge, events, state).
func (f *idleFix) rootPoll(b testing.TB) {
	b.Helper()
	turns := make([]daemonTurn, len(f.projs))
	var picks []daemon.ProjectPick
	running := 0
	for i, p := range f.projs {
		turns[i] = p.prepare(f.r, nil)
		if turns[i].err != nil {
			b.Fatal(turns[i].err)
		}
		picks = append(picks, turns[i].pick)
		running += turns[i].pick.Running
	}
	picked := daemon.RoundRobin(picks, running, daemon.DefaultMaxRuns, "")
	share := map[string]int{}
	for _, pk := range picked {
		share[pk.Project]++
	}
	var resp DaemonResponse
	for i, p := range f.projs {
		if err := p.poll(f.r, &resp, turns[i], daemon.DefaultMaxRuns, share[p.name]); err != nil {
			b.Fatal(err)
		}
	}
}

// driverPoll is one idle tick of a run driver waiting on an agent step: the
// session's status, the runs in flight and what each is blocked on.
func (f *idleFix) driverPoll(b testing.TB) {
	b.Helper()
	st, err := f.step.AgentStatus(f.sid)
	if err != nil || st.Done {
		b.Fatalf("agent status %+v, %v", st, err)
	}
	runs, err := f.projs[0].rs.List()
	if err != nil || len(runs) != idleRuns {
		b.Fatalf("%d runs, %v", len(runs), err)
	}
	for _, j := range runs {
		runBlockedNow(j.State(), f.sessions)
	}
}

// BenchmarkIdleRootPoll times one poll of the root daemon with two
// registered projects and nothing changed. Profile it with
//
//	go test -trimpath ./internal/api/ -run '^$' -bench IdleRootPoll -cpuprofile $TMPDIR/idle.pprof
func BenchmarkIdleRootPoll(b *testing.B) {
	f := newIdleFix(b)
	f.rootPoll(b) // the first poll folds and fills the caches
	b.ResetTimer()
	for range b.N {
		f.rootPoll(b)
	}
}

// BenchmarkIdleDriverPoll times one idle poll of a run driver over the same
// fixture.
func BenchmarkIdleDriverPoll(b *testing.B) {
	f := newIdleFix(b)
	f.driverPoll(b)
	b.ResetTimer()
	for range b.N {
		f.driverPoll(b)
	}
}
