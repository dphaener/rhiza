# ship-gate-base-env

capability · `01M78TAJC7RKC9QNK7T52NFT62` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShipLocalGate' shows a gate command that prints its environment receiving RHIZA_GATE_HEAD equal to the worktree's HEAD and RHIZA_GATE_BASE equal to git merge-base HEAD origin/main, and an empty RHIZA_GATE_BASE in a repo with no origin/main; go test -trimpath ./cmd/rhiza -run 'TestScripts/describe-golden' passes with the ship.gate summary naming both variables

## statement

rhiza ship gate runs the project's gate command with RHIZA_GATE_BASE (the merge-base of HEAD with origin/main) and RHIZA_GATE_HEAD (the commit it checks) in its environment, RHIZA_GATE_BASE empty when no merge-base can be found, so a project's own gate command can check only what the branch changes while Rhiza stays language-neutral

## Links

### Out

- depends-on [ship-local-gate](ship-local-gate.md)
- decided-by [local-gate-reuse-and-scope](../../../../../decision/local-gate-reuse-and-scope.md)

### In

- depends-on [rhiza-scoped-gate](../../../integrity/self-check/rhiza-scoped-gate.md)
- verifies [ship-gate-base-env-verified](../../../../../evidence/ship-gate-base-env-verified.md)
- verifies [ship-gate-base-env-qa-1](../../../../../evidence/ship-gate-base-env-qa-1.md)
- verifies [pr-513-a-passing-local-gate-is-reused](../../../../../evidence/pr-513-a-passing-local-gate-is-reused.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
