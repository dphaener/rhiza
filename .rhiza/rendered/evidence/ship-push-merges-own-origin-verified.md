# ship-push-merges-own-origin-verified

evidence · `01M7929VD2C49YYPRA4CZQ0KXG` · actor [claude-code-089110d2](../actor/claude-code-089110d2.md)

## result

PASS

## statement

Verified at 5dccff371 (attempt 2), detached checkout. Ran `go test -trimpath ./internal/api -run ShipPushDiverged -v -count=1`: PASS for TestShipPushDivergedMergesOrigin/{ship.open-pr,ship.rebase} and TestShipPushDivergedConflictStops/{ship.open-pr,ship.rebase}. These use scratch repos with a bare origin whose cs/ship-it branch is one commit ahead (pushed from another clone). In the clean case (origin touches theirs.txt, local touches mine.txt) the step merges and pushes: origin's branch equals local HEAD, history holds both 'earlier push' and 'mine', and a merge commit exists. In the conflict case (both touch clash.txt) the error names 'earlier push', there is no MERGE_HEAD, HEAD is unchanged and nothing is pushed. Every git call made through ctx.Git and ctx.Exec is recorded, and noForceOrRebase asserts that none is a rebase, --force, -f push, --force-with-lease or + refspec push. That closes attempt 1's findings f1 (rebase not covered) and f2 (no force/rebase assertion). Read the diff: shipFollowPR now calls shipMergeOrigin on divergence (dirty-tree refusal, merge --no-edit, merge --abort plus Conflictf naming origin's commits on failure). ship.rebase reaches it through shipFollowPR. Neither file uses --force or rebase. Also ran: `go test -trimpath ./internal/api -run 'ShipOpenPR|FollowPR|ShipRebase|ShipClose|ShipLand|ShipMerge' -count=1`: ok. `go vet ./internal/api`: clean.

## Links

### Out

- verifies [ship-push-merges-own-origin](../capability/rhiza/orchestration/shipping/ship-push-merges-own-origin.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
