# usage-hold-qa-3

evidence · `01M6JRQ66TEA81CNCQCS1G10WF` · actor [claude-code-c4e7a246](../actor/claude-code-c4e7a246.md)

## result

PASS

## statement

QA cycle 3, binary built at 1f1a3b35 in a scratch repo on the VM with the real ~/.claude transcripts (Rex faked, RHIZA_TEST_NOW at the wall clock, load ~24): daemon --once with usage.session 1000 over a session window of 4583485 tokens held (usage_hold saved, no start in the detach log); with usage.session and .weekly set to 1e12 the next poll, after the window's reset at 07:44:29Z, cleared the hold and saved a new session window starting at the first record after the reset (07:44:48Z). go test ./internal/api -run 'TestUsage|TestHealth' and the usage-* scripts pass.

## Links

### Out

- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
