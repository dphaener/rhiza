# pr-421-fix-rhiza-inbox-said-nothing

evidence · `01M71N7X5KSQ6SF638QR29Q0N6` · actor [dphaener](../actor/dphaener.md)

## head

64d24881c7d3166d0c481f7aa5ad9865f309c763

## pr

https://github.com/dphaener/rhiza/pull/421

## result

11 checks passed: lint SUCCESS, go-install SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #421 (https://github.com/dphaener/rhiza/pull/421) passed CI at 64d2488

## status

verified

## Links

### Out

- verifies [inbox-lists-gate-stopped-plan-run](../capability/rhiza/inbox-lists-gate-stopped-plan-run.md)

## History

- `01M71HJW9W123PXWWAGZTHGE6Z` Fix rhiza inbox said 'Nothing waits on you' while plan run 01M71E6F15HMTME3KKF96NJC5R (plan-milestone, owner autonomy) was stopped at step gate on open question review-01m71e6f15hmtme3kkf96njc5r-r3-f1 (left by record 01M71HBAFK92B6EQR0FGKN9H8M, plan review round 3). It was found only because rhiza run send-back refused with 'stopped at step gate: open questions'. Expected: the inbox lists the stopped run and its question, worded as what is asked. Unstuck by asking the owner, recording a Decision answering the question (record change set, commit 0ef646e in the run worktree), then rhiza run resume; the run went on to open PR #419. Daemon f7c85e91, 2026-11-10 ~01:35Z.
