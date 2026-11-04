# usage-settings-qa-2

evidence · `01M6JM1FRRFZ9RCGW86M4CNQW9` · actor [claude-code-3977a5fc](../actor/claude-code-3977a5fc.md)

## result

PASS

## statement

QA cycle 2, built binary at 46cc1d58 in a scratch repo on the VM with the real ~/.claude transcripts (Rex faked, starts captured by RHIZA_TEST_DETACH_LOG with RHIZA_TEST_NOW set to the wall clock): rhiza daemon settings showed usage.session and usage.weekly unset and usage.reserve 20 (default); rhiza daemon set wrote session 1000, 5000000 and 1000000000000 and weekly 1000000000000, and daemon status listed them as git config.

## Links

### Out

- verifies [usage-settings](../capability/rhiza/orchestration/usage-governor/usage-settings.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
