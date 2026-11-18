package api

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/web"
	"github.com/dphaener/rhiza/internal/workflow"
)

// gitConfigFake answers a fake Git's config reads from cfg, as git does:
// git config --get key, and the one snapshot read settingsOf makes (git
// config -z --get-regexp ^rhiza\.), each entry key, newline, value, NUL. A
// key of cfg is written as git prints it, canonical or not. ok is false for
// any other git call.
func gitConfigFake(cfg map[string]string, args []string) (out string, err error, ok bool) {
	switch {
	case len(args) == 3 && args[0] == "config" && args[1] == "--get":
		for k, v := range cfg {
			if canonicalConfigKey(k) == canonicalConfigKey(args[2]) {
				return v + "\n", nil, true
			}
		}
		return "", errors.New("unset"), true
	case len(args) == 4 && args[0] == "config" && args[1] == "-z" && args[2] == "--get-regexp" && args[3] == `^rhiza\.`:
		var b strings.Builder
		for k, v := range cfg {
			if strings.HasPrefix(canonicalConfigKey(k), "rhiza.") {
				b.WriteString(canonicalConfigKey(k) + "\n" + v + "\x00")
			}
		}
		if b.Len() == 0 {
			return "", errors.New("exit status 1"), true
		}
		return b.String(), nil, true
	}
	return "", nil, false
}

// TestSettingsOfOneGitRead shows settingsOf reading git config once, with the
// same rows, sources and errors as a git config --get per key gave.
func TestSettingsOfOneGitRead(t *testing.T) {
	t.Parallel()
	// The snapshot git config -z --get-regexp prints: keys canonical (section
	// and variable lowercased, subsection kept), a key set twice listed
	// twice, the last value winning.
	snapshot := "rhiza.ship.prepushcheck\nmake self-check\x00" +
		"rhiza.daemon.maxruns\n2\x00" +
		"rhiza.daemon.maxruns\n5\x00" +
		"rhiza.project.watchdog.enabled\nfalse\x00" +
		"rhiza.project.Watchdog.fileBugs\nfalse\x00" + // another subsection: not rhiza.project.watchdog
		"rhiza.autonomy.ship\nowner\x00" +
		"rhiza.web.token\nsecret-token-value\x00" +
		"rhiza.bare\x00"
	// What git config --get answers for each key of the same config.
	get := map[string]string{
		"rhiza.ship.prePushCheck":        "make self-check",
		"rhiza.daemon.maxRuns":           "5",
		"rhiza.project.watchdog.enabled": "false",
		"rhiza.autonomy.ship":            "owner",
		"rhiza.web.token":                "secret-token-value",
	}
	counting := func(snap string, get map[string]string) (*Ctx, *[][]string) {
		var calls [][]string
		ctx := &Ctx{Git: func(args ...string) (string, error) {
			calls = append(calls, args)
			if len(args) == 4 && args[0] == "config" && args[2] == "--get-regexp" {
				return snap, nil
			}
			if len(args) == 3 && args[0] == "config" && args[1] == "--get" {
				if v, ok := get[args[2]]; ok {
					return v + "\n", nil
				}
			}
			return "", errors.New("unset")
		}}
		return ctx, &calls
	}
	// want is what a git config --get per key gives: the rows settingsOf
	// returned before the snapshot.
	want := func(get map[string]string) ([]Setting, error) {
		var out []Setting
		for _, d := range settingDefs {
			s := Setting{Key: d.key, Value: d.def, Default: d.def, Source: SettingDefault, Scope: d.scope()}
			if v := strings.TrimSpace(get[d.key]); v != "" {
				if err := d.validate(v); err != nil {
					return nil, Usagef("git config %s is %q: %v", d.key, v, err)
				}
				s.Value, s.Source = v, SettingConfig
			}
			if d.key == web.TokenKey && s.Value != "" {
				s.Value = web.TokenSet
			}
			out = append(out, s)
		}
		return out, nil
	}

	ctx, calls := counting(snapshot, get)
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(*calls) != 1 {
		t.Fatalf("settingsOf ran git %d times, want once: %q", len(*calls), *calls)
	}
	wantRows, _ := want(get)
	if !reflect.DeepEqual(rows, wantRows) {
		t.Fatalf("rows differ from git config --get's:\n got %+v\nwant %+v", rows, wantRows)
	}
	for key, w := range map[string]struct{ value, source string }{
		"rhiza.ship.prePushCheck":         {"make self-check", SettingConfig},
		"rhiza.daemon.maxRuns":            {"5", SettingConfig},
		"rhiza.project.watchdog.enabled":  {"false", SettingConfig},
		"rhiza.project.watchdog.fileBugs": {"true", SettingDefault},
		"rhiza.autonomy.ship":             {"owner", SettingConfig},
	} {
		for _, s := range rows {
			if s.Key == key && (s.Value != w.value || s.Source != w.source) {
				t.Errorf("%s = %q (%s), want %q (%s)", key, s.Value, s.Source, w.value, w.source)
			}
		}
	}
	if v := settingValue(rows, "rhiza.web.token"); v == "secret-token-value" {
		t.Fatalf("the web token's value is shown: %q", v)
	}

	t.Run("real git", func(t *testing.T) {
		// A real repo, its global and system config shut out: settingsOf's
		// one snapshot against a git config --get per key.
		dir := t.TempDir()
		git := func(args ...string) (string, error) {
			c := exec.Command("git", append([]string{"-C", dir}, args...)...)
			c.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1")
			out, err := c.Output()
			return string(out), err
		}
		set := func(args ...string) {
			t.Helper()
			if out, err := git(append([]string{"config"}, args...)...); err != nil {
				t.Fatalf("git config %v: %v %s", args, err, out)
			}
		}
		if out, err := git("init", "-q"); err != nil {
			t.Fatalf("git init: %v %s", err, out)
		}
		set("rhiza.ship.prePushCheck", "make self-check")
		set("rhiza.daemon.maxRuns", "2")
		set("--add", "rhiza.daemon.maxRuns", "5")
		set("rhiza.project.watchdog.enabled", "false")
		set("rhiza.project.Watchdog.fileBugs", "false")
		set("RHIZA.autonomy.SHIP", "owner")
		var calls int
		ctx := &Ctx{Git: func(args ...string) (string, error) { calls++; return git(args...) }}
		perKey := func() map[string]string {
			m := map[string]string{}
			for _, d := range settingDefs {
				if v, err := git("config", "--get", d.key); err == nil {
					m[d.key] = v
				}
			}
			return m
		}
		rows, err := settingsOf(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if calls != 1 {
			t.Fatalf("settingsOf ran git %d times, want once", calls)
		}
		wantRows, _ := want(perKey())
		if !reflect.DeepEqual(rows, wantRows) {
			t.Fatalf("rows differ from git config --get's:\n got %+v\nwant %+v", rows, wantRows)
		}
		if settingValue(rows, "rhiza.daemon.maxRuns") != "5" || settingValue(rows, "rhiza.project.watchdog.enabled") != "false" || settingValue(rows, "rhiza.autonomy.ship") != AutonomyOwner || settingValue(rows, "rhiza.ship.prePushCheck") != "make self-check" {
			t.Fatalf("the set keys did not read: %+v", rows)
		}
		set("rhiza.daemon.minFree", "roomy")
		calls = 0
		_, err = settingsOf(ctx, nil)
		_, wantErr := want(perKey())
		if err == nil || wantErr == nil || err.Error() != wantErr.Error() || !strings.Contains(err.Error(), "rhiza.daemon.minFree") {
			t.Fatalf("got error %v, want %v", err, wantErr)
		}
		if calls != 1 {
			t.Fatalf("settingsOf with a bad value ran git %d times, want once", calls)
		}
	})

	t.Run("bad value", func(t *testing.T) {
		bad := map[string]string{"rhiza.daemon.maxRuns": "lots"}
		ctx, calls := counting("rhiza.daemon.maxruns\nlots\x00", bad)
		_, err := settingsOf(ctx, nil)
		_, wantErr := want(bad)
		if err == nil || wantErr == nil || err.Error() != wantErr.Error() {
			t.Fatalf("got error %v, want %v", err, wantErr)
		}
		if !strings.Contains(err.Error(), "rhiza.daemon.maxRuns") {
			t.Fatalf("error %q does not name the key", err)
		}
		if len(*calls) != 1 {
			t.Fatalf("ran git %d times, want once", len(*calls))
		}
	})

	t.Run("no keys", func(t *testing.T) {
		// git config --get-regexp exits 1 when nothing matches.
		ctx := &Ctx{Git: func(args ...string) (string, error) { return "", errors.New("exit status 1") }}
		rows, err := settingsOf(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		wantRows, _ := want(nil)
		if !reflect.DeepEqual(rows, wantRows) {
			t.Fatalf("got %+v, want the defaults %+v", rows, wantRows)
		}
	})

	t.Run("autonomy", func(t *testing.T) {
		ctx, calls := counting(snapshot, get)
		for stage, want := range map[string]string{"plan": AutonomyOwner, "accept": AutonomyAuto, "ship": AutonomyOwner, "signoff": AutonomyAuto} {
			*calls = nil
			m, err := autonomyStage(ctx, stage)
			if err != nil || m != want {
				t.Fatalf("autonomy %s = %q, %v; want %q", stage, m, err, want)
			}
			if len(*calls) != 1 {
				t.Fatalf("autonomy %s ran git %d times, want once", stage, len(*calls))
			}
		}
		*calls = nil
		modes, err := runAutonomy(ctx, &workflow.Definition{Name: "plain"})
		if err != nil {
			t.Fatal(err)
		}
		if len(modes) != len(workflow.AutonomyStages) || len(*calls) != 1 {
			t.Fatalf("runAutonomy read %d stages with %d git calls, want %d with one", len(modes), len(*calls), len(workflow.AutonomyStages))
		}
		ctx, _ = counting("rhiza.autonomy.plan\nsometimes\x00", nil)
		if _, err := autonomyStage(ctx, "plan"); err == nil || !strings.Contains(err.Error(), "rhiza.autonomy.plan") {
			t.Fatalf("a bad autonomy value: got %v", err)
		}
	})
}

func TestCanonicalConfigKey(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]string{
		"rhiza.ship.prePushCheck":        "rhiza.ship.prepushcheck",
		"RHIZA.eventCommand":             "rhiza.eventcommand",
		"rhiza.project.Watchdog.enabled": "rhiza.project.Watchdog.enabled",
		"rhiza.project.watchdog.ENABLED": "rhiza.project.watchdog.enabled",
		"rhiza":                          "rhiza",
	} {
		if got := canonicalConfigKey(in); got != want {
			t.Errorf("canonicalConfigKey(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestDaemonMaxRunsAuto shows rhiza.daemon.maxRuns taking auto beside a
// whole number: daemon set writes it, daemon settings reads it back, and
// daemonConf is marked auto with the default cap kept for its readers; a
// number pins the cap. 0, -1 and autox are refused, and auto is refused for
// rhiza.project.maxRuns.
func TestDaemonMaxRunsAuto(t *testing.T) {
	t.Parallel()
	cfg := map[string]string{}
	ctx := &Ctx{Git: func(args ...string) (string, error) {
		if len(args) == 3 && args[0] == "config" && args[1] != "--get" && args[1] != "-z" {
			cfg[args[1]] = args[2]
			return "", nil
		}
		out, err, ok := gitConfigFake(cfg, args)
		if !ok {
			t.Fatalf("unexpected git %v", args)
		}
		return out, err
	}}
	set := func(key, v string) error {
		_, err := runDaemonSet(ctx, DaemonSetRequest{Key: key, Value: v})
		return err
	}

	if err := set("rhiza.daemon.maxRuns", "auto"); err != nil {
		t.Fatal(err)
	}
	out, err := runDaemonSettings(ctx, DaemonSettingsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if v := settingValue(out.(DaemonSettingsResponse).Settings, "rhiza.daemon.maxRuns"); v != "auto" {
		t.Fatalf("daemon settings read maxRuns %q, want auto", v)
	}
	c, err := daemonConfOf(ctx, DaemonRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !c.maxRunsAuto || c.maxRuns != daemon.DefaultMaxRuns {
		t.Fatalf("auto: maxRunsAuto %v maxRuns %d, want true and %d", c.maxRunsAuto, c.maxRuns, daemon.DefaultMaxRuns)
	}

	if err := set("rhiza.daemon.maxRuns", "3"); err != nil {
		t.Fatal(err)
	}
	if c, err = daemonConfOf(ctx, DaemonRequest{}); err != nil {
		t.Fatal(err)
	}
	if c.maxRunsAuto || c.maxRuns != 3 {
		t.Fatalf("3: maxRunsAuto %v maxRuns %d, want false and 3", c.maxRunsAuto, c.maxRuns)
	}
	// --max-runs still pins over auto in config.
	cfg["rhiza.daemon.maxRuns"] = "auto"
	if c, err = daemonConfOf(ctx, DaemonRequest{MaxRuns: 4}); err != nil {
		t.Fatal(err)
	}
	if c.maxRunsAuto || c.maxRuns != 4 {
		t.Fatalf("--max-runs 4 over auto: maxRunsAuto %v maxRuns %d", c.maxRunsAuto, c.maxRuns)
	}

	for _, v := range []string{"0", "-1", "autox"} {
		err := set("rhiza.daemon.maxRuns", v)
		if err == nil || !strings.Contains(err.Error(), "a whole number of at least 1 or auto") {
			t.Errorf("maxRuns %q: error %v, want one naming 'a whole number of at least 1 or auto'", v, err)
		}
	}
	if cfg["rhiza.daemon.maxRuns"] != "auto" {
		t.Errorf("a refused value was written: %q", cfg["rhiza.daemon.maxRuns"])
	}
	if err := set("rhiza.project.maxRuns", "auto"); err == nil || !strings.Contains(err.Error(), "must be a whole number of at least 1") || strings.Contains(err.Error(), "or auto") {
		t.Errorf("project maxRuns auto: error %v, want a whole-number refusal", err)
	}
}

// TestDaemonMaxRunsAutoCap checks the cap under rhiza.daemon.maxRuns=auto:
// one over the runs running while every pre-start check has room, the runs
// running while the load or memory hold, the tmp-low hold or the disk-low
// check holds.
func TestDaemonMaxRunsAutoCap(t *testing.T) {
	room := daemonConf{maxRunsAuto: true}
	machine := room
	machine.machine = &daemon.MachineHold{Reason: daemon.SkipLoadHigh}
	tmp := room
	tmp.tmpLow = &daemon.DiskLow{Free: 1, Min: 2}
	disk := room
	disk.diskLow = true
	for _, c := range []struct {
		name    string
		conf    daemonConf
		running int
		want    int
	}{
		{"room, none running", room, 0, 1},
		{"room, two running", room, 2, 3},
		{"machine hold", machine, 2, 2},
		{"tmp low", tmp, 1, 1},
		{"disk low", disk, 0, 0},
	} {
		if got := c.conf.autoMaxRuns(c.running); got != c.want {
			t.Errorf("%s: autoMaxRuns(%d) = %d, want %d", c.name, c.running, got, c.want)
		}
	}
}
