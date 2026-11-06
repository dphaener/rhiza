# web-overview-runs-qa-2

evidence · `01M6QC68QFGQ7BYS5G4WJP2E47` · actor [claude-code-1f2862fd](../actor/claude-code-1f2862fd.md)

## result

PASS

## statement

At 8166e669 built rhiza (go build -trimpath) to a scratch binary, cloned the branch to a scratch repo with a copy of the real .git/rhiza/runs journals, set rhiza.daemon.paused=true, maxRuns=3, rhiza.web.port=18766 and ran the real rhiza daemon --permission-mode acceptEdits. The Live runs panel listed the live and waiting runs from the real journals (e.g. 'Web UI Overview page', ship-change, review/review; 'Let the web UI listen on all', plan, interview), each with title, workflow, current step and a /runs/<id> link and an All runs link; succeeded and failed runs, including a crafted failed run, were not listed. go test ./internal/web passes, incl. the new journal-based overview runs test.

## Links

### Out

- verifies [web-overview-runs](../capability/rhiza/web-monitor/web-overview-runs.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
