# question-9a4a7e78-1

question · `01M7AHKCZADEQ2MFGYWYMMSKBQ`

## session

9a4a7e78-6eb7-4185-b879-01a1d51484a5

## statement

The new builtin ship-change (conduct, gate, signoff, push/PR, CI, merge, main CI, close-out): should it keep the classic's last two steps, `rhiza ship command` (.rhiza/workflows/ship.sh) and `rhiza ship deploy` (rhiza.ship.deployUrl), after close-out? The milestone leaf caps the builtin at 10 steps, which leaves no room for them. / In the new builtin ship-change, what should happen when main's CI is red after the merge? Classic stops the run there. The change is already merged, so a conductor on the branch cannot fix it. / Rhiza's repo shadows ship-change only to build rhiza from the PR head (SHIP_HEAD_RHIZA / SHIP_HEAD_BUILD in [env]), so ship renders with the branch's own rhiza. If the shadow is removed, where should that [env] go? / About 50 testscripts and several Go tests (internal/api run_test, workflow_test, daemon pick tests) drive the classic ship-change's steps: land, rebase, ci-fix, the infra waits and so on. When the builtin becomes the conductor default, what should happen to them?

## Links

### In

- answers [builtin-ship-change-open-points](../decision/builtin-ship-change-open-points.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
