# plan-run-merges-closes-out-verified

evidence · `01M694ZV3PQANDRY7GPE9F88XQ` · actor [claude-code-180b619b](../actor/claude-code-180b619b.md)

## result

PASS

## statement

At 26af883 (on 10d4d82) ran go test ./cmd/rhiza -run TestScripts/plan-workflow -count=1 (ok) and go test ./internal/api -run 'TestShipMergeAfterPlanMerge|TestShipCloseOut' -count=1 -v (all PASS); make fmt-check vet test pass (self-check only flagged the render of this verifier's own session actor, not in the commit). Read plan-workflow-pr.txtar: with fake gh green checks the plan run calls 'pr merge 7 --squash --match-head-commit <checked head>' once, the planned change set stays proposed, cs/plan-reports is gone locally and on origin, the worktree is removed and the ship state dir is empty; then a ship-change run on the change set succeeds, logs a second gh pr merge ('pr merge 8 --squash --match-head-commit <head>'), and its merge step's JSON has a merge_commit and no already_merged (AlreadyMerged false). TestShipMergeAfterPlanMerge starts from ship state holding a plan's merge commit, runs open-pr, check, merge, and asserts AlreadyMerged false with one gh pr merge call. Attempt-1 finding f1 is fixed.

## Links

### Out

- verifies [plan-run-merges-closes-out](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-merges-closes-out.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
