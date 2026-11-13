# pr-513-a-passing-local-gate-is-reused

evidence · `01M7942JJAVG9Y4GVR84QD2G2X` · actor [dphaener](../actor/dphaener.md)

## head

7a96ac202aae8f55b2f0adfcd46d0e0c70574172

## pr

https://github.com/dphaener/rhiza/pull/513

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #513 (https://github.com/dphaener/rhiza/pull/513) passed CI at 7a96ac2

## status

verified

## Links

### Out

- verifies [ship-gate-reuse-code-tree](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-reuse-code-tree.md)
- verifies [ship-gate-base-env](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-base-env.md)
- verifies [rhiza-scoped-gate](../capability/rhiza/integrity/self-check/rhiza-scoped-gate.md)
- verifies [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
