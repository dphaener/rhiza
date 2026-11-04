# usage-hold-reported-qa-3

evidence · `01M6JRP2M9S7NXN2DXNX5NQ846` · actor [claude-code-c4e7a246](../actor/claude-code-c4e7a246.md)

## result

PASS

## statement

QA cycle 3, binary built at 1f1a3b35 in a scratch repo on the VM with the real ~/.claude transcripts (Rex faked, RHIZA_TEST_NOW at the wall clock, load ~24): with usage.session 1000, daemon --once saved usage_hold (session, 4583485 used, limit 1000, 20% reserve) and the windows with their used. rhiza health showed warn usage with that detail and inbox one usage-reserve item; health took 0.6 s real / 0.1 s user, and strace -e openat of it opened no ~/.claude/projects file. After raising the limit, the next poll cleared the hold and opened a new session window after the reset; health (text and --json) showed 'ok usage session 1721 of 1000000000000; weekly 18945650 of 1000000000000', exactly the used the poll saved, in under 1 s, and the inbox said Nothing waits on you. With the saved session window's reset edited into the past, health showed session 0.

## Links

### Out

- verifies [usage-hold-reported](../capability/rhiza/orchestration/usage-governor/usage-hold-reported.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
