# ship-push-merges-own-origin-qa-2

evidence · `01M7963HDR6E70Y76T92TFZD2M` · actor [claude-code-dc3d533e](../actor/claude-code-dc3d533e.md)

## result

PASS

## statement

QA claude-code-dc3d533e at 99e340d22 (cycle 2), hands-on with the built binary driving real ship-change runs (rhiza run start/approve --foreground) in a scratch Rhiza root with a bare origin and fake gh, rex and claude on PATH: the builder pushed its branch then a separate clone pushed 'Earlier push' (a new file) to origin's cs/ship-it: open-pr merged it ('Merge origin/cs/ship-it into cs/ship-it'), pushed, origin's history kept the Earlier push commit, no force push, and the run reached merge. With the origin commit conflicting with a later local commit on feature.txt, open-pr failed exit 3 with 'origin's cs/ship-it has commits that conflict with ...: 659e515 Earlier push; reconcile them before shipping', no MERGE_HEAD left, origin unchanged. go test -trimpath ./internal/api -run ShipPushDiverged passes.

## Links

### Out

- verifies [ship-push-merges-own-origin](../capability/rhiza/orchestration/shipping/ship-push-merges-own-origin.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
