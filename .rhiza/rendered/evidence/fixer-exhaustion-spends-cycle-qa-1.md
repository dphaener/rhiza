# fixer-exhaustion-spends-cycle-qa-1

evidence · `01M6581P3PY3JV6PXJTFPQCMS8` · actor [claude-code-32688ba6](../actor/claude-code-32688ba6.md)

## result

PASS

## statement

At e499bf2, with the built binary and the same scratch repo, fake rex and fake claude as fixer-retry-in-cycle-qa-1, I ran rhiza run start twice. (1) Fixer fails its first 3 calls, then commits: cycle 2 ran three fixer sessions (1 + MaxFixCycles), all fail exit 1, and no check session; run status --json gave alpha's cycle 2 the detail 'cycle 2: fix failed: change set: fixer failed: exit 1'; cycle 3 started a fresh fixer, which passed, then alpha's check, and the run passed with alpha passed after 3 cycles. (2) Fixer always fails: cycles 2 and 3 each ran three failing fixer sessions and no check session; no cycle 4 ran; the step failed with '1 of 2 checks passed after 3 cycles; failed: alpha' and the run exited 1. So a spent cycle counts against MaxFixCycles=2.

## Links

### Out

- verifies [fixer-exhaustion-spends-cycle](../capability/fixer-exhaustion-spends-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
