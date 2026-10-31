# leaf-fix-retry-in-cycle-verified

evidence · `01M688P65F5FAASR85QGXA6B2J` · actor [claude-code-3541095c](../actor/claude-code-3541095c.md)

## result

PASS

## statement

At f44c812 ran go test ./internal/run -run 'TestFixLeaf|TestFixActions|TestNextCheckActionsFixPhase|TestFixChangeSet' -count=1 -v: all pass, including new TestFixLeafFixerRetriedInCycle (fake runner; fixer fails attempt 2, passes retry attempt 3; review passes in 2 cycles; 3 sessions and 3 attempts recorded apart, judgements fail then pass; 'leaf-attempt leaf-a verifier 3 pass' journalled before 'check-session alpha check 2') and TestFixLeafVerifierRetriedInCycle (verifier fails attempt 2; fresh fixer and verifier on attempt 3; same ordering check). Read the diff: fixActions loops 1+MaxFixCycles tries per leaf, keyed by leafFixAttempt(cycle,try), failing the phase only after the last try; startLeaf now takes the cycle from leafAction.cycle. make check exits 0.

## Links

### Out

- verifies [leaf-fix-retry-in-cycle](../capability/fixer-retry-in-cycle/leaf-fix-retry-in-cycle.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
