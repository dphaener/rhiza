# leaf-fix-retry-in-cycle

capability · `01M68584KCHFKS2RDF7M3E8XFG` · parent [fixer-retry-in-cycle](../fixer-retry-in-cycle.md)

## acceptance

New tests in internal/run/fixleaf_test.go with a fake runner as in fixcs_test.go: the leaf fixer fails attempt 1 and passes attempt 2 of cycle 2 and the step passes in 2 cycles with both attempts' sessions recorded apart and the second's judgement before cycle 2's check session; the same with the verifier failing first

## statement

A failed review leaf fixer, or its verifier, is retried within its cycle as a fresh fixer and verifier attempt, up to MaxFixCycles more times, and the cycle's checks start only once every leaf fix and the change-set fixer pass

## Links

### Out

- decided-by [leaf-fix-retry-fresh-attempt](../../decision/leaf-fix-retry-fresh-attempt.md)

### In

- verifies [leaf-fix-retry-in-cycle-verified](../../evidence/leaf-fix-retry-in-cycle-verified.md)
- verifies [leaf-fix-retry-in-cycle-qa-1](../../evidence/leaf-fix-retry-in-cycle-qa-1.md)
- verifies [leaf-fix-retry-in-cycle-fix-verified](../../evidence/leaf-fix-retry-in-cycle-fix-verified.md)
- verifies [leaf-fix-retry-in-cycle-fix2-verified](../../evidence/leaf-fix-retry-in-cycle-fix2-verified.md)
- verifies [pr-125-a-failed-review-leaf-fixer](../../evidence/pr-125-a-failed-review-leaf-fixer.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
