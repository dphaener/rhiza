package api

import (
	"github.com/dphaener/rhiza/internal/web"
)

// WebDaemonSource reads the daemon and its health for the web UI's Overview,
// from the repository base works in. Each call looks at the files as they
// are.
func WebDaemonSource(base *Ctx) web.DaemonSource {
	return func() web.DaemonView {
		c := *base
		c.Record, c.DryRun, c.Out, c.JSON = false, false, nil, true
		c.Stdin, c.Rex, c.Start, c.Detach, c.Interrupts = nil, nil, nil, nil, nil
		var v web.DaemonView
		if out, err := runDaemonStatus(&c, DaemonStatusRequest{}); err == nil {
			s := out.(DaemonStatusResponse)
			v.Running, v.Limit = s.Running, s.MaxRuns
			for _, run := range s.Runs {
				if run.State == "running" {
					v.Slots++
				}
			}
			for _, set := range s.Settings {
				if set.Key == "rhiza.daemon.paused" {
					v.Paused = set.Value == "true"
				}
			}
		}
		if out, err := runHealth(&c, nil); err == nil {
			for _, p := range out.(HealthResponse).Parts {
				v.Health = append(v.Health, web.HealthPill{Part: p.Part, State: p.State, Summary: p.Summary})
			}
		}
		return v
	}
}
