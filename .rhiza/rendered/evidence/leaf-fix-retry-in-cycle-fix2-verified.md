# leaf-fix-retry-in-cycle-fix2-verified

evidence · `01M68A91SYRH08DQXZAAHS1426` · actor [claude-code-3d36b1c7](../actor/claude-code-3d36b1c7.md)

## result

PASS

## statement

Verifier of fix attempt 3 at 544cb11 (cycle 2's confirmed model f1: two leaf fixers live at once after a verifier retry). Read the diff: fixActions now starts with fixerFree false when any leaf with findings has a started, unjudged fixer of the cycle (liveFixer over leafFixAttempt tries), so a leaf whose verifier failed waits for the next leaf's live fixer. Its retry then starts first in leaf order once the slot is free. Two new TestFixActions cases cover both: retry waits (only 7:leaf-b:fixer polled), then takes the slot (6:leaf-a:fixer plus leaf-b's verifier). Ran go test ./internal/run -run 'Fix|Leaf' -count=1 -v: all pass, including TestFixLeafFixerRetriedInCycle and TestFixLeafVerifierRetriedInCycle (the acceptance tests) and TestFixActions. make fmt-check vet test pass. make self-check fails only on stale .rhiza/rendered actor/index files from the run's own record commits, not on code. Model f2 (the fixActions doc still says 'both on attempt cycle') is a minor finding and remains.

## Links

### Out

- verifies [leaf-fix-retry-in-cycle](../capability/fixer-retry-in-cycle/leaf-fix-retry-in-cycle.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
