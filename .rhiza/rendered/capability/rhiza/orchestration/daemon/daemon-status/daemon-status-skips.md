# daemon-status-skips

capability · `01M686TKZGT5S1SHQSCK4QPPRR` · parent [daemon-status](../daemon-status.md)

## acceptance

daemon-status.txtar sets up an accepted change set serving no milestone, one with no leaves and one depending on an unlanded change set, runs daemon --once, and checks Skipped: lines with serves no milestone and the JSON skipped list with reason not-ready and its reasons; describe-golden.txtar is regenerated

## statement

rhiza daemon status and its --json list each skipped change set by ID and title with its reason, and its readiness reasons when it is not ready, as of the daemon's last poll

## Links

### Out

- decided-by [daemon-skips-saved-at-poll](../../../../../decision/daemon-skips-saved-at-poll.md)

### In

- verifies [daemon-status-skips-verified](../../../../../evidence/daemon-status-skips-verified.md)
- verifies [daemon-status-skips-qa-1](../../../../../evidence/daemon-status-skips-qa-1.md)
- verifies [pr-129-daemon-explains-skips](../../../../../evidence/pr-129-daemon-explains-skips.md)

## History

- `01M686RS47EK7F4H2WY15TNAQ5` rhiza daemon status names each accepted change set the daemon does not start and why
