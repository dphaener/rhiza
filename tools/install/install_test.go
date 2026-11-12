package install

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/release"
)

const script = "../../install.sh"

func TestInstallScript(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("install.sh supports linux and darwin")
	}
	for _, tool := range []string{"sh", "curl", "tar"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not installed", tool)
		}
	}

	t.Run("syntax", func(t *testing.T) {
		if out, err := exec.Command("sh", "-n", script).CombinedOutput(); err != nil {
			t.Fatalf("sh -n: %v\n%s", err, out)
		}
	})

	fx := newFixture(t)

	t.Run("release channel", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "bin")
		out, err := runScript(t, fx, dir)
		if err != nil {
			t.Fatalf("install: %v\n%s", err, out)
		}
		got, err := exec.Command(filepath.Join(dir, "rhiza"), "version").Output()
		if err != nil {
			t.Fatalf("run installed rhiza: %v", err)
		}
		if strings.TrimSpace(string(got)) != "fixture release" {
			t.Fatalf("installed rhiza printed %q, want the release build", got)
		}
	})

	t.Run("latest channel", func(t *testing.T) {
		dir := t.TempDir()
		out, err := runScript(t, fx, dir, "--channel", "latest")
		if err != nil {
			t.Fatalf("install: %v\n%s", err, out)
		}
		got, err := exec.Command(filepath.Join(dir, "rhiza"), "version").Output()
		if err != nil {
			t.Fatalf("run installed rhiza: %v", err)
		}
		if strings.TrimSpace(string(got)) != "fixture latest" {
			t.Fatalf("installed rhiza printed %q, want the latest build", got)
		}
	})

	t.Run("tampered archive", func(t *testing.T) {
		fx.tamper = true
		defer func() { fx.tamper = false }()
		dir := filepath.Join(t.TempDir(), "bin")
		out, err := runScript(t, fx, dir)
		if err == nil {
			t.Fatalf("install succeeded on a tampered archive:\n%s", out)
		}
		if !strings.Contains(out, "sha256 mismatch") {
			t.Fatalf("want a sha256 mismatch error, got:\n%s", out)
		}
		if _, err := os.Stat(filepath.Join(dir, "rhiza")); err == nil {
			t.Fatal("rhiza was installed despite the mismatch")
		}
	})

	// The documented form: the script on stdin, so $0 is the shell, with a
	// manifest on one line, so the split does not lean on tools/release's
	// indentation.
	t.Run("piped with a compact manifest", func(t *testing.T) {
		fx.compact = true
		defer func() { fx.compact = false }()
		dir := t.TempDir()
		out, err := pipeScript(t, fx, "--dir", dir)
		if err != nil {
			t.Fatalf("install: %v\n%s", err, out)
		}
		got, err := exec.Command(filepath.Join(dir, "rhiza"), "version").Output()
		if err != nil {
			t.Fatalf("run installed rhiza: %v", err)
		}
		if strings.TrimSpace(string(got)) != "fixture release" {
			t.Fatalf("installed rhiza printed %q, want the release build", got)
		}
	})

	t.Run("piped help", func(t *testing.T) {
		out, err := pipeScript(t, fx, "--help")
		if err != nil {
			t.Fatalf("--help: %v\n%s", err, out)
		}
		if !strings.Contains(out, "--channel release|latest") || strings.Contains(out, "sed:") {
			t.Fatalf("want the help text, got:\n%s", out)
		}
	})

	t.Run("bad channel", func(t *testing.T) {
		out, err := runScript(t, fx, t.TempDir(), "--channel", "nope")
		if err == nil || !strings.Contains(out, "unknown channel") {
			t.Fatalf("want an unknown channel error, got %v\n%s", err, out)
		}
	})
}

type fixture struct {
	srv     *httptest.Server
	tamper  bool
	compact bool // serve manifests as one-line JSON
}

// newFixture serves a channel: manifests for both channels and, per channel,
// an archive holding a shell script named rhiza that prints its channel.
func newFixture(t *testing.T) *fixture {
	t.Helper()
	fx := &fixture{}
	archives := map[string][]byte{}
	mux := http.NewServeMux()
	fx.srv = httptest.NewServer(mux)
	t.Cleanup(fx.srv.Close)

	for _, ch := range []string{release.ChannelRelease, release.ChannelLatest} {
		a := archive(t, "#!/bin/sh\necho fixture "+ch+"\n")
		archives[ch] = a
		sum := sha256.Sum256(a)
		m := release.Manifest{Channel: ch, Version: "v0.0.0-" + ch, Commit: "abc"}
		for _, p := range release.Platforms {
			p.URL = fmt.Sprintf("%s/archive/%s/%s_%s.tar.gz", fx.srv.URL, ch, p.OS, p.Arch)
			p.SHA256 = hex.EncodeToString(sum[:])
			m.Platforms = append(m.Platforms, p)
		}
		b, err := m.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		path := "/releases/latest/download/manifest.json"
		if ch == release.ChannelLatest {
			path = "/releases/download/latest/manifest.json"
		}
		var compact bytes.Buffer
		if err := json.Compact(&compact, b); err != nil {
			t.Fatal(err)
		}
		mux.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
			if fx.compact {
				w.Write(compact.Bytes())
				return
			}
			w.Write(b)
		})
		ch := ch
		mux.HandleFunc("/archive/"+ch+"/", func(w http.ResponseWriter, _ *http.Request) {
			a := archives[ch]
			if fx.tamper {
				a = append(append([]byte{}, a...), 0)
			}
			w.Write(a)
		})
	}
	return fx
}

func archive(t *testing.T, body string) []byte {
	t.Helper()
	var buf strings.Builder
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "rhiza", Mode: 0o755, Size: int64(len(body))}); err != nil {
		t.Fatal(err)
	}
	tw.Write([]byte(body))
	tw.Close()
	gz.Close()
	return []byte(buf.String())
}

func runScript(t *testing.T, fx *fixture, dir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command("sh", append([]string{script, "--dir", dir}, args...)...)
	cmd.Env = append(os.Environ(), "RHIZA_INSTALL_BASE_URL="+fx.srv.URL, "HOME="+t.TempDir())
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// pipeScript runs the script as curl | sh -s -- args would: on stdin.
func pipeScript(t *testing.T, fx *fixture, args ...string) (string, error) {
	t.Helper()
	f, err := os.Open(script)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cmd := exec.Command("sh", append([]string{"-s", "--"}, args...)...)
	cmd.Stdin = f
	cmd.Env = append(os.Environ(), "RHIZA_INSTALL_BASE_URL="+fx.srv.URL, "HOME="+t.TempDir())
	out, err := cmd.CombinedOutput()
	return string(out), err
}
