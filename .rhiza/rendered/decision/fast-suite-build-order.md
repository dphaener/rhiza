# fast-suite-build-order

decision · `01M7EQ75C7E0GZVKHFSNEK2FAE`

## reasoning

The proposal ordered the change sets by payoff; a depends-on link is real only where one needs another's output, and a timing file recorded before the clock change would weight the shards by sleeps that no longer happen.

## statement

The five change sets build in the proposal's order of payoff: virtual clock, opt-in budgets, parallel weighted scripts, parallel internal/api tests, cheap start and boundary TestFit; only the weighted shards depend on another, the virtual clock, because their timing file must record script times after the clock change

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [perf-milestone-interview](../source/perf-milestone-interview.md)

### In

- decided-by [milestone-38](../intent/milestone-38.md)

## History

- `01M7EQ19FNS1VDXT4750YX89RD` The test suite is fast. Owner, 2026-11-14/15, after the orchestrator measured the suite (cold make test 854 s wall for 972 CPU-s on 4 cores, warm 254 s): a cold make test within 4 minutes and a warm one within 60 s on the 4-core VM, the median PR CI run within 3 minutes, with no test skipped and no coverage lost; five change sets, the gate bug first, projection out of scope.
