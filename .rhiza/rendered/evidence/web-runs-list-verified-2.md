# web-runs-list-verified-2

evidence · `01M6Q8850E8342DM5HZP66D94C` · actor [claude-code-804a1831](../actor/claude-code-804a1831.md)

## result

PASS

## statement

At e2d82132 (fix attempt 2), rhiza show web-runs-list gives the amended statement with no short ID, which matches runRow in internal/web/runs.go (title, workflow, state pill, current step, start time, link to /runs/<id>, no ID text) and constraint owner-surfaces-no-ids. go test -trimpath -count=1 ./internal/web/... passed; TestRunsPageOrderFieldsAndLinks, TestRunsFragmentMatchesSection and TestRunsPageEmptyAndUnavailable pass and assert live-then-recent order, fields, links, IDs only in links, and that the runs section declares the run event. make self-check reports only stale-render failures, expected before landing, and no fold problems.

## Links

### Out

- verifies [web-runs-list](../capability/rhiza/web-monitor/web-runs-list.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
