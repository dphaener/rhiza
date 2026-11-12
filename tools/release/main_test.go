package main

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

func TestRunWritesManifest(t *testing.T) {
	dir := t.TempDir()
	var sums strings.Builder
	for _, p := range []string{"linux_amd64", "linux_arm64", "darwin_amd64", "darwin_arm64"} {
		name := "rhiza_0.1.0_" + p + ".tar.gz"
		body := []byte(name)
		os.WriteFile(filepath.Join(dir, name), body, 0o644)
		s := sha256.Sum256(body)
		fmt.Fprintf(&sums, "%s  %s\n", hex.EncodeToString(s[:]), name)
	}
	os.WriteFile(filepath.Join(dir, "checksums.txt"), []byte(sums.String()), 0o644)
	if err := run([]string{"-dist", dir, "-tag", "v0.1.0", "-commit", "abc"}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m struct{ Platforms []any }
	if err := json.Unmarshal(b, &m); err != nil || len(m.Platforms) != 4 {
		t.Fatalf("%v %s", err, b)
	}
}

func TestRunFailsWithoutAllPlatforms(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "checksums.txt"), nil, 0o644)
	if err := run([]string{"-dist", dir, "-tag", "v0.1.0", "-commit", "abc"}); err == nil {
		t.Fatal("want error")
	}
}
