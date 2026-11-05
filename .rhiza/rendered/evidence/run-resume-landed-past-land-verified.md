# run-resume-landed-past-land-verified

evidence · `01M6JJDEPVWWAS8K4FTYWT1YZE` · actor [claude-code-4482b6f9](../actor/claude-code-4482b6f9.md)

## result

PASS

## statement

At c092ae36 ran go test ./cmd/rhiza -run 'TestScripts/(run-resume-landed-past-land|run-resume-landed-plan|run-resume-unaccepted)$' -count=1: all three PASS. run-resume-landed-past-land.txtar runs a ship-change-shaped run (journal workflow rewritten to ship-change) that enters land and stops at its wait-main, lands the change set, then rhiza run resume goes on past the wait and prints succeeded. With internal/api/run.go reverted to HEAD~1 in a scratch worktree the new script fails: resume refused with exit 3, 'is landed, no longer accepted'. run-resume-landed-plan.txtar (BUILDLAND01 ship-change, journal never entered land: refused, names rhiza run abandon, cmp journal unchanged) and run-resume-unaccepted.txtar (withdrawn: refused same way, journal unchanged) are unmodified and pass. Diff read: the gate lets a non-plan run on a landed change set through only when st.Visits[land] > 0 (counted on entered events, so land entered or passed). make check: fmt-check, vet and go test ./... passed; self-check failed only on rendering this verifier's own untracked actor record, not on the leaf's commit.

## Links

### Out

- verifies [run-resume-landed-past-land](../capability/rhiza/workflows/workflow-runs/run-resume-landed-past-land.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
