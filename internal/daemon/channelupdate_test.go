package daemon

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/release"
)

// fakeFeed serves a latest channel's manifest and archive.
type fakeFeed struct {
	srv      *httptest.Server
	version  string
	binary   string
	badSum   bool
	manifest int
	archives int
}

func newFeed(t *testing.T, version, binary string) *fakeFeed {
	f := &fakeFeed{version: version, binary: binary}
	mux := http.NewServeMux()
	mux.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		f.manifest++
		a := f.archive()
		sum := sha256.Sum256(a)
		hexSum := hex.EncodeToString(sum[:])
		if f.badSum {
			hexSum = strings.Repeat("0", 64)
		}
		m := release.Manifest{Channel: release.ChannelLatest, Version: f.version, Commit: "abc", Platforms: []release.Platform{
			{OS: runtime.GOOS, Arch: runtime.GOARCH, URL: f.srv.URL + "/rhiza.tar.gz", SHA256: hexSum},
		}}
		json.NewEncoder(w).Encode(m)
	})
	mux.HandleFunc("/rhiza.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		f.archives++
		w.Write(f.archive())
	})
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeFeed) archive() []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	tw.WriteHeader(&tar.Header{Name: "rhiza_x/rhiza", Typeflag: tar.TypeReg, Mode: 0o755, Size: int64(len(f.binary))})
	tw.Write([]byte(f.binary))
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

func (f *fakeFeed) update(t *testing.T, current string) (*ChannelUpdate, string, *[]string) {
	exe := filepath.Join(t.TempDir(), "rhiza")
	if err := os.WriteFile(exe, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}
	var logs []string
	return &ChannelUpdate{
		Channel:  release.ChannelLatest,
		Resolver: &release.Resolver{URLs: map[string]string{release.ChannelLatest: f.srv.URL + "/manifest.json"}},
		Current:  current,
		Exe:      exe,
		Log:      func(format string, a ...any) { logs = append(logs, format) },
	}, exe, &logs
}

func TestChannelUpdateInstallsNewerBuild(t *testing.T) {
	f := newFeed(t, "v0.1.1-main.7+abc1234", "new binary")
	u, exe, _ := f.update(t, "v0.1.1-main.5+def5678")
	before, _ := StatImage(exe)
	if !u.Run() {
		t.Fatal("newer build not installed")
	}
	got, _ := os.ReadFile(exe)
	if string(got) != "new binary" {
		t.Fatalf("exe holds %q", got)
	}
	if fi, _ := os.Stat(exe); fi.Mode().Perm()&0o100 == 0 {
		t.Fatalf("mode %v", fi.Mode())
	}
	if after, _ := StatImage(exe); after == before {
		t.Fatal("the image did not change, so the re-exec would not fire")
	}
	// Installed: the next check, of the same build, installs nothing.
	u.Current = f.version
	if u.Run() || f.archives != 1 {
		t.Fatalf("second run installed again (archives %d)", f.archives)
	}
}

func TestChannelUpdateLeavesUpToDateInstallAlone(t *testing.T) {
	f := newFeed(t, "v0.1.1-main.5+abc1234", "new binary")
	u, exe, _ := f.update(t, "v0.1.1-main.5+abc1234")
	if u.Run() {
		t.Fatal("an up-to-date install was replaced")
	}
	// An older published build is not a downgrade.
	u.Current = "v0.1.1-main.9+abc1234"
	if u.Run() {
		t.Fatal("a newer install was replaced by an older build")
	}
	got, _ := os.ReadFile(exe)
	if string(got) != "old" || f.archives != 0 {
		t.Fatalf("exe %q, archives downloaded %d", got, f.archives)
	}
}

func TestChannelUpdateBadChecksumInstallsNothingAndIsNotRetried(t *testing.T) {
	f := newFeed(t, "v0.1.1-main.7+abc1234", "new binary")
	f.badSum = true
	u, exe, logs := f.update(t, "v0.1.1-main.5+def5678")
	if u.Run() {
		t.Fatal("installed despite a checksum mismatch")
	}
	if got, _ := os.ReadFile(exe); string(got) != "old" {
		t.Fatalf("exe holds %q", got)
	}
	if len(*logs) == 0 {
		t.Fatal("the failure was not logged")
	}
	u.Run()
	if f.archives != 1 {
		t.Fatalf("the failed version was downloaded %d times", f.archives)
	}
}

func TestChannelUpdateUnreachableFeed(t *testing.T) {
	f := newFeed(t, "v0.1.1-main.7+abc1234", "new binary")
	u, exe, logs := f.update(t, "v0.1.1-main.5+def5678")
	f.srv.Close()
	if u.Run() || len(*logs) != 1 {
		t.Fatalf("logs %v", *logs)
	}
	if got, _ := os.ReadFile(exe); string(got) != "old" {
		t.Fatalf("exe holds %q", got)
	}
}
