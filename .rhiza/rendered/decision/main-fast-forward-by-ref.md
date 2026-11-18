# main-fast-forward-by-ref

decision · `01M7Q2VZTECZATGK3FTWECSN15`

## statement

Close-out and the model-PR merge update the main checkout with git fetch origin main followed by git merge --ff-only refs/remotes/origin/main, through one shared helper, and never with git pull, so a concurrent fetch rewriting FETCH_HEAD and the user's pull.rebase or pull.ff settings cannot fail the update. Only the git commands change: the step is still refused when main is not on main or has tracked changes, and a failure is still recorded as main-stale.

## wording

the planner's design from its research, in an autonomous plan run with no owner present

## Links

### Out

- supersedes [close-out-fast-forwards-main](close-out-fast-forwards-main.md)
- derived-from [close-out-pull-race-plan](../source/close-out-pull-race-plan.md)

### In

- decided-by [main-update-without-pull](../capability/rhiza/orchestration/shipping/main-update-without-pull.md)

## History

- `01M7Q2Q68JE50NCRKGCT5Z3APZ` Fix watchdog transient failure at step plan/close-out: run 01M7PSG0M7APCTB4V1YKPTM394, request 01M7KHWSWG3ENZ81KH148SK1GD; log: main was not updated: git pull --ff-only origin main: From https://github.com/dphaener/rhiza | * branch main -> FETCH_HEAD | fatal: Cannot rebase onto multiple branches.
