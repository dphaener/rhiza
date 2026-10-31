# leaf-fix-exhaustion-spends-cycle

capability · `01M68585JZBX1K2RWAK29TVWSZ` · parent [fixer-exhaustion-spends-cycle](../fixer-exhaustion-spends-cycle.md)

## acceptance

A test where the leaf fixer fails all attempts of cycle 2: cycle 2's check outcome starts with cycle 2: fix failed: leaf leaf-a: fixer failed:, no check session starts in cycle 2, and cycle 3 starts a fresh first attempt; TestFixLeafNoCommitFailsCycle and TestFixLeafVerifierFailsCycle are updated to the new session lists

## statement

A leaf fixer or verifier failing every attempt of a cycle spends that cycle, and its checks fail with cycle N: fix failed: leaf <handle>: ...

## Links

### Out

- decided-by [leaf-fix-retry-fresh-attempt](../../decision/leaf-fix-retry-fresh-attempt.md)

### In

- verifies [leaf-fix-exhaustion-spends-cycle-verified](../../evidence/leaf-fix-exhaustion-spends-cycle-verified.md)
- verifies [leaf-fix-exhaustion-spends-cycle-qa-1](../../evidence/leaf-fix-exhaustion-spends-cycle-qa-1.md)
- verifies [pr-125-a-failed-review-leaf-fixer](../../evidence/pr-125-a-failed-review-leaf-fixer.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
