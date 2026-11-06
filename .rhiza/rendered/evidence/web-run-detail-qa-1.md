# web-run-detail-qa-1

evidence · `01M6Q7JXPWNMJ4XT50XW99WGZJ` · actor [claude-code-1544fbc8](../actor/claude-code-1544fbc8.md)

## result

FAIL

## statement

At 6d5ec8e2 with the built rhiza's real daemon web UI over copied real run journals: GET /runs/<id> answered 404 for every run (a ship-change run, a plan run, a succeeded run) as well as for an unknown one, and /fragments/run-detail?run=<id> rendered No such run. GET /api/run.status through the same server answers 404 no read-only operation "run.status": run.status (internal/api/run.go:55) is not marked ReadOnly, so ReadOnlyHandler refuses it and runDetail maps the 404 to a page 404. The package tests pass only because they serve run.status from a fake API handler.

## Links

### Out

- verifies [web-run-detail](../capability/rhiza/web-monitor/web-run-detail.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
