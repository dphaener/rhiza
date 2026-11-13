# ship-merge-refused-stops

capability · `01M76YXZ5T5BQGQXD3VNC24CFS` · parent [shipping](../shipping.md)

## acceptance

Go tests in internal/api and internal/workflow (go test -trimpath ./internal/api -run ShipMergeRefused and go test -trimpath ./internal/workflow -run ShipChange) with a fake Exec assert: a gh pr merge failing with 'the base branch policy prohibits the merge' makes rhiza ship merge report a refused merge naming that reason, and the ship-change workflow stops the run there without visiting gate-rebase-merge; a merge failing because the head branch is behind main or conflicting still goes to gate-rebase-merge as before.

## statement

When GitHub refuses a ship-change run's merge for a reason other than the PR being behind main or conflicting (a branch-protection rule, a required check, a merge method the repo does not allow), the merge step stops the run once with GitHub's reason instead of failing back through gate-rebase-merge and rebase-merge until max_visits.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [ship-merge-refused-stops-verified](../../../../evidence/ship-merge-refused-stops-verified.md)
- verifies [ship-merge-refused-stops-qa-1](../../../../evidence/ship-merge-refused-stops-qa-1.md)
- verifies [ship-merge-refused-stops-qa-2](../../../../evidence/ship-merge-refused-stops-qa-2.md)
- verifies [pr-515-held-run-error-sources](../../../../evidence/pr-515-held-run-error-sources.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
