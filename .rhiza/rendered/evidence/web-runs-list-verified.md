# web-runs-list-verified

evidence · `01M6Q633Q56M1B1N33T89G59ME` · actor [claude-code-080b4504](../actor/claude-code-080b4504.md)

## result

PASS

## statement

Verifier claude-code-080b4504 at b457b744: go test -trimpath -count=1 ./internal/web/... passes; TestRunsPageOrderFieldsAndLinks renders /runs over a run.ls fixture with running, blocked, succeeded and failed runs and asserts live-first-then-recent order, Live/Recent groups, title, workflow, state pill, current step, start time, /runs/<id> links and IDs only in links (owner-surfaces-no-ids, which is why no short ID is shown); TestRunsFragmentMatchesSection asserts the section declares exactly EventRun (data-rz-events="run") and the page embeds /fragments/runs. A scratch test in a copy (not committed) over real run journals via run.Store and api.ReadOnlyHandler (one locked live, one finished failed, one finished succeeded) rendered the live run under Live, then failed (newer) before done (older) under Recent, with correct pills, workflows and links. gofmt and go vet clean.

## Links

### Out

- verifies [web-runs-list](../capability/rhiza/web-monitor/web-runs-list.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
