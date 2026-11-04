# usage-hold-qa-2

evidence · `01M6JM2DDQV2YE9HHWH5QXNRWS` · actor [claude-code-3977a5fc](../actor/claude-code-3977a5fc.md)

## result

PASS

## statement

QA cycle 2, built binary at 46cc1d58 in a scratch repo on the VM with the real ~/.claude transcripts (Rex faked, starts captured by RHIZA_TEST_DETACH_LOG with RHIZA_TEST_NOW set to the wall clock): daemon --once with usage.session 1000 (and 5000000) while the session window had ~4.8M used held the ready change set as a usage-reserve skip and started nothing; a repeat poll kept the hold's since; with both limits set high the hold cleared and the same poll started the change set (rhiza run start --foreground captured); after the window reset at 06:29:50 the next poll started the held change set.

## Links

### Out

- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
