# test-suite-locks-scripts-only

decision · `01M6JR8Z53XKKT41D34ZQY90B8`

## options

```json
[
  "lock cmd/rhiza's scripts and the timing budgets, the rest unlocked beside them",
  "lock the whole suite and the budget phase",
  "no lock, -trimpath only"
]
```

## quote

"i think you should monitor the daemon and adjust accordingly. we likely have some performance issues to tackle." (owner, 2026-11-03/04)

## statement

make test builds with -trimpath and takes one machine-wide flock (beside GOCACHE) around cmd/rhiza's scripts and, separately, around the timing budgets; the other packages run unlocked in the background beside the scripts. Measured on the loaded 4-core VM (load 15-30): the non-script packages are latency-bound, not CPU-bound (internal/index, slice, gen and op three at once took 235 s against 407 s one after another, at the same 86 CPU-s each; internal/api spends 500-590 s of wall on far less CPU), so locking them would only queue waits; cmd/rhiza's scripts use about 430 CPU-s per run in parallel, so they are what thrashes

## wording

the owner's words; the design is the orchestrator's, as the owner's proxy

## Links

### In

- decided-by [test-suite-shared](../capability/rhiza/integrity/self-check/test-suite-shared.md)

## History

- `01M6JJMK3M50P34PPX1QNTHYCB` The test suite is cheap to run in many worktrees at once: unchanged packages hit the shared test cache across worktrees, and concurrent full suites on one machine queue instead of thrashing
