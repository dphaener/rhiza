package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestMakeUpdateDryRun checks the order of what `make update` would run: the
// in-place build into PREFIX comes first, then that binary's rex install, and
// no cp ever appears (a copied Go binary is SIGKILLed on arm64).
func TestMakeUpdateDryRun(t *testing.T) {
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make is not on PATH")
	}
	prefix := t.TempDir()
	cmd := exec.Command("make", "-n", "update", "PREFIX="+prefix)
	cmd.Dir = "../.."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make -n update: %v\n%s", err, out)
	}
	got := string(out)
	build := strings.Index(got, "go build")
	if build < 0 || !strings.Contains(got, "-o "+prefix+"/rhiza") {
		t.Fatalf("no go build -o %s/rhiza line:\n%s", prefix, got)
	}
	install := strings.Index(got, prefix+"/rhiza rex install")
	if install < 0 {
		t.Fatalf("no %s/rhiza rex install line:\n%s", prefix, got)
	}
	if build > install {
		t.Errorf("rex install comes before the build:\n%s", got)
	}
	for _, line := range strings.Split(got, "\n") {
		for _, f := range strings.Fields(line) {
			if f == "cp" {
				t.Errorf("cp appears in make -n update:\n%s", got)
			}
		}
	}
}
