package release

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fixture struct {
	srv       *httptest.Server
	manifests map[string]*Manifest
	raw       map[string]string
	bodies    map[string][]byte
}

func newFixture(t *testing.T) (*fixture, *Resolver) {
	t.Helper()
	f := &fixture{manifests: map[string]*Manifest{}, raw: map[string]string{}, bodies: map[string][]byte{}}
	mux := http.NewServeMux()
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if body, ok := f.bodies[r.URL.Path]; ok {
			w.Write(body)
			return
		}
		ch := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/manifest.json")
		if s, ok := f.raw[ch]; ok {
			w.Write([]byte(s))
			return
		}
		if m, ok := f.manifests[ch]; ok {
			json.NewEncoder(w).Encode(m)
			return
		}
		http.NotFound(w, r)
	})
	for ch, v := range map[string]string{"release": "v0.1.0", "latest": "v0.1.1-12-gabc"} {
		m := &Manifest{Channel: ch, Version: v, Commit: "abc"}
		for _, p := range Platforms {
			path := fmt.Sprintf("/dl/%s/%s_%s.tar.gz", ch, p.OS, p.Arch)
			body := []byte("archive " + path)
			sum := sha256.Sum256(body)
			f.bodies[path] = body
			p.URL = f.srv.URL + path
			p.SHA256 = hex.EncodeToString(sum[:])
			m.Platforms = append(m.Platforms, p)
		}
		f.manifests[ch] = m
	}
	r := &Resolver{
		URLs:   map[string]string{"release": f.srv.URL + "/release/manifest.json", "latest": f.srv.URL + "/latest/manifest.json"},
		Client: f.srv.Client(),
	}
	return f, r
}

func TestNewestBothChannels(t *testing.T) {
	_, r := newFixture(t)
	for ch, want := range map[string]string{"release": "v0.1.0", "latest": "v0.1.1-12-gabc"} {
		got, err := r.Newest(context.Background(), ch)
		if err != nil || got != want {
			t.Errorf("%s: %q, %v; want %q", ch, got, err, want)
		}
	}
	if _, err := r.Newest(context.Background(), "nightly"); err == nil {
		t.Error("unknown channel accepted")
	}
}

func TestDownloadEachPlatform(t *testing.T) {
	_, r := newFixture(t)
	for _, p := range Platforms {
		r.OS, r.Arch = p.OS, p.Arch
		b, m, err := r.Download(context.Background(), "release")
		if err != nil {
			t.Fatalf("%s/%s: %v", p.OS, p.Arch, err)
		}
		if want := fmt.Sprintf("archive /dl/release/%s_%s.tar.gz", p.OS, p.Arch); string(b) != want || m.Version != "v0.1.0" {
			t.Errorf("%s/%s: %q %s", p.OS, p.Arch, b, m.Version)
		}
	}
}

func TestDownloadRunningPlatformDefault(t *testing.T) {
	_, r := newFixture(t)
	// The running platform may be unlisted (e.g. windows); only then is it an error.
	if _, _, err := r.Download(context.Background(), "latest"); err != nil && !strings.Contains(err.Error(), "no archive") {
		t.Fatal(err)
	}
}

func TestDownloadRejectsBadSHA(t *testing.T) {
	f, r := newFixture(t)
	r.OS, r.Arch = "linux", "amd64"
	f.bodies["/dl/release/linux_amd64.tar.gz"] = []byte("tampered")
	_, _, err := r.Download(context.Background(), "release")
	if err == nil || !strings.Contains(err.Error(), "sha256") {
		t.Fatalf("err = %v", err)
	}
}

func TestDownloadRejectsMissingPlatform(t *testing.T) {
	f, r := newFixture(t)
	r.OS, r.Arch = "linux", "riscv64"
	_, _, err := r.Download(context.Background(), "release")
	if err == nil || !strings.Contains(err.Error(), "no archive for linux/riscv64") {
		t.Fatalf("err = %v", err)
	}
	f.manifests["release"].Platforms = f.manifests["release"].Platforms[1:]
	r.Arch = "amd64"
	if _, _, err := r.Download(context.Background(), "release"); err == nil {
		t.Fatal("platform removed from the manifest was still picked")
	}
}

func TestRejectsMalformedManifest(t *testing.T) {
	f, r := newFixture(t)
	for name, raw := range map[string]string{
		"not json":      "<html>",
		"empty":         "{}",
		"wrong channel": `{"channel":"latest","version":"v1","commit":"c","platforms":[]}`,
		"short sha":     `{"channel":"release","version":"v1","commit":"c","platforms":[{"os":"linux","arch":"amd64","url":"http://x","sha256":"ab"}]}`,
		"missing url":   `{"channel":"release","version":"v1","commit":"c","platforms":[{"os":"linux","arch":"amd64","sha256":"` + strings.Repeat("a", 64) + `"}]}`,
	} {
		f.raw["release"] = raw
		if _, err := r.Newest(context.Background(), "release"); err == nil || !strings.Contains(err.Error(), "manifest") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
}

func TestManifestHTTPError(t *testing.T) {
	_, r := newFixture(t)
	r.URLs["release"] += ".missing"
	if _, err := r.Newest(context.Background(), "release"); err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("err = %v", err)
	}
}
