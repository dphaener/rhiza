package release

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDist(t *testing.T, skip string) string {
	t.Helper()
	dir := t.TempDir()
	var sums strings.Builder
	for _, p := range Platforms {
		name := fmt.Sprintf("rhiza_0.1.0_%s_%s.tar.gz", p.OS, p.Arch)
		body := []byte("archive " + name)
		if err := os.WriteFile(filepath.Join(dir, name), body, 0o644); err != nil {
			t.Fatal(err)
		}
		if name == skip {
			continue
		}
		s := sha256.Sum256(body)
		fmt.Fprintf(&sums, "%s  %s\n", hex.EncodeToString(s[:]), name)
	}
	if err := os.WriteFile(filepath.Join(dir, ChecksumsFile), []byte(sums.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func opts(dir string) Options {
	return Options{Dist: dir, Project: "rhiza", Tag: "v0.1.0", Commit: "abc123", Repo: "dphaener/rhiza"}
}

func TestBuild(t *testing.T) {
	m, err := Build(opts(writeDist(t, "")))
	if err != nil {
		t.Fatal(err)
	}
	if m.Channel != "release" || m.Version != "v0.1.0" || m.Commit != "abc123" || len(m.Platforms) != 4 {
		t.Fatalf("manifest: %+v", m)
	}
	for _, p := range m.Platforms {
		name := fmt.Sprintf("rhiza_0.1.0_%s_%s.tar.gz", p.OS, p.Arch)
		s := sha256.Sum256([]byte("archive " + name))
		if p.SHA256 != hex.EncodeToString(s[:]) {
			t.Errorf("%s sha256 %s", name, p.SHA256)
		}
		if want := "https://github.com/dphaener/rhiza/releases/download/v0.1.0/" + name; p.URL != want {
			t.Errorf("url %s, want %s", p.URL, want)
		}
	}
	b, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var back Manifest
	if err := json.Unmarshal(b, &back); err != nil || len(back.Platforms) != 4 {
		t.Fatalf("round trip: %v %s", err, b)
	}
	b2, _ := m.Marshal()
	if string(b) != string(b2) {
		t.Error("marshal not deterministic")
	}
}

func TestBuildMissingPlatform(t *testing.T) {
	_, err := Build(opts(writeDist(t, "rhiza_0.1.0_darwin_arm64.tar.gz")))
	if err == nil || !strings.Contains(err.Error(), "darwin/arm64") {
		t.Fatalf("err = %v", err)
	}
}

func TestBuildChecksumMismatch(t *testing.T) {
	dir := writeDist(t, "")
	if err := os.WriteFile(filepath.Join(dir, "rhiza_0.1.0_linux_amd64.tar.gz"), []byte("tampered"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Build(opts(dir)); err == nil || !strings.Contains(err.Error(), "sha256") {
		t.Fatalf("err = %v", err)
	}
}

func TestBuildReleaseTagNamesDownloads(t *testing.T) {
	o := opts(writeDist(t, ""))
	o.ReleaseTag = "latest"
	o.Channel = "latest"
	m, err := Build(o)
	if err != nil {
		t.Fatal(err)
	}
	if m.Channel != "latest" || m.Version != o.Tag {
		t.Fatalf("channel %q version %q", m.Channel, m.Version)
	}
	for _, p := range m.Platforms {
		if !strings.Contains(p.URL, "/releases/download/latest/") {
			t.Fatalf("url %q", p.URL)
		}
	}
}

func TestBuildLatestAssetNamesHaveNoPlus(t *testing.T) {
	dir := t.TempDir()
	var sums strings.Builder
	for _, p := range Platforms {
		name := fmt.Sprintf("rhiza_0.0.1-main.5-abc1234_%s_%s.tar.gz", p.OS, p.Arch)
		body := []byte("archive " + name)
		if err := os.WriteFile(filepath.Join(dir, name), body, 0o644); err != nil {
			t.Fatal(err)
		}
		s := sha256.Sum256(body)
		fmt.Fprintf(&sums, "%s  %s\n", hex.EncodeToString(s[:]), name)
	}
	if err := os.WriteFile(filepath.Join(dir, ChecksumsFile), []byte(sums.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	o := opts(dir)
	o.Tag = "v0.0.1-main.5+abc1234"
	o.Version = "0.0.1-main.5+abc1234"
	o.Channel = "latest"
	o.ReleaseTag = "latest"
	m, err := Build(o)
	if err != nil {
		t.Fatal(err)
	}
	if m.Version != "v0.0.1-main.5+abc1234" {
		t.Fatalf("version %q", m.Version)
	}
	for _, p := range m.Platforms {
		want := fmt.Sprintf("https://github.com/dphaener/rhiza/releases/download/latest/rhiza_0.0.1-main.5-abc1234_%s_%s.tar.gz", p.OS, p.Arch)
		if p.URL != want {
			t.Errorf("url %s, want %s", p.URL, want)
		}
	}
}
