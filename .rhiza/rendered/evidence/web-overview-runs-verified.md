# web-overview-runs-verified

evidence · `01M6Q6SJ19YDZAQ40A2Q47D76J` · actor [claude-code-ae8f187e](../actor/claude-code-ae8f187e.md)

## result

PASS

## statement

At bf956078 ran go test -trimpath -count=1 ./internal/web/... (ok) and go test -trimpath -count=1 ./internal/api -run 'TestWebRuns|TestWeb|Daemon' (ok); go vet and gofmt -l clean. TestOverviewListsLiveRuns renders / and /fragments/runs and asserts each live run's title, workflow (<small>build</small>, <small>plan</small>), current step (leaves, gate) and link href=/runs/<id>, the waiting pill, and the No live runs empty state. TestWebRunsSourceListsLiveRuns builds real run journals with run.Store (one entered at step leaves, one finished) and asserts the Overview source returns only the unfinished run with its change set title, workflow and step. The finished-run filter lives in api.WebRunsSource (web has no internal deps), and daemon.go wires that source into web.HandlerWithRuns, so the rendered / lists only unfinished runs; the acceptance is covered across the two tests rather than one web-package test over journals. Diff read: IDs only build links, never shown; layout goldens updated.

## Links

### Out

- verifies [web-overview-runs](../capability/rhiza/web-monitor/web-overview-runs.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
