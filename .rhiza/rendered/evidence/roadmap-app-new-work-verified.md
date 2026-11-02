# roadmap-app-new-work-verified

evidence · `01M6A71573KGPZ00ZVM014D827` · actor [claude-code-69067aa5](../actor/claude-code-69067aa5.md)

## result

PASS

## statement

At 5933776f ran go test ./internal/tui/roadmapapp -run NewWork -v -count=1: TestNewWorkStartsPlanRunWithNoChangeSet passes with a fake Starter.Exec: m, typed text, enter, y makes exactly one call, rex window new ... -- rhiza run start --workflow plan --new <intent> --milestone, with no change set argument; n gives the same without --milestone; no call before y and the view shows y/n. TestNewWorkDeclinedOrEmptyCallsNothing passes: empty enter stays in the prompt, esc at the intent prompt, n at the y/n prompt and esc at the y/n prompt make no call. Read the diff (start.go typeIntent, pending kind new, roadmapapp.go m/n keys and footers): matches the acceptance. make check exits 0 with RHIZA_ROOT/CHANGE/RUN unset.

## Links

### Out

- verifies roadmap-app-new-work (retired)

## History

- `01M686RVAD1EHGERBBDJN3S45G` The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent
