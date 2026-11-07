# health-across-projects-qa-1

evidence · `01M6STV5D54N2ACS701HE06J0R` · actor [claude-code-8ba077ae](../actor/claude-code-8ba077ae.md)

## result

FAIL

## statement

QA at 538fc842 with a scratch build, read-only against the real root ~/code (nebari, rhiza). rhiza health at the root: daemon, autonomy, machine, disk and usage once untagged, then main, inbox, requests, roadmap and check once for nebari and once for rhiza, tagged (nebari: 5 behind origin/main; rhiza: 49 requests, 14 change sets, 7 check warnings), exit 0 with no bad part. Inside ~/code/rhiza and with --project nebari: one project's parts, no project tags. Defect: at the root the daemon part says 0 live runs while daemon status at the root lists 2 live runs under rhiza, and inside ~/code/rhiza health says 2 live runs. runHealth reads the machine parts through scopes[0] (nebari, first by name), so healthDaemon counts only nebari's runs; the usage part likewise reads nebari's daemon state only.

## Links

### Out

- verifies [health-across-projects](../capability/rhiza/orchestration/health-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
