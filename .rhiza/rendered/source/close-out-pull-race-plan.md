# close-out-pull-race-plan

source · `01M7Q2VE7VSZF9ABDDTATTW968`

## quotes

```json
[]
```

## raw_sha256

bb944fece358f595ff5ae0ad46c615e2083a82303b15779e153269ed24df5ba4

## session

85060fd4-900b-4356-803c-01a1ee2be88f

## summary

Autonomous plan run on change set 01M7Q2Q6, no owner present. Run 01M7PSG0M7AP failed at plan/close-out: 'git pull --ff-only origin main: fatal: Cannot rebase onto multiple branches'. The planner traced it. ship close-out (internal/api/ship_close.go) and the model-PR merge (pullMainAfterModelMerge in internal/api/model_prs.go) both fast-forward the main checkout with git pull --ff-only origin main. The machine's global git config sets pull.rebase=true, so git pull rebases onto what it reads from FETCH_HEAD. FETCH_HEAD sits in the main checkout's git dir, and other processes fetch there at the same time (the daemon's self-update runs git fetch origin main there each poll, and close-out's own dropIncomingLogCopies fetches just before). When the writes interleave, FETCH_HEAD holds more than one for-merge line, and a rebasing pull fails with 'Cannot rebase onto multiple branches'. A retry usually passes, which is why the watchdog saw the failure as transient. The fix the planner chose: never use git pull to update main. Fetch origin main, then git merge --ff-only refs/remotes/origin/main. That reads a ref, which git updates atomically, instead of FETCH_HEAD, and it ignores the user's pull.rebase and pull.ff settings. Both call sites share one helper. The planner judged this a technical choice that the code settles, not an owner question, so no question was raised. The change set serves m40-hands-off, because a close-out that needs a retry is an intervention. The bug is not severe: no data is lost, and main is only stale until the retry.

## title

Close-out pull race: plan research

## Links

### In

- derived-from [main-fast-forward-by-ref](../decision/main-fast-forward-by-ref.md)

## History

- `01M7Q2Q68JE50NCRKGCT5Z3APZ` Fix watchdog transient failure at step plan/close-out: run 01M7PSG0M7APCTB4V1YKPTM394, request 01M7KHWSWG3ENZ81KH148SK1GD; log: main was not updated: git pull --ff-only origin main: From https://github.com/dphaener/rhiza | * branch main -> FETCH_HEAD | fatal: Cannot rebase onto multiple branches.
