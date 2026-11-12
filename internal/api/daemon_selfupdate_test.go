package api

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/dphaener/rhiza/internal/release"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// A start pending at the end of a poll delays the self-update to the end of
// each later poll with none pending.
func TestDaemonLoopSelfUpdateWaitsForPending(t *testing.T) {
	common := t.TempDir()
	lock, err := daemon.TryLock(common, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	polls := 0
	var updatedAt []int
	// Polls 1 and 2 leave a start pending, polls 3 and 4 do not.
	pendingAt := map[int]int{1: 1, 2: 1}
	ctx := testCtx(t.TempDir())
	ticks, stop := make(chan time.Time), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		l, err := daemonLoop(ctx, common, lock,
			func() error { polls++; return nil },
			func() int { return pendingAt[polls] },
			nil, func() { updatedAt = append(updatedAt, polls) }, ticks, stop)
		if l != nil {
			l.Unlock()
		}
		done <- err
	}()
	for range 3 {
		ticks <- time.Now()
	}
	close(stop)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if polls != 4 || len(updatedAt) != 2 || updatedAt[0] != 3 || updatedAt[1] != 4 {
		t.Fatalf("polls %d, self-update at polls %v; want 4 polls, self-update at [3 4]", polls, updatedAt)
	}
}

// A root that still sets the removed self_update_project and
// rhiza.daemon.selfUpdate keeps working: its settings read autoUpdate =
// main, the daemon notes name rhiza.daemon.autoUpdate, and daemon set
// writes rhiza.toml without the old keys.
func TestDaemonSettingsMigrateRemovedSelfUpdate(t *testing.T) {
	ctx, _ := daemonRepo(t)
	path := filepath.Join(ctx.Root, "rhiza.toml")
	old := "self_update_project = \"repo\"\n\n[projects.repo]\npath = \".\"\n\n[daemon]\n\"rhiza.daemon.selfUpdate\" = \"true\"\n"
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v := settingValue(rows, "rhiza.daemon.autoUpdate"); v != "main" {
		t.Fatalf("autoUpdate = %q, want main", v)
	}
	_, reg, _ := rootOf(ctx)
	if notes := removedSettingNotes(ctx, reg); len(notes) != 2 || !strings.Contains(strings.Join(notes, "\n"), "rhiza.daemon.autoUpdate") {
		t.Fatalf("notes %q", notes)
	}
	resp, err := runDaemonSet(ctx, DaemonSetRequest{Key: "rhiza.daemon.maxRuns", Value: "2"})
	if err != nil || resp.(DaemonSetResponse).Scope != ScopeMachine {
		t.Fatalf("daemon set = %+v, %v", resp, err)
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "selfUpdate") || strings.Contains(string(data), "self_update_project") ||
		!strings.Contains(string(data), "rhiza.daemon.autoUpdate") {
		t.Fatalf("rhiza.toml after daemon set:\n%s", data)
	}
	_, err = runDaemonSet(ctx, DaemonSetRequest{Key: "rhiza.daemon.selfUpdate", Value: "true"})
	if err == nil || !strings.Contains(err.Error(), "rhiza.daemon.autoUpdate") {
		t.Fatalf("daemon set selfUpdate error = %v", err)
	}
}

// A machine setting under a root whose rhiza.toml does not load fails: it
// never falls back to the project's git config.
func TestDaemonSetBrokenRootWritesNoGitConfig(t *testing.T) {
	ctx, _ := daemonRepo(t)
	if err := os.WriteFile(filepath.Join(ctx.Root, "rhiza.toml"), []byte("not = = toml\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runDaemonSet(ctx, DaemonSetRequest{Key: "rhiza.daemon.autoUpdate", Value: "main"}); err == nil {
		t.Fatal("daemon set under a broken root succeeded")
	}
	if v, err := ctx.Git("config", "--get", "rhiza.daemon.autoUpdate"); err == nil && strings.TrimSpace(v) != "" {
		t.Fatalf("git config got %q", v)
	}
}

// hangProcess is a program whose output never ends until it is killed.
type hangProcess struct {
	r      *io.PipeReader
	w      *io.PipeWriter
	killed chan struct{}
	once   sync.Once
}

func (h *hangProcess) Output() io.Reader { return h.r }
func (h *hangProcess) Wait() (int, error) {
	<-h.killed
	return 137, nil
}
func (h *hangProcess) Kill() {
	h.once.Do(func() { close(h.killed); h.w.Close() })
}

// A self-update fetch that hangs is killed after its timeout, without a
// terminal prompt, and reported as failed.
func TestFetchMainWithinTimesOut(t *testing.T) {
	ctx := testCtx(t.TempDir())
	var got Program
	h := &hangProcess{killed: make(chan struct{})}
	h.r, h.w = io.Pipe()
	ctx.Start = func(p Program) (Process, error) { got = p; return h, nil }
	err := fetchMainWithin(ctx, "/main", 10*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("err %v", err)
	}
	if got.Name != "git" || strings.Join(got.Args, " ") != "fetch origin main" || got.Dir != "/main" ||
		strings.Join(got.Env, " ") != "GIT_TERMINAL_PROMPT=0" {
		t.Fatalf("program %+v", got)
	}
}

func TestFetchMainWithinExitCode(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.Start = (&shStart{codes: map[string]int{}}).start
	if err := fetchMainWithin(ctx, "/main", time.Second); err != nil {
		t.Fatalf("err %v", err)
	}
	ctx.Start = func(Program) (Process, error) { return &fakeProcess{out: strings.NewReader(""), code: 128}, nil }
	if err := fetchMainWithin(ctx, "/main", time.Second); err == nil {
		t.Fatal("exit 128 not an error")
	}
}

// autoUpdateRepo is daemonRepo whose git config sets rhiza.daemon.autoUpdate
// to value ("" for unset) and whose root registers one project at ".",
// holding a go.mod for module when module is not empty. Its origin/main has
// moved past it. It records every git and make run in that project.
func autoUpdateRepo(t *testing.T, value, module string) (*Ctx, func() []string) {
	t.Helper()
	ctx, common := daemonRepo(t)
	toml := "[projects.repo]\npath = \".\"\n"
	if err := os.WriteFile(filepath.Join(ctx.Root, "rhiza.toml"), []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}
	if module != "" {
		if err := os.WriteFile(filepath.Join(ctx.Root, "go.mod"), []byte("module "+module+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var mu sync.Mutex
	var ran []string
	ctx.Git = func(args ...string) (string, error) {
		switch {
		case len(args) == 2 && args[0] == "rev-parse" && args[1] == "--git-common-dir":
			return common + "\n", nil
		case len(args) == 3 && args[0] == "config" && args[2] == "rhiza.daemon.autoUpdate" && value != "":
			return value + "\n", nil
		}
		return "", errors.New("unexpected git " + strings.Join(args, " "))
	}
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if dir != ctx.Root {
			return nil, errors.New("unexpected " + name)
		}
		mu.Lock()
		ran = append(ran, name+" "+strings.Join(args, " "))
		mu.Unlock()
		if name == "make" {
			return nil, nil
		}
		switch args[0] {
		case "rev-parse":
			if args[2] == "HEAD" {
				return []byte("aaa\n"), nil
			}
			return []byte("bbb\n"), nil
		case "merge-base":
			return nil, errors.New("exit status 1")
		case "symbolic-ref":
			return []byte("main\n"), nil
		}
		return nil, nil
	}
	return ctx, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), ran...)
	}
}

// With rhiza.daemon.autoUpdate off or unset the daemon neither fetches nor
// runs make update; main with a Rhiza checkout fetches, fast-forwards and
// runs make update once over several polls; main with none registered logs
// it and updates nothing.
func TestDaemonAutoUpdateSetting(t *testing.T) {
	start := daemon.Image{Size: 10, Mod: time.Unix(1700000000, 0)}
	old := daemonReexec
	t.Cleanup(func() { daemonReexec = old })
	daemonReexec = func() *daemon.Reexec { return (&reexecFake{img: start}).reexec(start) }

	for _, value := range []string{"", "off"} {
		ctx, ran := autoUpdateRepo(t, value, daemon.RhizaModule)
		if err := runDaemonTicks(t, ctx, 2, nil); err != nil {
			t.Fatal(err)
		}
		if got := ran(); len(got) != 0 {
			t.Fatalf("autoUpdate %q: ran %q, want nothing", value, got)
		}
	}

	ctx, ran := autoUpdateRepo(t, "main", "example.com/other")
	var stderr strings.Builder
	ctx.Stderr = &stderr
	if err := runDaemonTicks(t, ctx, 2, nil); err != nil {
		t.Fatal(err)
	}
	if got := ran(); len(got) != 0 {
		t.Fatalf("no Rhiza checkout: ran %q, want nothing", got)
	}
	if !strings.Contains(stderr.String(), "no registered project is Rhiza's source checkout") {
		t.Fatalf("no Rhiza checkout: log %q", stderr.String())
	}

	ctx, ran = autoUpdateRepo(t, "main", daemon.RhizaModule)
	if err := runDaemonTicks(t, ctx, 2, nil); err != nil {
		t.Fatal(err)
	}
	got := strings.Join(ran(), "\n")
	for _, want := range []string{"git fetch origin main", "git merge --ff-only bbb", "make update"} {
		if !strings.Contains(got, want) {
			t.Fatalf("autoUpdate main: ran\n%s\nwant %q", got, want)
		}
	}
	if n := strings.Count(got, "make update"); n != 1 {
		t.Fatalf("autoUpdate main: make update ran %d times, want 1", n)
	}
}

// With rhiza.daemon.autoUpdate = latest the daemon reads the latest channel's
// feed and installs a newer build over its executable, running no git fetch
// and no make; an up-to-date daemon leaves the executable alone.
func TestDaemonAutoUpdateLatest(t *testing.T) {
	start := daemon.Image{Size: 10, Mod: time.Unix(1700000000, 0)}
	old := daemonReexec
	t.Cleanup(func() { daemonReexec = old })
	daemonReexec = func() *daemon.Reexec { return (&reexecFake{img: start}).reexec(start) }

	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	tw.WriteHeader(&tar.Header{Name: "rhiza", Typeflag: tar.TypeReg, Mode: 0o755, Size: 3})
	tw.Write([]byte("new"))
	tw.Close()
	gz.Close()
	sum := sha256.Sum256(archive.Bytes())

	var version string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/manifest.json" {
			json.NewEncoder(w).Encode(release.Manifest{Channel: release.ChannelLatest, Version: version, Commit: "abc",
				Platforms: []release.Platform{{OS: runtime.GOOS, Arch: runtime.GOARCH, URL: "http://" + r.Host + "/a.tar.gz", SHA256: hex.EncodeToString(sum[:])}}})
			return
		}
		w.Write(archive.Bytes())
	}))
	t.Cleanup(srv.Close)

	oldURLs, oldExe := channelUpdateURLs, channelUpdateExe
	t.Cleanup(func() { channelUpdateURLs, channelUpdateExe = oldURLs, oldExe })
	channelUpdateURLs = map[string]string{release.ChannelLatest: srv.URL + "/manifest.json"}

	for _, c := range []struct{ current, newest, want string }{
		{"v0.1.1-main.5+aaa1111", "v0.1.1-main.7+bbb2222", "new"},
		{"v0.1.1-main.7+bbb2222", "v0.1.1-main.7+bbb2222", "old"},
	} {
		exe := filepath.Join(t.TempDir(), "rhiza")
		if err := os.WriteFile(exe, []byte("old"), 0o755); err != nil {
			t.Fatal(err)
		}
		channelUpdateExe = func() (string, error) { return exe, nil }
		version = c.newest
		ctx, ran := autoUpdateRepo(t, "latest", daemon.RhizaModule)
		ctx.Version = c.current
		if err := runDaemonTicks(t, ctx, 2, nil); err != nil {
			t.Fatal(err)
		}
		if got, _ := os.ReadFile(exe); string(got) != c.want {
			t.Fatalf("running %s, feed %s: exe holds %q, want %q", c.current, c.newest, got, c.want)
		}
		if got := ran(); len(got) != 0 {
			t.Fatalf("latest ran %q, want no git and no make", got)
		}
	}
}
