# web-overview-runs-verified-2

evidence · `01M6Q9J4M57QSXW96GS21G8BQC` · actor [claude-code-1e0ffff9](../actor/claude-code-1e0ffff9.md)

## result

PASS

## statement

At dde95161 ran go test -trimpath -count=1 -run TestOverview -v ./internal/web/: all pass, including TestOverviewListsLiveRunsFromJournals (internal/web/overview_runs_test.go). It writes a running journal (build, step leaves) and a finished journal with run.Store in a fixture repo, serves / through web.HandlerWithRuns over api.WebRunsSource, and asserts the live run's title, <small>build</small> workflow, leaves step and href=/runs/01LIVE link, and that the finished run's ID and title are absent. Read the diff and WebRunsSource: succeeded/failed runs are skipped, so the absence comes from the finished filter. This answers compliance finding f1.

## Links

### Out

- verifies [web-overview-runs](../capability/rhiza/web-monitor/web-overview-runs.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
