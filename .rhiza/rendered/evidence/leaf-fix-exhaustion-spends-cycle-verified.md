# leaf-fix-exhaustion-spends-cycle-verified

evidence · `01M688PC3QSBVWBC22HTWCSXY5` · actor [claude-code-610b3e92](../actor/claude-code-610b3e92.md)

## result

PASS

## statement

At 1b5ba13, read the diff: TestFixLeafExhaustionSpendsCycle (internal/run/fixleaf_test.go) fails leaf-a's fixer on every attempt of cycle 2 (leafFixAttempt(2,1..1+MaxFixCycles) = 2,3,4) and asserts alpha's cycle 2 outcome is CheckFailed with detail prefix 'cycle 2: fix failed: leaf leaf-a: fixer failed:', check sessions are only alpha/beta/gamma check 1, alpha critic 1, alpha check 3 (none in cycle 2), and leaf sessions are fixer 2,3,4 then fixer 5 + verifier 5, i.e. cycle 3's fresh first attempt (leafFixAttempt(3,1)=5). TestFixLeafNoCommitFailsCycle now expects fixers 2-7 and TestFixLeafVerifierFailsCycle expects fixer+verifier 2-7, the new per-attempt session lists. Ran go test ./internal/run -run TestFixLeaf -count=1 -v: all 7 TestFixLeaf* pass. make check: fmt, vet, and tests pass apart from TestGeneratedModel in internal/tui/modelapp, which this change set does not touch; it took 41s under load and passes when run alone (4.8s). The self-check fails only on stale .rhiza/rendered files for this run's new actors, which is expected mid-run.

## Links

### Out

- verifies [leaf-fix-exhaustion-spends-cycle](../capability/fixer-exhaustion-spends-cycle/leaf-fix-exhaustion-spends-cycle.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
