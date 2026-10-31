# fixer-retries-within-cycle-v2

decision · `01M6899M1XNCHCSWP9HTM0R7GT`

## reasoning

Restates fixer-retries-within-cycle with its change-set fixer part unchanged and replaces its leaf sentence ('A leaf's fixer or verifier is not retried'), which change set 01M6580HBAKHN7XTR7XQV3XR0Q reversed: internal/run now retries a failed leaf fix within the cycle with the change-set fixer's budget, per leaf-fix-retry-fresh-attempt.

## statement

A failed change-set review fixer is retried within its cycle up to MaxFixCycles more times (the same cap as the checks' further cycles). Only when those retries are exhausted does the cycle count as spent, and the checks of that cycle fail with the fix-failed detail. A failed leaf fixer or verifier is retried the same way, as a fresh fixer and verifier attempt (leaf-fix-retry-fresh-attempt); it no longer fails the phase on its first failure.

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- relates-to [leaf-fix-retry-fresh-attempt](leaf-fix-retry-fresh-attempt.md) — The leaf part of this Decision is leaf-fix-retry-fresh-attempt
- supersedes [fixer-retries-within-cycle](fixer-retries-within-cycle.md)

### In

- decided-by [fixer-retry-in-cycle](../capability/fixer-retry-in-cycle.md)
- decided-by [fixer-exhaustion-spends-cycle](../capability/fixer-exhaustion-spends-cycle.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
