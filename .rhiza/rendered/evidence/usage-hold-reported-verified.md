# usage-hold-reported-verified

evidence · `01M6JB087V50C7HTECCQZQWR54` · actor [claude-code-84a74884](../actor/claude-code-84a74884.md)

## result

PASS

## statement

At f9b2784e, go test ./cmd/rhiza -run 'TestScripts/(usage-hold|describe-golden)$' -count=1 passes. The new usage-hold.txtar seeds a daemon state that holds and shows a usage-reserve skip with the window/use/limit/reset detail in daemon status --json and text, a warn usage part in health --json and text, and exactly one usage-reserve inbox item (inbox open refuses it, as there is no tab to open). With ok.json (no hold), none of these shows: health's usage part is ok and the inbox is empty. describe-golden.txtar is regenerated (usage-reserve skip/inbox kinds, usage health part). Read the diff: State.UsageHold is saved per poll and keeps Since while the same window/reset holds, carried across restarts, and fed into daemon notifications. make check passed fmt-check, vet and test before self-check failed only on stale .rhiza/rendered files (run actor records, rendered when the change set lands). A second, loaded test run failed only the read-timing benchmark, which was contending for the machine.

## Links

### Out

- verifies [usage-hold-reported](../capability/rhiza/orchestration/usage-governor/usage-hold-reported.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
