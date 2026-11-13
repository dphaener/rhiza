# question-aea6d8fc-1

question · `01M78T6GS9KPCS99BXNACD749P`

## session

aea6d8fc-7b79-4c73-8624-01a1d1a09694

## statement

Which landed intent should change set 01M78T0 (A passing local gate is reused) serve? It saves ~20 min of ship-change wall time per run on gate days; the perf report says milestone 35 has no change set covering it. Its sibling 01M78RM (gate switches on only for repeated named flakes) serves m4-runs-hold-up. / What does a passing local gate get reused on: which tree must be unchanged for ship gate to pass at once without rerunning it? / How does the gate check only what the change touches while the product stays language-neutral?

## Links

### In

- answers [local-gate-reuse-and-scope](../decision/local-gate-reuse-and-scope.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
