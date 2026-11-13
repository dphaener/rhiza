# land-conflict-to-fixer

capability · `01M76XZG5J7ETYBRT4W550P79M` · parent [shipping](../shipping.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run ShipLand) assert a land whose head build exits 3 with 'cannot land change set X: node N changed' reports a conflict naming node N (not a plain failure), and a testscript (go test -trimpath ./cmd/rhiza -run TestScripts/<the ship land conflict script>) runs ship-change with a fake gh to the conflict-fix step after that land failure; any other land failure still fails land with the head build's detail (ship-head-build-error-detail).

## statement

A ship land that fails because the PR head's rhiza refuses to land (exit 3, a node the change set edits changed on main) merges origin/main into the branch and hands the conflict to the run's conflict fixer, instead of failing land and holding the run.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../decision/efficient-cycle-milestone-shape.md)
- decided-by [land-conflict-script-named](../../../../decision/land-conflict-script-named.md)

### In

- blocks [review-01m76stxjb7cgga48pw43p487t-r3-f1](../../../../question/review-01m76stxjb7cgga48pw43p487t-r3-f1.md)
- verifies [land-conflict-to-fixer-verified](../../../../evidence/land-conflict-to-fixer-verified.md)
- verifies [land-conflict-to-fixer-qa-1](../../../../evidence/land-conflict-to-fixer-qa-1.md)
- verifies [land-conflict-to-fixer-verified-2](../../../../evidence/land-conflict-to-fixer-verified-2.md)
- verifies [land-conflict-to-fixer-qa-2](../../../../evidence/land-conflict-to-fixer-qa-2.md)
- verifies [pr-515-held-run-error-sources](../../../../evidence/pr-515-held-run-error-sources.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
