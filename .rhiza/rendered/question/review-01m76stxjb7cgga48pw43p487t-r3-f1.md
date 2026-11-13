# review-01m76stxjb7cgga48pw43p487t-r3-f1

question · `01M76ZQQ2KDNBTW2B0HEK7VWGS`

## statement

Plan review round 3 left finding f1 (missing-acceptance on land-conflict-to-fixer): The acceptance names its testscript only as a placeholder: 'go test -trimpath ./cmd/rhiza -run TestScripts/<the ship land conflict script>'. A builder or verifier cannot run that command as written. Name the script (for example ship-land-conflict) so the check is concrete. Also note that ship-change.toml's land step has no on.fail today, so the acceptance should say which step the conflict routes to (conflict-fix, or conflict-fix-landed). The critic confirmed it: The acceptance of land-conflict-to-fixer (in change set 01M76XSQKQJ8P78585NDYPATM7) names its testscript as 'TestScripts/<the ship land conflict script>', a placeholder no builder or verifier can run. The script needs a concrete name (e.g. run-ship-land-conflict, beside the existing run-ship-conflict-fix.txtar). The routing part is already stated: the acceptance says the run reaches the conflict-fix step, which exists in internal/workflow/builtin/ship-change.toml. How should it be settled?

## Links

### Out

- blocks [land-conflict-to-fixer](../capability/rhiza/orchestration/shipping/land-conflict-to-fixer.md)

### In

- answers [land-conflict-script-named](../decision/land-conflict-script-named.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
