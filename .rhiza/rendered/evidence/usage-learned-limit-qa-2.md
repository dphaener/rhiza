# usage-learned-limit-qa-2

evidence · `01M6JM26KS1X32C8GWW80Z07QE` · actor [claude-code-3977a5fc](../actor/claude-code-3977a5fc.md)

## result

PASS

## statement

QA cycle 2, built binary at 46cc1d58 in a scratch repo on the VM with the real ~/.claude transcripts (Rex faked, starts captured by RHIZA_TEST_DETACH_LOG with RHIZA_TEST_NOW set to the wall clock): with usage.session unset, .git/rhiza/usage.limits.jsonl hits with a reset 3h out, a reset 4d out and no reset taught exactly two learned limits (session 1872595, weekly 18273795, the window use at the hit), and daemon --once held on the learned session limit; setting usage.session 1000000000000 then held on the learned weekly limit instead, so the setting won.

## Links

### Out

- verifies [usage-learned-limit](../capability/rhiza/orchestration/usage-governor/usage-learned-limit.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
