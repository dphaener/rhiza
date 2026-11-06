# web-run-detail-verified

evidence · `01M6Q67DVYSWDF5XVQX9N0T4R7` · actor [claude-code-d80121d9](../actor/claude-code-d80121d9.md)

## result

PASS

## statement

At c2cb41f9 ran go test -trimpath -count=1 ./internal/web/... (ok), go vet and gofmt -l on internal/web (clean). internal/web/rundetail_test.go renders /runs/PLAN (fixture plan run) and asserts steps interview, draft, gate in order with done/failed/running pills, durations, times and the fail detail, the sessions per step (Session 1, finished · 15 tokens, running) with no session IDs, and the live section wired to the run event and /fragments/run-detail?run=PLAN; renders /runs/REVIEW (fixture review run) and asserts the Review section after the steps with check correctness, its failed pill and findings blocker Nil dereference a.go:3 then nit Spelling; and /runs/NOPE answers 404. Read the diff: runInfo's fields match api.RunInfo, RunSession, RunCheck rounds and RunFinding; the real run.status answers NotFoundf (HTTP 404) for an unknown run, which the page maps to 404. Noted, outside the acceptance: the session token struct has no json tags, so cache_creation and cache_read do not decode and the token sum shows input+output only.

## Links

### Out

- verifies [web-run-detail](../capability/rhiza/web-monitor/web-run-detail.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
