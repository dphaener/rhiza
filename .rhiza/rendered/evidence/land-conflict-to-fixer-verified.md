# land-conflict-to-fixer-verified

evidence · `01M7923T15GC0V00F5V2W6XMM4` · actor [claude-code-8452d0f0](../actor/claude-code-8452d0f0.md)

## result

PASS

## statement

Verifier at 609344f55. go test -trimpath ./internal/api -run ShipLand: ok; TestShipLandNodeChangedMergesMain asserts an exit-3 'cannot land change set X: node N changed' head-build failure returns a conflict naming node N and the conflicted file, with the merge conflict recorded; TestShipLandOtherHeadFailureStaysPlain asserts another exit-3 message and an exit-1 keep the head build's detail, merge nothing and record no conflict (TestShipLandHeadBuildErrorDetail also passes). go test -trimpath ./internal/workflow/...: ok (land-conflict-triage on.fail -> land-main-conflict-triage, which passes to conflict-fix only). go test -trimpath ./cmd/rhiza -run TestScripts/run-ship-land-conflict fails at this commit at build/check ('Evidence ... written by claude-code-..., who wrote its commit'), as does the untouched run-ship-model-conflict script: a base regression from d38f497d2 (shared actor handle), not this leaf. With 609344f55 cherry-picked onto d38f497d2^ in a scratch worktree, run-ship-land-conflict passes: case 1 goes land fail -> land-triage -> land-conflict-triage -> land-main-conflict-triage pass -> conflict-fix -> render-pr -> land visit 2 pass, never conflict-fix-landed; case 2 (other exit-3 message) fails the run at land-main-conflict-triage with no fixer. Diff read: shipLandNodeChanged requires exit status 3, no structured model conflict, and the node-changed message.

## Links

### Out

- verifies [land-conflict-to-fixer](../capability/rhiza/orchestration/shipping/land-conflict-to-fixer.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
