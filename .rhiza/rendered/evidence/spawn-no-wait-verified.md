# spawn-no-wait-verified

evidence · `01M7QEA6GTNC9KDRAT01J9YX1B` · actor [claude-code-587ab5c7](../actor/claude-code-587ab5c7.md)

## result

PASS

## statement

Verified commit 9cd172821 by re-running its Checks: block, all exit 0. go test -count=1 -trimpath -p 2 ./internal/api -run TestRunSpawn passes, including the new NoWait subtest (session child-1 and brief set, no outcome or detail, ExitCode 0, NoWait response, one spawn entry with no_wait, no spawn-ended entry, no rex calls so the pane stays open) and NoWaitStartFails (start error returned, no spawn entry journaled). The existing waiting cases now also assert !sp.NoWait. go test ./internal/run -run 'TestSpawn|TestFold' passes, and TestFoldSpawnNoWait shows the mark folded onto State.Spawns. go test ./cmd/rhiza -run TestScripts/describe-golden -update left no diff. TestScripts/(run-spawn|describe-golden)$ passes, and run-spawn checks a --no-wait spawn: exit 0, no_wait true, no outcome or detail, exactly one journaled spawn entry with no_wait. rhiza ship generate and rhiza render exit 0. In the diff, the --no-wait return comes after the hold wait and the start, and right after the spawn entry is journaled, so the waiting path is unchanged. No stubs or TODOs, and no messages sent. internal/run/journal_test.go is outside the leaf's files list but is the fold test the acceptance asks for.

## Links

### Out

- verifies [spawn-no-wait](../capability/rhiza/workflows/workflow-runs/spawn-no-wait.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)
