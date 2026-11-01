# plan-guide-new-change-set-verified

evidence · `01M6A8G2BNYM7X98DZY03NBJPK` · actor [claude-code-5265da2c](../actor/claude-code-5265da2c.md)

## result

PASS

## statement

At 125c2c9c, `go test ./internal/guide -run TestPlanText -count=1` passes; its plan-no-intent-yet case asserts plan.md covers an intention change set that serves no intent yet (detected by `"serves": []` in `rhiza change show --json`), picking a landed intent or writing a new `"kind":"intent"`, and writing the serves link with `rhiza link "$RHIZA_CHANGE" serves <intent>` before any leaf. Read the diff of all three Leaf commits: plan.md step 2 says exactly that. The guide's read is real: change show --json now carries a serves field (LinksFrom over the change set's own live serves links), and `go test ./cmd/rhiza -run 'TestScripts/(change-show-serves|describe-golden)'` passes, showing `"serves": []` on a fresh draft and the intent ID after linking. `rhiza check` on a clean clone at the commit prints ok; make check in the run checkout only flags render of this verifier's own actor record.

## Links

### Out

- verifies [plan-guide-new-change-set](../capability/rhiza/roadmap-app-start-runs/plan-guide-new-change-set.md)

## History

- `01M686RVAD1EHGERBBDJN3S45G` The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent
