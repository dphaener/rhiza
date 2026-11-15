# fast-suite-measuring-commands

decision · `01M7EQ757PJTZR13E5VE26D63E`

## reasoning

The owner approved the targets but not how to check them; a done point needs one runnable measure, and a fixed commit gives the test count a baseline that does not move with main.

## statement

Each done point is measured by the command its statement names, on the 4-core VM with no other test run going, against the 2026-11-15 baseline in /home/rhiza/tmp/claude-1000/perf-data; the test count baseline is commit 88516d956, the main this plan was made on

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [perf-milestone-interview](../source/perf-milestone-interview.md)

### In

- decided-by [m38-cold](../intent/milestone-38/m38-cold.md)
- decided-by [m38-warm](../intent/milestone-38/m38-warm.md)
- decided-by [m38-ci](../intent/milestone-38/m38-ci.md)
- decided-by [m38-nothing-lost](../intent/milestone-38/m38-nothing-lost.md)

## History

- `01M7EQ19FNS1VDXT4750YX89RD` The test suite is fast. Owner, 2026-11-14/15, after the orchestrator measured the suite (cold make test 854 s wall for 972 CPU-s on 4 cores, warm 254 s): a cold make test within 4 minutes and a warm one within 60 s on the 4-core VM, the median PR CI run within 3 minutes, with no test skipped and no coverage lost; five change sets, the gate bug first, projection out of scope.
