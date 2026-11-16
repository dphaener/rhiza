package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
)

// writeFixture writes a fixture script named name.txtar into a fresh
// directory and returns its path.
func writeFixture(t *testing.T, name, script string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), name+".txtar")
	if err := os.WriteFile(f, []byte(script), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

// callRhiza is a fake bin/ command that starts rhiza as a process and
// prints its exit status.
const callRhiza = `
-- bin/callrhiza --
#!/bin/sh
rhiza "$@"
echo "exit=$?"
`

// TestProcessScripts shows the process list enforced at the rhiza process
// entry: a script off processScripts whose fake bin/ command starts rhiza
// gets exit 2 and a message naming the list, and fails at its next bare
// rhiza line (or at its end) with that message, while its bare rhiza lines
// run in-process; the same script on the list starts rhiza as a process.
func TestProcessScripts(t *testing.T) {
	if _, ok := processScripts["process-list-fixture"]; !ok {
		t.Fatal("process-list-fixture is not on processScripts")
	}
	if _, ok := processScripts["off-list-fixture"]; ok {
		t.Fatal("off-list-fixture is on processScripts")
	}
	const listMsg = "rhiza: script off-list-fixture started rhiza as a process but is not on processScripts, the process list in cmd/rhiza/main_test.go"
	// Its bare rhiza line runs in-process, so the list does not apply.
	inProcess := writeFixture(t, "off-list-fixture", `
rhiza version
stdout '^rhiza '
! stderr .
`)
	on := writeFixture(t, "process-list-fixture", `
env PATH=$WORK/bin:$PATH
chmod 755 bin/callrhiza
exec callrhiza version
stdout '^rhiza '
stdout '^exit=0$'
! stderr .
exec rhiza version
stdout '^rhiza '
`+callRhiza)
	testscript.Run(t, scriptParams(t, []string{inProcess, on}))

	for _, c := range []struct{ name, script string }{
		// The refused callback's caller passes on; the next bare rhiza
		// line fails with the list's message.
		{"next rhiza line", `
env PATH=$WORK/bin:$PATH
chmod 755 bin/callrhiza
exec callrhiza version
stdout '^exit=2$'
! stdout '^rhiza '
stderr '^` + regexp.QuoteMeta(listMsg) + `'
rhiza version
`},
		// No later rhiza line: the script fails at its end.
		{"script end", `
env PATH=$WORK/bin:$PATH
chmod 755 bin/callrhiza
exec callrhiza version
stdout '^exit=2$'
`},
	} {
		t.Run(c.name, func(t *testing.T) {
			rt := &recordT{}
			testscript.RunT(rt, scriptParams(t, []string{writeFixture(t, "off-list-fixture", c.script+callRhiza)}))
			if !rt.failed || !strings.Contains(rt.log.String(), listMsg) {
				t.Errorf("failed %v, want a failure naming the list:\n%s", rt.failed, rt.log.String())
			}
		})
	}
}

// TestProcessRefusedFailsFast shows an off-list script whose in-process
// rhiza line waits on a fake bin/ callback that rhiza refused failing in
// seconds with the list's message, not at the go test timeout: the hook
// command calls rhiza, which exits 2, then waits a minute under a
// five-minute hook timeout.
func TestProcessRefusedFailsFast(t *testing.T) {
	f := writeFixture(t, "off-list-fixture", `
env PATH=$WORK/bin:$PATH
env RHIZA_TEST_NOTIFY_TIMEOUT=5m
chmod 755 bin/callrhiza
rhiza hook run notify --command 'callrhiza version; sleep 60'
`+callRhiza)
	rt := &recordT{}
	start := time.Now()
	testscript.RunT(rt, scriptParams(t, []string{f}))
	took := time.Since(start)
	if !rt.failed || !strings.Contains(rt.log.String(), "started rhiza as a process but is not on processScripts") {
		t.Errorf("failed %v, want a failure naming the list:\n%s", rt.failed, rt.log.String())
	}
	if took > 20*time.Second {
		t.Errorf("the script took %s to fail, want seconds", took)
	}
}

// TestScriptEnv shows the environment the in-process rhiza gets: setup's
// names, a value an env line set, NAME= for a name an env line cleared (and
// for one a later env line sets), and a git child of run() reading the
// script's HOME and GIT_CONFIG_GLOBAL.
func TestScriptEnv(t *testing.T) {
	f := writeFixture(t, "script-env-fixture", `
env FOO=x
env RHIZA_TEST_SEED=
script-env
stdout '^WORK=.*script-script-env-fixture$'
stdout '^HOME=.*script-env-fixture/home$'
stdout '^GIT_CONFIG_GLOBAL=.*script-env-fixture/home/\.gitconfig$'
stdout '^RHIZA_TEST_NOW=2026-10-24T12:00:00Z$'
stdout '^RHIZA_TEST_SCRIPT=script-env-fixture$'
stdout '^FOO=x$'
stdout '^RHIZA_TEST_SEED=$'
stdout '^LATER=$'
env LATER=y
script-env
stdout '^LATER=y$'

# run()'s git child reads the script's identity: the first write records it.
exec git init -q
rhiza init
rhiza change new 'Try it'
rhiza change ls
stdout ' landed record: Record actor test@rhiza.invalid$'
`)
	p := scriptParams(t, []string{f})
	p.Cmds["script-env"] = func(ts *testscript.TestScript, neg bool, args []string) {
		for _, kv := range scriptEnv(ts) {
			fmt.Fprintln(ts.Stdout(), kv)
		}
	}
	testscript.Run(t, p)
	if got := envLineNames("env A=1 B= 'C=x y'\n[unix] env D=2\n! env E\n# env F=1\nexec env G=1\n"); strings.Join(got, ",") != "A,B,C,D" {
		t.Errorf("envLineNames = %q, want A,B,C,D", got)
	}
}

// TestInProcessCommand shows a bare rhiza line running run() with no
// process started (the rhiza process entry logs none) and reporting the
// same stdout, stderr and exit status as exec rhiza.
func TestInProcessCommand(t *testing.T) {
	f := writeFixture(t, "in-process-fixture", `
env RHIZA_TEST_PROCESS_LOG=$WORK/procs
rhiza version
cp stdout in.stdout
cp stderr in.stderr
! rhiza nosuch
cp stdout in-fail.stdout
cp stderr in-fail.stderr
stderr .
! exists procs

exec rhiza version
cmp stdout in.stdout
cmp stderr in.stderr
! exec rhiza nosuch
cmp stdout in-fail.stdout
cmp stderr in-fail.stderr
exists procs
`)
	testscript.Run(t, scriptParams(t, []string{f}))

	// The exit status, from the log of a failing line.
	status := regexp.MustCompile(`\[exit status (\d+)\]`)
	for _, line := range []string{"rhiza nosuch", "exec rhiza nosuch"} {
		rt := &recordT{}
		testscript.RunT(rt, scriptParams(t, []string{writeFixture(t, "in-process-fixture", line+"\n")}))
		if !rt.failed {
			t.Errorf("%s: passed, want a failure", line)
		}
		m := status.FindStringSubmatch(rt.log.String())
		if m == nil {
			t.Errorf("%s: no exit status in the log:\n%s", line, rt.log.String())
			continue
		}
		if m[1] != "2" {
			t.Errorf("%s: exit status %s, want 2", line, m[1])
		}
	}
}

// recordT is a testscript.T that keeps the log and whether the script
// failed, so a test can read a failing script's exit status.
type recordT struct {
	mu     sync.Mutex
	log    strings.Builder
	failed bool
}

func (r *recordT) Skip(a ...any)  { r.Log(a...); runtime.Goexit() }
func (r *recordT) Fatal(a ...any) { r.Log(a...); r.FailNow() }
func (r *recordT) Parallel()      {}
func (r *recordT) Verbose() bool  { return true }

func (r *recordT) Log(a ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Fprintln(&r.log, a...)
}

func (r *recordT) FailNow() {
	r.mu.Lock()
	r.failed = true
	r.mu.Unlock()
	runtime.Goexit()
}

func (r *recordT) Run(name string, f func(testscript.T)) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f(r)
	}()
	<-done
}

// TestProcessListConverted keeps every processScripts entry converted: an
// entry whose reason is notConverted fails, naming the script, as the
// fixture list shows; the real list has none.
func TestProcessListConverted(t *testing.T) {
	fixture := map[string]string{
		"stays-a-process": "a stdin feed",
		"left-behind":     notConverted,
	}
	if got := unconverted(fixture); strings.Join(got, ",") != "left-behind" {
		t.Errorf("unconverted(fixture) = %q, want [left-behind]", got)
	}
	for _, name := range unconverted(processScripts) {
		t.Errorf("processScripts: %s keeps the reason %q: run its rhiza lines in-process, or give the reason it starts rhiza as a process", name, notConverted)
	}
}
