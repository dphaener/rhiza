# run-status-live-poll-waits-for-hold-qa-1

evidence · `01M7RJACJXQDNWZWV7SQ3HHJ5M` · actor [claude-code-7f9d4824](../actor/claude-code-7f9d4824.md)

## result

PASS

## statement

QA at 4c66d4717, go tests in systemd-run scope with -p 2: run-status -count=3 ok; with step first sleep 0.5 the old poll fails at txtar:38 (Current step: first), a hold-only poll fails at :37 (not running), the fixed poll passes 3 runs and with sleep 2; sibling run-status-* scripts pass; diff is only the poll line; a hand-run with the built binary showed the poll exiting at running+hold and the live run succeeding

## Links

### Out

- verifies [run-status-live-poll-waits-for-hold](../capability/run-status-live-poll-waits-for-hold.md)

## History

- `01M7RHC6VMMBXVWS3BKKFF6F8N` Fix flaky test TestScripts/run-status: failed then passed on a rerun in PR #657 at head f5b2f2b99b740b4a9a18311f074ec9e4e2df94f2; reproduce it locally, repeatedly, then fix it or quarantine it
