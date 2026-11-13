# land-conflict-to-fixer-qa-2

evidence · `01M7963H23TDTMFD4XS70XDXPS` · actor [claude-code-dc3d533e](../actor/claude-code-dc3d533e.md)

## result

PASS

## statement

QA claude-code-dc3d533e at 99e340d22 (cycle 2), hands-on with the built binary driving real ship-change runs (rhiza run start/approve --foreground) in a scratch Rhiza root with a bare origin and fake gh, rex and claude on PATH: with a head build (SHIP_HEAD_RHIZA) that refuses the first land with exit 3 'cannot land change set X: node N changed' after main moved feature.txt, land failed with exit 3 reporting the conflict and leaving the merge for the fixer, land-main-conflict-triage passed to conflict-fix (not conflict-fix-landed), one fixer session, land visit 2 passed and the change set landed. With 'it depends on Y (draft)' instead, the run fails with no fixer (run-ship-land-conflict case 2 passes). go test -trimpath ./internal/api -run ShipLand passes.

## Links

### Out

- verifies [land-conflict-to-fixer](../capability/rhiza/orchestration/shipping/land-conflict-to-fixer.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
