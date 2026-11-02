# roadmap-new-work-requests-verified-2

evidence · `01M6D5SJZ73FCH1H2K1Q5JTAYA` · actor [claude-code-4ddcafb6](../actor/claude-code-4ddcafb6.md)

## result

PASS

## statement

At 3ae2673f (attempt 2) ran go test ./internal/tui/roadmapapp -run NewWork -v -count=1: TestNewWorkFilesPlanRequest passes on a fake runner, typing 'Ship the thing' after n calls request 'Ship the thing --kind feature', after m 'Ship the thing --kind milestone', and the fake Exec records no rex call. The fix commit adds a fallback only when request answers not_found (y/n then rex window new ... rhiza run start --workflow plan --new <intent> [--milestone], flags that exist on RunStartRequest), covered by TestNewWorkWithoutRequestOpStartsPlanRun; the success path still runs no rex. make fmt-check vet test exit 0; make self-check fails only on stale .rhiza/rendered actor files from review sessions, not this leaf's code.

## Links

### Out

- verifies [roadmap-new-work-requests](../capability/rhiza/orchestration/daemon/plan-requests/roadmap-new-work-requests.md)

## History

- `01M6CRZQYG5PJV7X2B234DJCV4` Feature, bug and milestone requests queue in the clone and the daemon turns each into a plan run
