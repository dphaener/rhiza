# ship-push-merges-own-origin

capability · `01M76XZG4SPZ8481AWNBER0JJS` · parent [shipping](../shipping.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run ShipPushDiverged) on scratch repos with a bare origin assert: open-pr and rebase, with origin's branch one commit ahead on a file the local branch did not touch, merge it and push, and origin's history keeps that commit; with a conflicting origin commit the step stops with a detail naming the commit and leaves no merge in progress; no call uses --force or rebase.

## statement

When a ship step's push is refused because the run's branch on origin holds commits the local branch lacks, the step fetches origin's branch and merges it (never force-pushes or rebases) when it merges cleanly, then pushes; a conflicting merge stops the run naming the origin commits.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [ship-push-merges-own-origin-verified](../../../../evidence/ship-push-merges-own-origin-verified.md)
- verifies [ship-push-merges-own-origin-qa-1](../../../../evidence/ship-push-merges-own-origin-qa-1.md)
- verifies [ship-push-merges-own-origin-verified-2](../../../../evidence/ship-push-merges-own-origin-verified-2.md)
- verifies [ship-push-merges-own-origin-qa-2](../../../../evidence/ship-push-merges-own-origin-qa-2.md)
- verifies [pr-515-held-run-error-sources](../../../../evidence/pr-515-held-run-error-sources.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
