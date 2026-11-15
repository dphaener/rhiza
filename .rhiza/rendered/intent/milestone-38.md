# milestone-38

intent · `01M7EQ74Z9CZ12B122F4WW6M6C`

## milestone

38

## statement

The test suite is fast: on the 4-core VM a cold make test finishes within 4 minutes and a warm one within 60 s, and the median CI run on a pull request finishes within 3 minutes, with no test skipped and no coverage lost (2026-11-15: cold 854 s wall for 972 CPU-s, warm 254 s, CI median 283 s).

## Links

### Out

- decided-by [fast-suite-proposal-approved](../decision/fast-suite-proposal-approved.md)
- decided-by [fast-suite-build-order](../decision/fast-suite-build-order.md)
- decided-by [m38-closed-cold-interim](../decision/m38-closed-cold-interim.md)
- decided-by [m38-ci-target-carried](../decision/m38-ci-target-carried.md)

### In

- serves change set `01M7FRV11EM44XJAJS6VXWC5DK` (Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets)

## History

- `01M7EQ19FNS1VDXT4750YX89RD` The test suite is fast. Owner, 2026-11-14/15, after the orchestrator measured the suite (cold make test 854 s wall for 972 CPU-s on 4 cores, warm 254 s): a cold make test within 4 minutes and a warm one within 60 s on the 4-core VM, the median PR CI run within 3 minutes, with no test skipped and no coverage lost; five change sets, the gate bug first, projection out of scope.
