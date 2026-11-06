package api

import (
	"github.com/dphaener/rhiza/internal/web"
)

// WebRunsSource reads the live runs for the web UI's Overview from the
// repository base works in: every run not finished, newest first. Each call
// looks at the journals as they are.
func WebRunsSource(base *Ctx) web.RunsSource {
	return func() []web.RunView {
		c := *base
		c.Record, c.DryRun, c.Out, c.JSON = false, false, nil, true
		c.Stdin, c.Rex, c.Start, c.Detach, c.Interrupts = nil, nil, nil, nil, nil
		out, err := runRunLs(&c, RunLsRequest{Limit: runLsMaxLimit})
		if err != nil {
			return nil
		}
		var runs []web.RunView
		for _, it := range out.(RunLsResponse).Runs {
			var state string
			switch it.State {
			case runRunning:
				state = "running"
			case runBlocked, runStopped:
				state = "waiting"
			default:
				continue // succeeded or failed: finished
			}
			title := it.Title
			if title == "" {
				title = it.Workflow
			}
			runs = append(runs, web.RunView{ID: it.ID, Title: oneLine(title), Workflow: it.Workflow, Step: it.Current, State: state, Started: it.StartedAt})
		}
		return runs
	}
}
