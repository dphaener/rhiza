# usage-windows-qa-2

evidence · `01M6JM1Z56QT2D4M2MY4P17C6T` · actor [claude-code-3977a5fc](../actor/claude-code-3977a5fc.md)

## result

PASS

## statement

QA cycle 2, built binary at 46cc1d58 in a scratch repo on the VM with the real ~/.claude transcripts (Rex faked, starts captured by RHIZA_TEST_DETACH_LOG with RHIZA_TEST_NOW set to the wall clock): daemon state saved a session window 01:29:50.689Z-06:29:50.689Z and a weekly window; a second poll kept both unchanged (no re-anchoring); the first poll after the 06:29:50 reset opened a new session window at 06:30:01.758Z (the first record after the reset) with 4736 tokens used, weekly unchanged; limit hits with resets 3h and 4d out set the session window to end at the 3h reset and the weekly at the 4d reset.

## Links

### Out

- verifies [usage-windows](../capability/rhiza/orchestration/usage-governor/usage-windows.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
