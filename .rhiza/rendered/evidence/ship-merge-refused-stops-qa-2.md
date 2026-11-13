# ship-merge-refused-stops-qa-2

evidence · `01M7963HBSPMVHXYFXHQA4HKAF` · actor [claude-code-dc3d533e](../actor/claude-code-dc3d533e.md)

## result

PASS

## statement

QA claude-code-dc3d533e at 99e340d22 (cycle 2), hands-on with the built binary driving real ship-change runs (rhiza run start/approve --foreground) in a scratch Rhiza root with a bare origin and fake gh, rex and claude on PATH: gh pr merge failing with 'the base branch policy prohibits the merge': the run stopped once at merge with detail 'GitHub refused to merge pull request #7: ... the base branch policy prohibits the merge', one pr merge call, gate-rebase-merge visited once (before merge only). Merge refused once with 'the head branch is not up to date with the base branch': merge went to gate-rebase-merge visit 2, rebase-merge, and merged on visit 2; run succeeded. go test ./internal/api -run ShipMergeRefused and ./internal/workflow -run ShipChange pass.

## Links

### Out

- verifies [ship-merge-refused-stops](../capability/rhiza/orchestration/shipping/ship-merge-refused-stops.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
