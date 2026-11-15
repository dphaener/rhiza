# fast-suite-proposal-approved

decision · `01M7EQ7534CRTCC8BW89A04P8G`

## options

```json
[
  "targets as proposed, or tighter or looser",
  "budgets opt-in in CI's budgets job, or kept in make test",
  "projection rewrite in this milestone, or out of scope",
  "gate bug first, or the milestone first"
]
```

## quote

those  are all fine. continue without me

## quote_where

owner, 2026-11-15T02:44:54Z, reply to the proposal's four questions (targets; budgets opt-in; projection in or out; gate bug first)

## statement

The milestone is planned as the orchestrator proposed it on 2026-11-15: the targets as written (cold make test within 4 minutes, warm within 60 s on the 4-core VM, median PR CI run within 3 minutes, nothing skipped, no coverage lost); the read budgets move out of make test into an opt-in run that CI's budgets job runs; the projection rewrite of internal/api's FoldLog model stays out of scope, with only parallel tests and snapshot fixtures here; and the gate bug is fixed before the milestone starts. The five change sets run on the conductor default with the orchestrator arbitrating.

## Links

### Out

- derived-from [perf-milestone-interview](../source/perf-milestone-interview.md)

### In

- decided-by [milestone-38](../intent/milestone-38.md)
- decided-by [m38-cold](../intent/milestone-38/m38-cold.md)
- decided-by [m38-warm](../intent/milestone-38/m38-warm.md)
- decided-by [m38-ci](../intent/milestone-38/m38-ci.md)
- decided-by [m38-nothing-lost](../intent/milestone-38/m38-nothing-lost.md)
- decided-by [scaled-test-clock-env](../capability/rhiza/integrity/self-check/scaled-test-clock-env.md)
- decided-by [budget-tests-backdate-mtimes](../capability/rhiza/integrity/self-check/budget-tests-backdate-mtimes.md)
- decided-by [api-snapshot-fixtures](../capability/rhiza/integrity/self-check/api-snapshot-fixtures.md)
- decided-by [cheap-process-start](../capability/rhiza/integrity/self-check/cheap-process-start.md)
- decided-by [testfit-boundaries](../capability/rhiza/integrity/self-check/testfit-boundaries.md)
- decided-by [api-tests-parallel-session](../capability/rhiza/integrity/self-check/api-tests-parallel-session.md)
- decided-by [api-tests-parallel-runtime](../capability/rhiza/integrity/self-check/api-tests-parallel-runtime.md)
- decided-by [api-tests-parallel-rest](../capability/rhiza/integrity/self-check/api-tests-parallel-rest.md)
- decided-by [api-tests-parallel-model](../capability/rhiza/integrity/self-check/api-tests-parallel-model.md)
- decided-by [api-tests-parallel-finish](../capability/rhiza/integrity/self-check/api-tests-parallel-finish.md)

## History

- `01M7EQ19FNS1VDXT4750YX89RD` The test suite is fast. Owner, 2026-11-14/15, after the orchestrator measured the suite (cold make test 854 s wall for 972 CPU-s on 4 cores, warm 254 s): a cold make test within 4 minutes and a warm one within 60 s on the 4-core VM, the median PR CI run within 3 minutes, with no test skipped and no coverage lost; five change sets, the gate bug first, projection out of scope.
