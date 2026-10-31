# daemon-skip-reasons-verified

evidence · `01M694DHQRR616F0J1B68HTD1Z` · actor [claude-code-0c403fe6](../actor/claude-code-0c403fe6.md)

## result

PASS

## statement

At ef244d2: go test ./internal/daemon -run TestSkips passes; the table test of the pure Skips (beside Pick, sharing one decide function) covers live (running and interrupted runs), failed-not-cleared, held and no-slot, and shows an abandoned failure is picked when a slot is free and a retried failure is never skipped as failed (only as live, its retry running); it also checks Pick and Skips split the candidates. go test ./internal/api -run TestDaemonPicks passes, asserting daemonCandidates reports no-milestone (backlog), not-ready and no-leaves skips. Diff read: State gains Skips saved each successful poll, with unjournaled skips; Pick's picks are unchanged. fmt-check, vet and go test ./... pass; make self-check fails only on stale render of actor entries written by this run's sessions after the commit, not by the leaf.

## Links

### Out

- verifies [daemon-skip-reasons](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skip-reasons.md)

## History

- `01M686RS47EK7F4H2WY15TNAQ5` rhiza daemon status names each accepted change set the daemon does not start and why
