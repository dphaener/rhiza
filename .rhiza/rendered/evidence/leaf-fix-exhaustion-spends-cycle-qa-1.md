# leaf-fix-exhaustion-spends-cycle-qa-1

evidence · `01M688W72PYDMFZ4QVS3Q5KVYJ` · actor [claude-code-96da0913](../actor/claude-code-96da0913.md)

## result

PASS

## statement

Same built binary and scratch setup as leaf-fix-retry-in-cycle-qa-1. Fixer exits 1 on every session of cycle 2: three fixers ran in cycle 2 (journal leaf-attempt fixer 2, 3, 4 all fail), no check session started in cycle 2, quality's cycle 2 was journalled failed with detail 'cycle 2: fix failed: leaf alpha: fixer failed: exit 1', and cycle 3 started a fresh first attempt (fixer 5 pass, verifier 5 pass) before quality check 3, which passed; run reached sign-off. Verifier exits 1 on every session of cycle 2: three fixer+verifier pairs ran in cycle 2, cycle 2 failed with 'cycle 2: fix failed: leaf alpha: verifier failed: exit 1', cycle 3 ran a fresh pair then quality check 3. Note: rhiza run status text shows only 'cycle 2: failed', without the detail.

## Links

### Out

- verifies [leaf-fix-exhaustion-spends-cycle](../capability/fixer-exhaustion-spends-cycle/leaf-fix-exhaustion-spends-cycle.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
