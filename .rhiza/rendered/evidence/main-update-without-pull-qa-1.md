# main-update-without-pull-qa-1

evidence · `01M7Q3KASPF1HABBFBR5BFK6Q5` · actor [claude-code-77ded467](../actor/claude-code-77ded467.md)

## result

PASS

## statement

QA: a scratch build's ship close-out, under pull.rebase=true with a multi-line FETCH_HEAD written by a concurrent fetch, updated main via git fetch origin main + git merge --ff-only refs/remotes/origin/main with no git pull; an injected merge or fetch failure reported update-main failed with a 'main was not updated: ' detail and a main-stale inbox item, and a failed fetch ran no merge; the scoped internal/api tests and two close-out scripts pass; no "pull" in non-test internal/api.

## Links

### Out

- verifies [main-update-without-pull](../capability/rhiza/orchestration/shipping/main-update-without-pull.md)

## History

- `01M7Q2Q68JE50NCRKGCT5Z3APZ` Fix watchdog transient failure at step plan/close-out: run 01M7PSG0M7APCTB4V1YKPTM394, request 01M7KHWSWG3ENZ81KH148SK1GD; log: main was not updated: git pull --ff-only origin main: From https://github.com/dphaener/rhiza | * branch main -> FETCH_HEAD | fatal: Cannot rebase onto multiple branches.
