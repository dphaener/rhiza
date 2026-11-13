# ship-merge-refused-stops-verified

evidence · `01M792KKZ2EAH7CHEY8AEGQWN3` · actor [claude-code-2d126e52](../actor/claude-code-2d126e52.md)

## result

PASS

## statement

Verified at ec2d8a5bf (attempt 2). go test -trimpath ./internal/api -run ShipMergeRefused (and -run ShipMerge) pass: TestShipMergeRefusedStops, with a fake gh pr merge failing 'the base branch policy prohibits the merge' (bare, and as gh's full messages 'X Pull request owner/repo#7 is not mergeable: the base branch policy prohibits the merge.' with and without the X), asserts ship merge exits 75 with a stop naming that reason and records no merge commit; gh's full 'not up to date' and 'cannot be cleanly created' messages, plus behind/conflict variants, still fail (exit 1, to gate-rebase-merge). go test -trimpath ./internal/workflow -run ShipChange passes: TestShipChangeMergeRefusedStops checks merge's on.fail is gate-rebase-merge and on.pass wait-main, and go test ./internal/run -run StopExitStopsRun passes, showing an exit-75 command stops the run without taking on.fail, so a refused merge never visits gate-rebase-merge. Diff read: the attempt-1 finding is fixed ('not mergeable' dropped from shipMergeCurable, which now matches only the reason after gh's prefix); execIn passes gh's stderr into the error text, so the real message reaches shipMergeRefusal; describe-golden and cli.md updated for the stop field.

## Links

### Out

- verifies [ship-merge-refused-stops](../capability/rhiza/orchestration/shipping/ship-merge-refused-stops.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
