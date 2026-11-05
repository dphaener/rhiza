package api_test

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	_ "github.com/dphaener/rhiza/internal/tui/autonomyapp"
	_ "github.com/dphaener/rhiza/internal/tui/daemonapp"
	_ "github.com/dphaener/rhiza/internal/tui/inboxapp"
	_ "github.com/dphaener/rhiza/internal/tui/requestapp"
	_ "github.com/dphaener/rhiza/internal/tui/runsapp"
)

// With the five apps linked in, as in the rhiza binary, rex install writes
// exactly five Rex actions with their keys in the rhiza key mode.
func TestRexInstallRegistersExactlyFiveApps(t *testing.T) {
	dir := t.TempDir()
	ctx := &api.Ctx{
		Root: t.TempDir(),
		Getenv: func(k string) string {
			if k == "REX_CONFIG" {
				return filepath.Join(dir, "init.lua")
			}
			return ""
		},
		Exec: func(string, []string, string, []string) ([]byte, error) { return nil, nil },
	}
	o, ok := api.Lookup("rex.install")
	if !ok {
		t.Fatal("rex.install is not registered")
	}
	req, err := api.DecodeRequest(o, []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := o.Run(ctx, req); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "rhiza.lua"))
	if err != nil {
		t.Fatal(err)
	}
	script := string(b)
	// The package's own tests register zz-* apps in this test binary; the
	// real binary has none, so they are not counted.
	var names []string
	for _, m := range regexp.MustCompile(`name = "(rhiza_app_[a-z0-9_]+)"`).FindAllStringSubmatch(script, -1) {
		if !strings.HasPrefix(m[1], "rhiza_app_zz_") {
			names = append(names, m[1])
		}
	}
	if n := strings.Count(script, "rex.action{") - strings.Count(script, `name = "rhiza_app_zz_`); n != 5 {
		t.Errorf("%d rex.action calls, want 5", n)
	}
	sort.Strings(names)
	want := []string{"rhiza_app_autonomy", "rhiza_app_daemon", "rhiza_app_inbox", "rhiza_app_request", "rhiza_app_runs"}
	if strings.Join(names, ",") != strings.Join(want, ",") {
		t.Errorf("actions = %v, want %v", names, want)
	}
	for app, key := range map[string]string{"autonomy": "y", "daemon": "d", "inbox": "i", "request": "q", "runs": "v"} {
		if !strings.Contains(script, `["`+app+`"] = "`+key+`",`) {
			t.Errorf("rhiza.lua lacks key %s for %s:\n%s", key, app, script)
		}
	}
}
