# roadmap-new-work-requests-verified

evidence · `01M6D4PAD824WPQMTV197PYBP1` · actor [claude-code-59f3f7e8](../actor/claude-code-59f3f7e8.md)

## result

PASS

## statement

At f9b9b42c ran go test ./internal/tui/roadmapapp -run TestNewWorkFilesPlanRequest -v -count=1: PASS. The test drives the roadmap app on a fake runner and a fake rex Exec: key n, type 'Ship the thing', enter gives exactly one request call with args 'Ship the thing --kind feature'; key m gives 'Ship the thing --kind milestone'; in both cases the rex Exec records zero calls (no rex window new) and the view shows 'Queued'. Read the diff: typeIntent's enter now calls fileRequest (m.run("request", ...)) instead of setting a pending 'new' run, and the 'new' branches of pending.prompt/args/confirm are removed. go test ./internal/tui/... ./cmd/rhiza, go vet and gofmt -l also clean. Note: the request operation itself is not registered in this build; it belongs to the parent plan-requests capability, outside this leaf's acceptance.

## Links

### Out

- verifies roadmap-new-work-requests (retired)

## History

- `01M6CRZQYG5PJV7X2B234DJCV4` Feature, bug and milestone requests queue in the clone and the daemon turns each into a plan run
