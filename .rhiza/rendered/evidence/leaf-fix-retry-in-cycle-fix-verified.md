# leaf-fix-retry-in-cycle-fix-verified

evidence · `01M689F85M4VD5T36GXZTBBPRN` · actor [claude-code-e49bc6ac](../actor/claude-code-e49bc6ac.md)

## result

PASS

## statement

At a4a8c87 (the review fix 'Journal the review cycle on a leaf's fix sessions'), I re-ran the acceptance with go test ./internal/run -run TestFixLeaf -count=1. All 7 pass. That includes TestFixLeafFixerRetriedInCycle (fixer 2 fails, fixer and verifier at leafFixAttempt(2,2)=3 pass, the step passes in 2 cycles, 3 sessions and 3 judgements kept apart, now each journalled with Cycle 2 and CycleAttempts [2 3], and verifier 3's pass comes before cycle 2's check-session) and TestFixLeafVerifierRetriedInCycle (verifier 2 fails, fixer and verifier 3 pass, same ordering). go test ./internal/api -run TestReviewSummary passes, with cycle 3's two fixer tries rendered on separate lines. make fmt-check, vet and test (with describe-golden carrying the new cycle fields and corrected attempt descriptions) all pass. From reading the diff: the confirmed quality f1 (attempt documented as the cycle) and model f1 (the summary keyed by cycle) are fixed. leaf-session entries carry Cycle from startLeaf, RunLeafSession and RunLeafAttempt expose cycle, run status text shows 'in cycle N', fixStatus reports s.Cycle, and summaryCycles groups by the session's cycle. make self-check fails only on stale .rhiza/rendered actor and index files left by the run's chore record commits, which this commit does not touch.

## Links

### Out

- verifies [leaf-fix-retry-in-cycle](../capability/fixer-retry-in-cycle/leaf-fix-retry-in-cycle.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
