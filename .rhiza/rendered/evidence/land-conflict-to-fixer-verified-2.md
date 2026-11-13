# land-conflict-to-fixer-verified-2

evidence · `01M7955SSW577Z3GNXJPWZJWF0` · actor [claude-code-153822b1](../actor/claude-code-153822b1.md)

## result

PASS

## statement

Verified at 3f47460c6 (attempt 2) in a detached checkout. go test -trimpath ./internal/api -run ShipLand -v: all PASS, including TestShipLandNodeChangedMergesMain (head build exiting 3 with 'cannot land change set X: node N changed' returns CodeConflict naming node N and the conflicted file and records the merge conflict) and TestShipLandOtherHeadFailureStaysPlain (other message / other exit stay a plain land failure with the head build's detail; TestShipLandHeadBuildErrorDetail also passes). go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-land-conflict$': ok; the script runs ship-change with fake gh through land fail -> land-triage -> land-conflict-triage -> land-main-conflict-triage pass -> conflict-fix -> land visit 2 pass, and case 2 fails with no fixer. The cycle-1 blocker (actors shared across sessions stopping runs at build/check) no longer reproduces at this head. Also ran go test -trimpath ./internal/api -run ShipFollow: ok (quality f2's TestShipFollowRefusesLocalCommits passes). The leaf commit itself only adds .rhiza log entries.

## Links

### Out

- verifies [land-conflict-to-fixer](../capability/rhiza/orchestration/shipping/land-conflict-to-fixer.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
