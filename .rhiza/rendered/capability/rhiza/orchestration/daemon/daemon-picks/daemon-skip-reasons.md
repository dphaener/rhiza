# daemon-skip-reasons

capability · `01M686TJN8KRFR48423FXP6BMM` · parent [daemon-picks](../daemon-picks.md)

## acceptance

A table test in internal/daemon/pick_test.go for a pure Skips function beside Pick covers live, failed-not-cleared, held and no-slot, and shows abandoned or retried failures are not skipped; an internal/api test shows daemonCandidates also reports backlog, not-ready and no-leaves skips

## statement

At each poll the daemon works out, for every accepted change set it does not start, a fixed reason (serves no milestone, no leaves, not ready with its readiness reasons, live, failed and not cleared, held, start not journaled, waiting for a slot) and saves the list in its state file

## Links

### Out

- decided-by [daemon-skips-saved-at-poll](../../../../../decision/daemon-skips-saved-at-poll.md)

### In

- verifies [daemon-skip-reasons-verified](../../../../../evidence/daemon-skip-reasons-verified.md)
- verifies [daemon-skip-reasons-qa-1](../../../../../evidence/daemon-skip-reasons-qa-1.md)
- verifies [pr-129-daemon-explains-skips](../../../../../evidence/pr-129-daemon-explains-skips.md)

## History

- `01M686RS47EK7F4H2WY15TNAQ5` rhiza daemon status names each accepted change set the daemon does not start and why
