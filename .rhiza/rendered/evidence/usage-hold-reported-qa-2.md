# usage-hold-reported-qa-2

evidence · `01M6JM2MZ85JNPFZYCFRJ83EJJ` · actor [claude-code-3977a5fc](../actor/claude-code-3977a5fc.md)

## result

PASS

## statement

QA cycle 2, built binary at 46cc1d58 in a scratch repo on the VM with the real ~/.claude transcripts (Rex faked, starts captured by RHIZA_TEST_DETACH_LOG with RHIZA_TEST_NOW set to the wall clock): while held, rhiza daemon status listed 'usage reserve: held until the window resets (session window: 4765656 tokens used, ... over the limit of 1000 less a 20% reserve; resets 2026-11-04T06:29:50Z)', rhiza inbox (text and --json) one usage-reserve item, rhiza health (text and --json) a warn usage part with the same detail; after the hold cleared inbox said Nothing waits on you and health showed 'ok usage session 4736 of 5000000; weekly 18338538, no limit known'.

## Links

### Out

- verifies [usage-hold-reported](../capability/rhiza/orchestration/usage-governor/usage-hold-reported.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
