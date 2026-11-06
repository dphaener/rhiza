package api

// This file holds the daemon's settings: they live in the repo's git config
// and the daemon re-reads them every poll. A command-line flag wins over
// config for that run. rhiza daemon settings lists each with its value and
// source; rhiza daemon set validates and writes one.

import (
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/web"
)

// The sources of a setting's value.
const (
	SettingDefault = "default"
	SettingConfig  = "git config"
	SettingFlag    = "flag"
)

// settingDef is one setting: its git config key, its default (the one every
// reader uses, as shown), and the check of a value.
type settingDef struct {
	key      string
	def      string
	validate func(v string) error
	// flag names the daemon flag that overrides the key, if any.
	flag string
}

func validateChoice(opts ...string) func(string) error {
	return func(v string) error {
		for _, o := range opts {
			if v == o {
				return nil
			}
		}
		return fmt.Errorf("must be %s", strings.Join(opts, " or "))
	}
}

var settingDefs = []settingDef{
	{key: "rhiza.daemon.paused", def: "false", validate: validateChoice("true", "false")},
	{key: "rhiza.daemon.selfUpdate", def: "false", validate: validateChoice("true", "false")},
	{key: "rhiza.daemon.maxRuns", def: strconv.Itoa(daemon.DefaultMaxRuns), flag: "--max-runs", validate: func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return fmt.Errorf("must be a whole number of at least 1")
		}
		return nil
	}},
	{key: "rhiza.daemon.minFree", def: "10G", flag: "--min-free", validate: func(v string) error {
		_, err := daemon.ParseSize(v)
		return err
	}},
	{key: "rhiza.daemon.trimGoCache", def: "true", validate: validateChoice("true", "false")},
	{key: "rhiza.daemon.permissionMode", def: "", flag: "--permission-mode", validate: func(v string) error {
		if strings.TrimSpace(v) == "" || strings.ContainsAny(v, " \t\r\n") {
			return fmt.Errorf("must be one word")
		}
		return nil
	}},
	// The usage governor's windows: a limit is a token count, unset by
	// default (a window with no known limit never holds); the reserve is the
	// percent of a limit kept back.
	{key: "rhiza.daemon.usage.session", def: "", validate: validateTokenLimit},
	{key: "rhiza.daemon.usage.weekly", def: "", validate: validateTokenLimit},
	{key: "rhiza.daemon.usage.reserve", def: "20", validate: func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 || n > 100 {
			return fmt.Errorf("must be a whole number from 0 to 100")
		}
		return nil
	}},
	// The web UI's address: a literal IP, 127.0.0.1 by default.
	{key: web.HostKey, def: web.DefaultHost, validate: func(v string) error {
		_, err := web.ParseHost(v)
		return err
	}},
	// The web UI's port, or off.
	{key: web.PortKey, def: web.DefaultPort, validate: func(v string) error {
		_, err := web.ParsePort(v)
		return err
	}},
	// The ship deploy check: the health URL ship-change waits on after a
	// merge (unset means no deploy check) and how long it waits.
	{key: DeployURLKey, def: "", validate: validateDeployURL},
	{key: DeployTimeoutKey, def: DefaultDeployTimeout, validate: validateDeployTimeout},
	{key: "rhiza.autonomy.plan", def: AutonomyOwner, validate: validateChoice(AutonomyAuto, AutonomyOwner)},
	// accept governs every plan run's gate, autonomous or owner-run, which
	// accepts once the review passed unless it is owner; other change sets
	// always wait for the owner.
	{key: "rhiza.autonomy.accept", def: AutonomyAuto, validate: validateChoice(AutonomyAuto, AutonomyOwner)},
	{key: "rhiza.autonomy.ship", def: AutonomyAuto, validate: validateChoice(AutonomyAuto, AutonomyOwner)},
	{key: "rhiza.autonomy.signoff", def: AutonomyAuto, validate: validateChoice(AutonomyAuto, AutonomyOwner)},
}

// The ship deploy check's keys and the timeout's default.
const (
	DeployURLKey         = "rhiza.ship.deployUrl"
	DeployTimeoutKey     = "rhiza.ship.deployTimeout"
	DefaultDeployTimeout = "15m"
)

func validateDeployURL(v string) error {
	u, err := url.Parse(v)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("must be an http or https URL")
	}
	return nil
}

func validateDeployTimeout(v string) error {
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return fmt.Errorf("must be a positive duration such as 15m")
	}
	return nil
}

func validateTokenLimit(v string) error {
	if n, err := strconv.ParseUint(v, 10, 64); err != nil || n < 1 {
		return fmt.Errorf("must be a whole number of tokens, at least 1")
	}
	return nil
}

func findSetting(key string) (settingDef, bool) {
	for _, d := range settingDefs {
		if d.key == key {
			return d, true
		}
	}
	return settingDef{}, false
}

// Setting is one daemon setting with its effective value.
type Setting struct {
	Key     string `json:"key" desc:"The git config key"`
	Value   string `json:"value" desc:"The effective value"`
	Default string `json:"default" desc:"The value when nothing sets it"`
	Source  string `json:"source" desc:"Where the value comes from: default, git config or flag"`
}

// settingsOf reads every setting from git config, a flag in flags (key to
// value) winning. A bad value in config is a usage error naming the key.
func settingsOf(ctx *Ctx, flags map[string]string) ([]Setting, error) {
	var out []Setting
	for _, d := range settingDefs {
		s := Setting{Key: d.key, Value: d.def, Default: d.def, Source: SettingDefault}
		if ctx.Git != nil {
			if v, err := ctx.Git("config", "--get", d.key); err == nil {
				if v = strings.TrimSpace(v); v != "" {
					if err := d.validate(v); err != nil {
						return nil, Usagef("git config %s is %q: %v", d.key, v, err)
					}
					s.Value, s.Source = v, SettingConfig
				}
			}
		}
		if v, ok := flags[d.key]; ok {
			s.Value, s.Source = v, SettingFlag
		}
		out = append(out, s)
	}
	return out, nil
}

func settingValue(rows []Setting, key string) string {
	for _, s := range rows {
		if s.Key == key {
			return s.Value
		}
	}
	return ""
}

// daemonConf is the daemon's effective settings for one poll.
type daemonConf struct {
	rows        []Setting
	maxRuns     int
	minFree     uint64
	mode        string
	paused      bool
	selfUpdate  bool
	trimGoCache bool
	webHost     string
	webPort     string
}

// daemonConfOf resolves the settings for a poll: config, then the request's
// flags over it.
func daemonConfOf(ctx *Ctx, r DaemonRequest) (daemonConf, error) {
	flags := map[string]string{}
	if r.MaxRuns > 0 {
		flags["rhiza.daemon.maxRuns"] = strconv.Itoa(r.MaxRuns)
	}
	if r.MinFree != "" {
		if _, err := daemon.ParseSize(r.MinFree); err != nil {
			return daemonConf{}, Usagef("--min-free: %v", err)
		}
		flags["rhiza.daemon.minFree"] = r.MinFree
	}
	if r.PermissionMode != "" {
		flags["rhiza.daemon.permissionMode"] = r.PermissionMode
	}
	rows, err := settingsOf(ctx, flags)
	if err != nil {
		return daemonConf{}, err
	}
	c := daemonConf{rows: rows}
	c.maxRuns, _ = strconv.Atoi(settingValue(rows, "rhiza.daemon.maxRuns"))
	c.minFree, _ = daemon.ParseSize(settingValue(rows, "rhiza.daemon.minFree"))
	c.paused = settingValue(rows, "rhiza.daemon.paused") == "true"
	c.selfUpdate = settingValue(rows, "rhiza.daemon.selfUpdate") == "true"
	c.trimGoCache = settingValue(rows, "rhiza.daemon.trimGoCache") != "false"
	c.mode = settingValue(rows, "rhiza.daemon.permissionMode")
	c.webHost = settingValue(rows, web.HostKey)
	c.webPort = settingValue(rows, web.PortKey)
	return c, nil
}

func init() {
	register(Operation{
		Name:     "daemon.settings",
		ReadOnly: true,
		Summary:  "List the daemon's settings, each with its value and where it comes from",
		Request:  reflect.TypeFor[DaemonSettingsRequest](),
		Response: reflect.TypeFor[DaemonSettingsResponse](),
		Run:      runDaemonSettings,
	})
	register(Operation{
		Name:     "daemon.set",
		Summary:  "Validate and write one daemon setting to the repo's git config; a running daemon reads it at its next poll",
		Request:  reflect.TypeFor[DaemonSetRequest](),
		Response: reflect.TypeFor[DaemonSetResponse](),
		Run:      runDaemonSet,
	})
}

// DaemonSettingsRequest takes nothing.
type DaemonSettingsRequest struct{}

// DaemonSettingsResponse lists the settings.
type DaemonSettingsResponse struct {
	Settings []Setting `json:"settings" desc:"Every setting, in a fixed order"`
}

// Human writes the settings as text.
func (r DaemonSettingsResponse) Human(w io.Writer) {
	writeSettings(w, r.Settings)
}

func writeSettings(w io.Writer, rows []Setting) {
	for _, s := range rows {
		v := s.Value
		if v == "" {
			v = `""`
		}
		fmt.Fprintf(w, "%s = %s (%s)\n", s.Key, v, s.Source)
	}
}

func runDaemonSettings(ctx *Ctx, req any) (any, error) {
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return nil, err
	}
	return DaemonSettingsResponse{Settings: rows}, nil
}

// DaemonSetRequest names a setting and its new value.
type DaemonSetRequest struct {
	Key   string `json:"key" arg:"0" desc:"The setting: rhiza.daemon.paused, selfUpdate, maxRuns, minFree, permissionMode, usage.session, usage.weekly or usage.reserve, rhiza.web.host, rhiza.web.port, rhiza.ship.deployUrl or deployTimeout, or rhiza.autonomy.plan, accept, ship or signoff"`
	Value string `json:"value" arg:"1" desc:"The new value: true or false, a whole number, a token count, a percent, a size, a permission mode, an IP address, a port number or off, an http(s) URL, a duration, or auto or owner, as the key takes"`
}

// DaemonSetResponse is the setting written.
type DaemonSetResponse struct {
	Key   string `json:"key" desc:"The setting"`
	Value string `json:"value" desc:"The value written"`
}

// Human writes the response as text.
func (r DaemonSetResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "%s = %s\n", r.Key, r.Value)
}

func runDaemonSet(ctx *Ctx, req any) (any, error) {
	r := req.(DaemonSetRequest)
	d, ok := findSetting(r.Key)
	if !ok {
		keys := make([]string, 0, len(settingDefs))
		for _, s := range settingDefs {
			keys = append(keys, s.key)
		}
		return nil, Usagef("unknown setting %q: one of %s", r.Key, strings.Join(keys, ", ")).WithNext("rhiza daemon settings")
	}
	v := strings.TrimSpace(r.Value)
	if err := d.validate(v); err != nil {
		return nil, Usagef("%s: %q %v", d.key, r.Value, err).WithNext("rhiza daemon settings")
	}
	if ctx.DryRun {
		return DaemonSetResponse{Key: d.key, Value: v}, nil
	}
	if ctx.Git == nil {
		return nil, Errorf("no git to write %s", d.key)
	}
	if _, err := ctx.Git("config", d.key, v); err != nil {
		return nil, Errorf("write %s: %v", d.key, err)
	}
	return DaemonSetResponse{Key: d.key, Value: v}, nil
}

func (c daemonConf) stateRows() []daemon.SettingValue {
	var out []daemon.SettingValue
	for _, s := range c.rows {
		out = append(out, daemon.SettingValue{Key: s.Key, Value: s.Value, Source: s.Source})
	}
	return out
}

func settingDefault(key string) string {
	d, _ := findSetting(key)
	return d.def
}
