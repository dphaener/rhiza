# spawn-no-wait

capability · `01M7QE202T3TYENJVMH40K1XPT` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/api -run TestRunSpawn passes with a new case: a --no-wait spawn returns while its child still runs, with the session set, no outcome or detail, exit code 0, a spawn entry marked no-wait and no spawn-ended entry journaled, and the child's pane left open; the existing waiting cases still pass unchanged (their spawn entries carry no no-wait mark); (failure case) a --no-wait spawn whose child cannot be started returns the start error and journals no spawn entry. go test -trimpath ./internal/run -run 'TestSpawn|TestFold' shows the no-wait mark folded onto State.Spawns. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-spawn|describe-golden)' passes, describe-golden regenerated with -update.

## files

internal/api/run_spawn.go, internal/api/run_spawn_test.go, internal/run/journal.go, cmd/rhiza/testdata/script/run-spawn.txtar, cmd/rhiza/testdata/script/describe-golden.txtar, docs/reference/cli.md

## must_not

Change the default (waiting) behaviour, output or exit codes of rhiza run spawn; journal a spawn-ended entry or close the child's pane from a --no-wait spawn; send any message

## statement

rhiza run spawn --no-wait returns once the child session is started and its spawn entry is journaled, the entry marked no-wait (a new journal field, folded onto run.Spawn), its response carrying the child's session and brief and no outcome, exit 0; the hold wait before the start is unchanged, and without --no-wait spawn waits and prints the outcome exactly as today.

## Links

### Out

- decided-by [spawn-no-wait-wake](../../../../decision/spawn-no-wait-wake.md)

### In

- depends-on [spawn-wake-conductor](spawn-wake-conductor.md)
- verifies [spawn-no-wait-verified](../../../../evidence/spawn-no-wait-verified.md)
- verifies [spawn-no-wait-qa-1](../../../../evidence/spawn-no-wait-qa-1.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)
