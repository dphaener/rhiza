# ship-push-merges-own-origin-verified-2

evidence · `01M795CYPXYYE04KG6D5Z99FB9` · actor [claude-code-a76a3b3a](../actor/claude-code-a76a3b3a.md)

## result

PASS

## statement

At 3414d9ac7 (attempt 2), in a detached checkout: go test -trimpath ./internal/api -run 'ShipPushDiverged|TestShipFollow|TestShipRebase' -v -count=1 PASS. TestShipPushDivergedMergesOrigin (ship.open-pr and ship.rebase subtests) uses scratch repos with a bare origin. Another clone pushes a commit on theirs.txt to origin's cs/ship-it, and the local branch commits mine.txt. The step merges (a merge commit exists), pushes, origin's ref equals local HEAD, and the history keeps 'earlier push'. TestShipPushDivergedConflictStops (both steps) makes a conflicting clash.txt: the error names 'earlier push', there is no MERGE_HEAD, HEAD is unmoved and nothing is pushed. noForceOrRebase records every git call through ctx.Git and ctx.Exec and finds no rebase, --force, -f, --force-with-lease or +refspec. Read shipMergeOrigin in internal/api/ship_close.go: merge --no-edit, then merge --abort plus a Conflictf naming the HEAD..origin commits on failure. Attempt 1's finding f2 (TestShipFollowRefusesLocalCommits red) is fixed: the test is rewritten as TestShipFollowRefusesConflictingOrigin and PASSes. The full go test -trimpath ./internal/api -count=1 also passes (ok, 594s). The other attempt-1 findings (actor reuse in session.go) are outside this leaf.

## Links

### Out

- verifies [ship-push-merges-own-origin](../capability/rhiza/orchestration/shipping/ship-push-merges-own-origin.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
