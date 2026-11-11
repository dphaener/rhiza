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
	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/web"
)

// The scopes of a setting: a machine setting lives in the root's rhiza.toml,
// a project setting in that project's git config.
const (
	ScopeMachine = "machine"
	ScopeProject = "project"
)

// The sources of a setting's value.
const (
	SettingDefault = "default"
	SettingConfig  = "git config"
	SettingRoot    = "rhiza.toml"
	SettingFlag    = "flag"
)

// settingDef is one setting: its git config key, its default (the one every
// reader uses, as shown), and the check of a value.
type settingDef struct {
	key      string
	def      string
	validate func(v string) error
	// typ names the kind of value and desc says what the setting does; the
	// generated reference prints both.
	typ  string
	desc string
	// flag names the daemon flag that overrides the key, if any.
	flag string
	// project marks a project setting; the rest are machine settings.
	project bool
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
	{key: "rhiza.daemon.paused", def: "false", typ: "boolean", desc: "Pause the whole daemon: no project is worked while true", validate: validateChoice("true", "false")},
	// The project's own pause and run cap: the project is not worked while
	// either pause holds, and runs at most this many at once.
	{project: true, key: "rhiza.project.paused", def: "false", typ: "boolean", desc: "Pause this project: it is not worked while true", validate: validateChoice("true", "false")},
	{project: true, key: "rhiza.project.maxRuns", def: "", typ: "whole number", desc: "The most runs this project has at once; unset means only the daemon's cap applies", validate: func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return fmt.Errorf("must be a whole number of at least 1")
		}
		return nil
	}},
	// The hook commands run through sh -c, per project.
	{project: true, key: "rhiza.eventCommand", def: "", typ: "shell command", desc: "A command run through sh -c for each event; unset means none", validate: validateCommand},
	{project: true, key: "rhiza.notifyCommand", def: "", typ: "shell command", desc: "A command run through sh -c to notify the owner; unset means none", validate: validateCommand},
	{key: "rhiza.daemon.selfUpdate", def: "false", typ: "boolean", desc: "Let the daemon update and restart itself when main has a newer build", validate: validateChoice("true", "false")},
	{key: "rhiza.daemon.maxRuns", def: strconv.Itoa(daemon.DefaultMaxRuns), flag: "--max-runs", typ: "whole number", desc: "The most runs the daemon works at once, across all projects", validate: func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return fmt.Errorf("must be a whole number of at least 1")
		}
		return nil
	}},
	{key: "rhiza.daemon.minFree", def: "10G", flag: "--min-free", typ: "size", desc: "The free disk space below which the daemon starts no run, such as 10G", validate: func(v string) error {
		_, err := daemon.ParseSize(v)
		return err
	}},
	{key: "rhiza.daemon.maxLoad", def: "3.0", flag: "--max-load", typ: "number", desc: "The 1-minute load per core above which the daemon starts no run until the 5-minute load is back under it; 0 turns it off", validate: func(v string) error {
		_, err := parseMaxLoad(v)
		return err
	}},
	{key: "rhiza.daemon.minMemory", def: "2G", flag: "--min-memory", typ: "size", desc: "The available memory below which the daemon starts no run, such as 2G; 0 turns it off", validate: func(v string) error {
		_, err := daemon.ParseSize(v)
		return err
	}},
	{key: "rhiza.daemon.trimGoCache", def: "true", typ: "boolean", desc: "Trim the Go build cache when disk space runs low", validate: validateChoice("true", "false")},
	{key: "rhiza.daemon.permissionMode", def: "", flag: "--permission-mode", typ: "one word", desc: "The permission mode agents run in; unset leaves the agent's own default", validate: func(v string) error {
		if strings.TrimSpace(v) == "" || strings.ContainsAny(v, " \t\r\n") {
			return fmt.Errorf("must be one word")
		}
		return nil
	}},
	// The usage governor's windows: a limit is a token count, unset by
	// default (a window with no known limit never holds); the reserve is the
	// percent of a limit kept back.
	{key: "rhiza.daemon.usage.session", def: "", typ: "token count", desc: "The token limit of the usage session window; unset means no known limit, so the window never holds", validate: validateTokenLimit},
	{key: "rhiza.daemon.usage.weekly", def: "", typ: "token count", desc: "The token limit of the usage weekly window; unset means no known limit, so the window never holds", validate: validateTokenLimit},
	{key: "rhiza.daemon.usage.reserve", def: "20", typ: "percent, 0 to 100", desc: "The percent of a usage limit kept back before the governor holds runs", validate: func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 || n > 100 {
			return fmt.Errorf("must be a whole number from 0 to 100")
		}
		return nil
	}},
	// The web UI's address: a literal IP, 127.0.0.1 by default.
	{key: web.HostKey, def: web.DefaultHost, typ: "IP address", desc: "The address the web UI is served on, a literal IP", validate: func(v string) error {
		_, err := web.ParseHost(v)
		return err
	}},
	// The web UI's port, or off.
	{key: web.PortKey, def: web.DefaultPort, typ: "port or off", desc: "The port the web UI is served on, or off to serve none", validate: func(v string) error {
		_, err := web.ParsePort(v)
		return err
	}},
	// The ship deploy check: the health URL ship-change waits on after a
	// merge (unset means no deploy check) and how long it waits.
	{project: true, key: DeployURLKey, def: "", typ: "http or https URL", desc: "The health URL ship-change waits on after a merge; unset means no deploy check", validate: validateDeployURL},
	{project: true, key: DeployTimeoutKey, def: DefaultDeployTimeout, typ: "duration", desc: "How long ship-change waits for the deploy check, such as 15m", validate: validateDeployTimeout},
	// The CI stability window: how many of the latest settled CI runs in the
	// stability record are checked for flakes.
	{project: true, key: CIWindowKey, def: strconv.Itoa(DefaultCIWindow), typ: "whole number", desc: "How many of the latest settled CI runs are checked for flakes when judging CI stability", validate: validateWholeFrom(1)},
	// The local gate: the command ship-change runs in the run's worktree
	// before every push while CI is unstable (unset means none).
	{project: true, key: LocalGateKey, def: "", typ: "shell command", desc: "The command ship-change runs in the run's worktree before every push while CI is unstable; unset means none", validate: validateCommand},
	{project: true, key: "rhiza.autonomy.plan", def: AutonomyOwner, typ: "auto or owner", desc: "Who plans: auto lets an agent plan a request, owner waits for the owner", validate: validateChoice(AutonomyAuto, AutonomyOwner)},
	// accept governs every plan run's gate, autonomous or owner-run, which
	// accepts once the review passed unless it is owner; other change sets
	// always wait for the owner.
	{project: true, key: "rhiza.autonomy.accept", def: AutonomyAuto, typ: "auto or owner", desc: "Who accepts a change set once its review passed: auto accepts a plan run's change set, owner waits for the owner", validate: validateChoice(AutonomyAuto, AutonomyOwner)},
	{project: true, key: "rhiza.autonomy.ship", def: AutonomyAuto, typ: "auto or owner", desc: "Who ships an accepted change set: auto merges it, owner waits for the owner", validate: validateChoice(AutonomyAuto, AutonomyOwner)},
	{project: true, key: "rhiza.autonomy.signoff", def: AutonomyAuto, typ: "auto or owner", desc: "Who signs off a shipped change set: auto signs off, owner waits for the owner", validate: validateChoice(AutonomyAuto, AutonomyOwner)},
}

// The CI stability window's key and default.
const (
	CIWindowKey     = "rhiza.ship.ciWindow"
	DefaultCIWindow = 20
)

// The ship deploy check's keys and the timeout's default.
const (
	DeployURLKey         = "rhiza.ship.deployUrl"
	DeployTimeoutKey     = "rhiza.ship.deployTimeout"
	DefaultDeployTimeout = "15m"
)

func validateCommand(v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("must not be empty")
	}
	return nil
}

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
	Source  string `json:"source" desc:"Where the value comes from: default, rhiza.toml, git config or flag"`
	Scope   string `json:"scope" desc:"machine (kept in the root's rhiza.toml) or project (kept in the project's git config)"`
}

func (d settingDef) scope() string {
	if d.project {
		return ScopeProject
	}
	return ScopeMachine
}

// rootOf finds the Rhiza root that ctx works under and loads its config.
func rootOf(ctx *Ctx) (string, rootcfg.Config, bool) {
	if ctx.Wd == "" {
		return "", rootcfg.Config{}, false
	}
	root, ok := rootcfg.Find(ctx.Wd)
	if !ok {
		return "", rootcfg.Config{}, false
	}
	c, err := rootcfg.Load(root)
	if err != nil {
		return "", rootcfg.Config{}, false
	}
	return root, c, true
}

// settingsOf reads every setting from git config, a flag in flags (key to
// value) winning. A bad value in config is a usage error naming the key.
func settingsOf(ctx *Ctx, flags map[string]string) ([]Setting, error) {
	var out []Setting
	_, cfg, haveRoot := rootOf(ctx)
	for _, d := range settingDefs {
		s := Setting{Key: d.key, Value: d.def, Default: d.def, Source: SettingDefault, Scope: d.scope()}
		if !d.project && haveRoot {
			if v := strings.TrimSpace(cfg.Daemon[d.key]); v != "" {
				if err := d.validate(v); err != nil {
					return nil, Usagef("rhiza.toml %s is %q: %v", d.key, v, err)
				}
				s.Value, s.Source = v, SettingRoot
			}
		}
		// A machine key not in rhiza.toml still reads an older git config
		// value; a project key reads only git config.
		if s.Source == SettingDefault && ctx.Git != nil {
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
	rows    []Setting
	maxRuns int
	// projectMax is the project's own cap on its running runs within the
	// machine-wide maxRuns, 0 for none.
	projectMax int
	minFree    uint64
	maxLoad    float64
	minMemory  uint64
	// machine is the hold of new runs for load or memory this poll found.
	machine     *daemon.MachineHold
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
	if r.MaxLoad != "" {
		if _, err := parseMaxLoad(r.MaxLoad); err != nil {
			return daemonConf{}, Usagef("--max-load: %v", err)
		}
		flags["rhiza.daemon.maxLoad"] = r.MaxLoad
	}
	if r.MinMemory != "" {
		if _, err := daemon.ParseSize(r.MinMemory); err != nil {
			return daemonConf{}, Usagef("--min-memory: %v", err)
		}
		flags["rhiza.daemon.minMemory"] = r.MinMemory
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
	c.maxLoad, _ = parseMaxLoad(settingValue(rows, "rhiza.daemon.maxLoad"))
	c.minMemory, _ = daemon.ParseSize(settingValue(rows, "rhiza.daemon.minMemory"))
	c.paused = settingValue(rows, "rhiza.daemon.paused") == "true" || settingValue(rows, "rhiza.project.paused") == "true"
	c.projectMax, _ = strconv.Atoi(settingValue(rows, "rhiza.project.maxRuns"))
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

// DaemonSettingsRequest names the project whose settings to list.
type DaemonSettingsRequest struct {
	Project string `json:"project,omitempty" flag:"project" desc:"The registered project whose project settings to list; default the current project"`
}

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
		fmt.Fprintf(w, "%s = %s (%s, %s)\n", s.Key, v, s.Source, s.Scope)
	}
}

// scopeCtx is the context a setting's scope is read or written in: the
// named project's for a project key given --project, else ctx itself. At
// the root with no project, a project key is a usage error.
func scopeCtx(ctx *Ctx, project string, needProject bool) (*Ctx, error) {
	root, cfg, ok := rootOf(ctx)
	if project != "" {
		if !ok {
			return nil, Usagef("--project needs a Rhiza root").WithNext("rhiza root init")
		}
		if _, reg := cfg.Projects[project]; !reg {
			return nil, Usagef("project %q is not registered", project).WithNext("rhiza project ls")
		}
		return projectCtx(ctx, cfg.Dir(root, project)), nil
	}
	if needProject && ok {
		if name, in := cfg.ProjectOf(root, ctx.Wd); !in {
			_ = name
			return nil, Usagef("%s is not inside a project: name one with --project", ctx.Wd).WithNext("rhiza daemon set --project <name> <key> <value>")
		}
	}
	return ctx, nil
}

func runDaemonSettings(ctx *Ctx, req any) (any, error) {
	pc, err := scopeCtx(ctx, req.(DaemonSettingsRequest).Project, false)
	if err != nil {
		return nil, err
	}
	rows, err := settingsOf(pc, nil)
	if err != nil {
		return nil, err
	}
	return DaemonSettingsResponse{Settings: rows}, nil
}

// DaemonSetRequest names a setting and its new value.
type DaemonSetRequest struct {
	Key     string `json:"key" arg:"0" desc:"The setting: rhiza.daemon.paused, selfUpdate, maxRuns, minFree, maxLoad, minMemory, permissionMode, usage.session, usage.weekly or usage.reserve, rhiza.daemon.watchdog.* or rhiza.project.watchdog.* (enabled, agentNudge, agentRetry, commandLimit, waitLimit, retryBudget, ciReruns, pollBehind, fileBugs, autonomousOnly), rhiza.web.host, rhiza.web.port, rhiza.ship.deployUrl, deployTimeout or ciWindow, or rhiza.autonomy.plan, accept, ship or signoff"`
	Project string `json:"project,omitempty" flag:"project" desc:"The registered project a project setting is written to; default the current project"`
	Value   string `json:"value" arg:"1" desc:"The new value: true or false, a whole number, a token count, a percent, a size, a permission mode, an IP address, a port number or off, an http(s) URL, a duration, or auto or owner, as the key takes"`
}

// DaemonSetResponse is the setting written.
type DaemonSetResponse struct {
	Key   string `json:"key" desc:"The setting"`
	Value string `json:"value" desc:"The value written"`
	Scope string `json:"scope" desc:"Where it was written: machine (rhiza.toml) or project (git config)"`
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
	if !d.project {
		if root, _, ok := rootOf(ctx); ok {
			if ctx.DryRun {
				return DaemonSetResponse{Key: d.key, Value: v, Scope: ScopeMachine}, nil
			}
			if err := rootcfg.SetDaemon(root, d.key, v); err != nil {
				return nil, Errorf("write %s: %v", d.key, err)
			}
			return DaemonSetResponse{Key: d.key, Value: v, Scope: ScopeMachine}, nil
		}
	}
	wc := ctx
	if d.project {
		var err error
		if wc, err = scopeCtx(ctx, r.Project, true); err != nil {
			return nil, err
		}
	}
	if ctx.DryRun {
		return DaemonSetResponse{Key: d.key, Value: v, Scope: d.scope()}, nil
	}
	if wc.Git == nil {
		return nil, Errorf("no git to write %s", d.key)
	}
	if _, err := wc.Git("config", d.key, v); err != nil {
		return nil, Errorf("write %s: %v", d.key, err)
	}
	return DaemonSetResponse{Key: d.key, Value: v, Scope: d.scope()}, nil
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

func settingScope(key string) string {
	d, _ := findSetting(key)
	return d.scope()
}

// parseMaxLoad reads rhiza.daemon.maxLoad: a load per core, 0 for off.
func parseMaxLoad(v string) (float64, error) {
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil || f < 0 {
		return 0, fmt.Errorf("must be a load per core of 0 or more, such as 3.0")
	}
	return f, nil
}
