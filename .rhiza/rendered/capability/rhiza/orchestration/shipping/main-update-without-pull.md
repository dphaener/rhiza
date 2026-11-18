# main-update-without-pull

capability · `01M7Q2VZY1G17DVAGMK5H1JZ1R` · parent [shipping](../shipping.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShipCloseOut|TestDaemonModelPR' passes, with: (1) the close-out unit tests in internal/api/ship_test.go (the fakes now answering git fetch origin main and git merge --ff-only refs/remotes/origin/main in place of git pull --ff-only origin main) asserting that a successful close-out runs that fetch and that merge in the main checkout, runs no git command whose first argument is pull, and reports update-main done; (2) (failure case) a close-out whose git merge --ff-only refs/remotes/origin/main fails reports update-main failed, main_updated false, a detail starting 'main was not updated: ' with git's error, and records main-stale; (3) (failure case) a close-out whose git fetch origin main fails does the same and runs no merge; (4) TestDaemonModelPR (or a new TestDaemonModelPR* test) in internal/api/model_prs_test.go where a model PR merges asserts the main checkout is updated by the same fetch and merge, with no git pull, and (failure case) a failing merge records main-stale. go test -trimpath ./cmd/rhiza -run 'TestScripts/(ship-close-out-forget-retry|ship-close-out-log-copies)' passes, and grep -rn '"pull"' internal/api/*.go | grep -v _test.go prints nothing.

## files

internal/api/ship_close.go, internal/api/model_prs.go, internal/api/ship_test.go, internal/api/model_prs_test.go

## must_not

Change when the update is refused (shipMainBlocked) or what dropIncomingLogCopies drops; change the main-stale record, the update-main step's name or results, or the ShipCloseOutResponse JSON shape; rebase, reset or force main; retry in a loop; touch the worktree fetches in ship.go, ship_rebase.go or watchdog_noci.go, or the daemon self-update.

## statement

ship close-out and pullMainAfterModelMerge fast-forward the main checkout through one helper that runs git fetch origin main and then git merge --ff-only refs/remotes/origin/main, never git pull, so FETCH_HEAD written by a concurrent fetch and a global pull.rebase=true no longer fail the update.

## Links

### Out

- decided-by [main-fast-forward-by-ref](../../../../decision/main-fast-forward-by-ref.md)

### In

- verifies [main-update-without-pull-verified](../../../../evidence/main-update-without-pull-verified.md)
- verifies [main-update-without-pull-qa-1](../../../../evidence/main-update-without-pull-qa-1.md)

## History

- `01M7Q2Q68JE50NCRKGCT5Z3APZ` Fix watchdog transient failure at step plan/close-out: run 01M7PSG0M7APCTB4V1YKPTM394, request 01M7KHWSWG3ENZ81KH148SK1GD; log: main was not updated: git pull --ff-only origin main: From https://github.com/dphaener/rhiza | * branch main -> FETCH_HEAD | fatal: Cannot rebase onto multiple branches.
