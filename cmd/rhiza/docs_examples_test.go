package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Every example workflow under docs/examples/workflows passes workflow check.
func TestDocsExampleWorkflows(t *testing.T) {
	files, err := filepath.Glob("../../docs/examples/workflows/*.toml")
	if err != nil || len(files) < 4 {
		t.Fatalf("examples: %v %v", files, err)
	}
	for i, f := range files {
		files[i], _ = filepath.Abs(f)
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, out)
	}
	t.Chdir(dir)
	env := func(k string) string {
		switch k {
		case "HOME", "XDG_CONFIG_HOME":
			return dir
		}
		return ""
	}
	for _, f := range files {
		var out, errb bytes.Buffer
		if code := run([]string{"workflow", "check", f}, strings.NewReader(""), &out, &errb, env); code != 0 || !strings.Contains(out.String(), "is valid") {
			t.Errorf("%s: exit %d\n%s%s", filepath.Base(f), code, out.String(), errb.String())
		}
	}
}
