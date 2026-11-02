# health-report-qa-1

evidence · `01M6DG493Q6A5MBMHTTWP3RSW7` · actor [claude-code-6e06158d](../actor/claude-code-6e06158d.md)

## result

PASS

## statement

At ee31e06b, QA built rhiza from the branch head and ran it by hand. On this repo, rhiza health printed all eight parts and exited 0 (inbox warn for a blocker). health --json and api call health gave the same bytes. In a fresh scratch repo with a bare origin: with no daemon, daemon was warn with 'start it with rhiza daemon', next rhiza daemon --autonomous. A missing rendered folder made check bad with fail render, and health exited 1. With RHIZA_TEST_FREE_BYTES=100 (and RHIZA_TEST_NOW), disk was bad '100 B free, under the minimum of 10.0 GiB', exit 1. A dirty README made main warn 'uncommitted changes', and a push from another clone made main warn '1 behind origin/main'. The inbox-failed case was covered only by health.txtar (passes), not by hand.

## Links

### Out

- verifies [health-report](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-report.md)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
