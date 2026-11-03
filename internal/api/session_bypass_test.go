package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// With neither HOME nor XDG_CONFIG_HOME the machine is trust daily: the
// placeholder path in the message is never read, so a file planted at it
// under the current directory cannot turn bypass on.
func TestBypassCheckNoConfigPathIsDaily(t *testing.T) {
	dir := t.TempDir()
	planted := filepath.Join(dir, "$XDG_CONFIG_HOME", "rhiza", "machine.toml")
	if err := os.MkdirAll(filepath.Dir(planted), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(planted, []byte("trust = \"dedicated\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	ctx := &Ctx{Getenv: func(string) string { return "" }}
	err := bypassCheck(ctx, "bypassPermissions")
	if err == nil {
		t.Fatal("bypassCheck = nil, want a refusal with no machine config path")
	}
	if !strings.Contains(err.Error(), "$XDG_CONFIG_HOME/rhiza/machine.toml") {
		t.Errorf("error %q does not name the placeholder path", err)
	}
}
