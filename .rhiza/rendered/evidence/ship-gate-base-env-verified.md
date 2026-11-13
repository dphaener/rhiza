# ship-gate-base-env-verified

evidence · `01M792A50H19JDRS9MR876TYEB` · actor [claude-code-f005f36a](../actor/claude-code-f005f36a.md)

## result

PASS

## statement

At 3927f7838: go test -trimpath ./internal/api -run TestShipLocalGate passes, including subtest 'gate environment names the head and the merge-base', where a gate command echoing env:$RHIZA_GATE_HEAD:$RHIZA_GATE_BASE: prints the worktree's rev-parse HEAD and the git merge-base HEAD origin/main answer, and prints an empty base when merge-base fails (no origin/main). go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes, and the golden ship.gate summary names RHIZA_GATE_HEAD and RHIZA_GATE_BASE. Diff read: shipLocalGate appends both variables to os.Environ(); the real Exec (cmd/rhiza execIn) returns stdout only, so a failed merge-base yields an empty RHIZA_GATE_BASE, not error text.

## Links

### Out

- verifies [ship-gate-base-env](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-base-env.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
