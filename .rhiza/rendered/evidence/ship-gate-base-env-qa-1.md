# ship-gate-base-env-qa-1

evidence · `01M792XMY89CNZ48YA5GRMVHJP` · actor [claude-code-7ac56481](../actor/claude-code-7ac56481.md)

## result

PASS

## statement

QA with the head build in the same scratch repo, the gate echoing HEAD=$RHIZA_GATE_HEAD BASE=$RHIZA_GATE_BASE: rhiza ship gate --json output was HEAD=aa1864b... (git rev-parse HEAD) BASE=b8d5ae1... (git merge-base HEAD origin/main); after removing the origin remote and origin/main, output was HEAD=ae5ef44... BASE= (empty) and the gate still passed. The ship.gate summary in docs/reference/cli.md and describe-golden names both variables.

## Links

### Out

- verifies [ship-gate-base-env](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-base-env.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
