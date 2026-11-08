package api

// This file is the daemon watchdog's policy: settings declared in settingDefs
// like the others. Each name has a machine key, rhiza.daemon.watchdog.<name>
// in the root's rhiza.toml, and, but for pollBehind, a project key,
// rhiza.project.watchdog.<name> in the project's git config. A project's
// effective value is its project key when set, else the machine key, else the
// default. A workflow step's own stall field wins for that step.

import (
	"fmt"
	"strconv"
	"time"
)

// The watchdog's machine and project key prefixes.
const (
	watchdogMachinePrefix = "rhiza.daemon.watchdog."
	watchdogProjectPrefix = "rhiza.project.watchdog."
)

// Shipped defaults of the watchdog policy.
const (
	DefaultAgentNudge   = 20 * time.Minute
	DefaultAgentRetry   = 40 * time.Minute
	DefaultCommandLimit = 30 * time.Minute
	DefaultWaitLimit    = 90 * time.Minute
	DefaultRetryBudget  = 2
	DefaultCIReruns     = 1
	DefaultPollBehind   = 3
)

func validateWholeFrom(min int) func(string) error {
	return func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil || n < min {
			return fmt.Errorf("must be a whole number of at least %d", min)
		}
		return nil
	}
}

// watchdogNames lists each watchdog setting by name, in listing order.
var watchdogNames = []struct {
	name        string
	def         string
	validate    func(string) error
	machineOnly bool
}{
	{"enabled", "true", validateChoice("true", "false"), false},
	{"agentNudge", "20m", validateDeployTimeout, false},
	{"agentRetry", "40m", validateDeployTimeout, false},
	{"commandLimit", "30m", validateDeployTimeout, false},
	{"waitLimit", "90m", validateDeployTimeout, false},
	{"retryBudget", strconv.Itoa(DefaultRetryBudget), validateWholeFrom(0), false},
	{"ciReruns", strconv.Itoa(DefaultCIReruns), validateWholeFrom(0), false},
	{"pollBehind", strconv.Itoa(DefaultPollBehind), validateWholeFrom(1), true},
	{"fileBugs", "true", validateChoice("true", "false"), false},
	{"autonomousOnly", "true", validateChoice("true", "false"), false},
}

func init() {
	for _, w := range watchdogNames {
		settingDefs = append(settingDefs, settingDef{key: watchdogMachinePrefix + w.name, def: w.def, validate: w.validate})
	}
	for _, w := range watchdogNames {
		if !w.machineOnly {
			settingDefs = append(settingDefs, settingDef{project: true, key: watchdogProjectPrefix + w.name, def: w.def, validate: w.validate})
		}
	}
}

// WatchdogPolicy is the watchdog's effective policy for one project.
type WatchdogPolicy struct {
	Enabled        bool
	AgentNudge     time.Duration
	AgentRetry     time.Duration
	CommandLimit   time.Duration
	WaitLimit      time.Duration
	RetryBudget    int
	CIReruns       int
	PollBehind     int
	FileBugs       bool
	AutonomousOnly bool
}

// watchdogPolicyOf resolves the policy from rows: a project key set in git
// config over the machine key over the default.
func watchdogPolicyOf(rows []Setting) WatchdogPolicy {
	get := func(name string) string {
		for _, s := range rows {
			if s.Key == watchdogProjectPrefix+name && s.Source != SettingDefault {
				return s.Value
			}
		}
		return settingValue(rows, watchdogMachinePrefix+name)
	}
	dur := func(name string) time.Duration {
		d, _ := time.ParseDuration(get(name))
		return d
	}
	num := func(name string) int {
		n, _ := strconv.Atoi(get(name))
		return n
	}
	return WatchdogPolicy{
		Enabled:        get("enabled") == "true",
		AgentNudge:     dur("agentNudge"),
		AgentRetry:     dur("agentRetry"),
		CommandLimit:   dur("commandLimit"),
		WaitLimit:      dur("waitLimit"),
		RetryBudget:    num("retryBudget"),
		CIReruns:       num("ciReruns"),
		PollBehind:     num("pollBehind"),
		FileBugs:       get("fileBugs") == "true",
		AutonomousOnly: get("autonomousOnly") == "true",
	}
}

// WatchdogPolicyOf reads the policy ctx's project works under; a bad value
// is a usage error naming its key.
func WatchdogPolicyOf(ctx *Ctx) (WatchdogPolicy, error) {
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return WatchdogPolicy{}, err
	}
	return watchdogPolicyOf(rows), nil
}

// StallLimit is how long a step of kind may go without progress: its own
// stall field when set, else the policy's limit for its kind (the agent
// retry limit for an agent step; nudge is the earlier AgentNudge).
func (p WatchdogPolicy) StallLimit(kind string, stall time.Duration) time.Duration {
	if stall > 0 {
		return stall
	}
	switch kind {
	case "command":
		return p.CommandLimit
	case "wait", "signoff":
		return p.WaitLimit
	}
	return p.AgentRetry
}
