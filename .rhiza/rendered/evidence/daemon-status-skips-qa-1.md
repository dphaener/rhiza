# daemon-status-skips-qa-1

evidence · `01M69510VQ0RAM71XBEV82BQX7` · actor [claude-code-150282ee](../actor/claude-code-150282ee.md)

## result

PASS

## statement

At 9ea6aab with the built binary in the same scratch repo: rhiza daemon status printed 'Skipped: <id> Leafless: has no leaves', 'Skipped: <id> Waiting: not ready (depends on 01M694VFHJFPQ0 (accepted))', 'Skipped: <id> Orphan: serves no milestone', 'Skipped: <id> Beta: waiting for a slot', and after the planted runs '... Alpha: has a live run' and '... Beta: its latest run failed and waits in the inbox'. rhiza daemon status --json listed skipped entries with change, title, reason (no-leaves, not-ready, no-milestone, no-slot, then live, failed) and detail (the readiness reason for not-ready, [] otherwise). Before any daemon it printed 'Daemon not running'.

## Links

### Out

- verifies [daemon-status-skips](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-skips.md)

## History

- `01M686RS47EK7F4H2WY15TNAQ5` rhiza daemon status names each accepted change set the daemon does not start and why
