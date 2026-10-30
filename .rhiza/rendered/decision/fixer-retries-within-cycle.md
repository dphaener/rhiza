# fixer-retries-within-cycle

decision · `01M65054XFRT6Z9CWSFS9MYH2R`

## reasoning

The change set intent says a failed fixer must not use up a review fix cycle, so the retry happens inside the cycle and the cycle's checks re-run only after its fixers pass. The intent sets no retry count, and the code today (internal/run/checks.go nextCheckActions, fixActions failure path) spends the cycle on the first failed fix. Retries must be bounded or a broken fixer loops forever, so I reused MaxFixCycles (2) as the in-cycle retry budget rather than adding a new constant: one knob, same meaning as the existing cap, and the existing tests that count cycles against MaxFixCycles stay valid. Unbounded retries and a zero-retry answer were rejected: the first never terminates, the second is today's behaviour the intent asks to change.

## statement

A failed change-set review fixer is retried within its cycle up to MaxFixCycles more times (the same cap as the checks' further cycles). Only when those retries are exhausted does the cycle count as spent, and the checks of that cycle fail with the fix-failed detail. A leaf's fixer or verifier is not retried: its failure still fails the phase and spends the cycle (tracked in discovered-from draft change set 01M6580HB).

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-33462c8e-1](../question/question-33462c8e-1.md)

### In

- decided-by [fixer-retry-in-cycle](../capability/fixer-retry-in-cycle.md)
- decided-by [fixer-exhaustion-spends-cycle](../capability/fixer-exhaustion-spends-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
