# perf-milestone-interview

source · `01M7EQ1R1FB9Z9TGZW6575TKF2`

## quotes

```json
[
  {
    "by": "owner",
    "text": "We will then work through test suite performance, it is terrible.",
    "where": "2026-11-14T14:33Z, opening request"
  },
  {
    "by": "owner",
    "text": "This should all still go through the Rhiza process. I want to get this all straightened out before we start the daemon up again.",
    "where": "2026-11-14T14:33Z, opening request"
  },
  {
    "by": "owner",
    "text": "We really need to focus on Rhiz performance. It's TERRIBLE.",
    "where": "2026-11-14"
  },
  {
    "by": "owner",
    "text": "I only have 4 cores available. I can give the VM 24GB at most",
    "where": "2026-11-14"
  },
  {
    "by": "owner",
    "text": "those  are all fine. continue without me",
    "where": "2026-11-15T02:44:54Z, reply to the proposal and its four questions (targets, budgets opt-in, projection in or out, gate bug first)"
  }
]
```

## raw_sha256

9006273335dd062abf881a9266591a9d5e96be3fc0347a29374528f26f88f9c1

## summary

Chat between the owner and the orchestrator, 2026-11-14/15. The owner asked to work through test suite performance through the Rhiza process before the daemon starts again, on a VM of 4 cores. The orchestrator measured the suite: a cold make test takes 854 s of wall time for 972 CPU-s across 4 cores (about 28% use), a warm one 254 s; compiling is 70 s cold and 7 s warm. Causes found: testscripts sleep through the builtin workflows' real 30 s polls and 2 min backoff (run-ship-ci-fix: 165 s wall, 18 s CPU); the read budgets run every time, alone, at the end (188 s even warm), and TestReadBudget mostly sleeps 2.2 s to wait out the index's 2 s racy window; internal/api's 749 tests run serially (no t.Parallel in the repo) and 58% of their CPU is store.ReadLog, since each write re-reads the log; scripts use about 1.2 cores at -parallel 4 (381 s instead of 639 s at -parallel 12, no failures); CI script shards are round-robin by name and took 61-239 s; each rhiza process start costs 39 ms, 22 ms of it go-runewidth's eager init, across 6,600+ starts; internal/slice TestFit sweeps every token budget (about 40 s). The orchestrator proposed the milestone 'The test suite is fast' (cold make test within 4 minutes and warm within 60 s on the 4-core VM, median CI run within 3 minutes, nothing skipped, no coverage lost) with five change sets in order of payoff: a virtual clock for scripts; read budgets as opt-in benchmarks with mtime backdating instead of sleeps; parallel scripts, longest first, and CI shards weighted by recorded timings; internal/api tests in parallel with snapshot fixtures; cheap process start and TestFit at its boundaries. It asked four questions: the targets; budgets moving out of make test so CI runs them only in its budgets job; whether projection replacing the FoldLog model is in this milestone or out of scope; and the gate bug first, then the milestone. The owner approved all four as proposed and left the rest to the orchestrator. Nothing was left open.

## title

Planning the test-perf milestone with the owner

## Links

### In

- derived-from [fast-suite-proposal-approved](../decision/fast-suite-proposal-approved.md)
- derived-from [fast-suite-measuring-commands](../decision/fast-suite-measuring-commands.md)
- derived-from [fast-suite-build-order](../decision/fast-suite-build-order.md)
- derived-from [scaled-test-clock](../decision/scaled-test-clock.md)
- derived-from [budgets-opt-in](../decision/budgets-opt-in.md)
- derived-from [api-fixture-templates](../decision/api-fixture-templates.md)

## History

- `01M7EQ19FNS1VDXT4750YX89RD` The test suite is fast. Owner, 2026-11-14/15, after the orchestrator measured the suite (cold make test 854 s wall for 972 CPU-s on 4 cores, warm 254 s): a cold make test within 4 minutes and a warm one within 60 s on the 4-core VM, the median PR CI run within 3 minutes, with no test skipped and no coverage lost; five change sets, the gate bug first, projection out of scope.
