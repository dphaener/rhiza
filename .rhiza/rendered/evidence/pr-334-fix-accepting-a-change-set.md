# pr-334-fix-accepting-a-change-set

evidence · `01M6S842M3FM0JP8WG5NFSDYGC` · actor [dphaener](../actor/dphaener.md)

## head

928f0b6a53f04d6fa7b272e0b642b3cde26d629e

## pr

https://github.com/dphaener/rhiza/pull/334

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #334 (https://github.com/dphaener/rhiza/pull/334) passed CI at 928f0b6

## status

verified

## Links

### Out

- verifies [run-worktree-carries-accept](../capability/rhiza/orchestration/shipping/run-worktree-carries-accept.md)

## History

- `01M6RXNQFZX8YM80WFZJCXBRAJ` Fix Accepting a change set from the Inbox app never reaches the run that builds it. On 2026-11-06 the owner accepted Import procedure (01M6J36X) in the Inbox app at 14:31:10Z; the accept entry was written to the main clone's .rhiza/log but stayed uncommitted and unpushed. The daemon started the ship-change run 11 s later on a branch from origin, without the entry, and the run built, passed review and CI, then failed at land: 'cannot ship change set 01M6J36X...: it is proposed, not accepted' (the orchestrator copied the entry onto the run branch by hand to get it merged as #328). An accept from the Inbox app (and rhiza accept on the main clone generally) must reach the run's branch, e.g. by committing and pushing the accept to main or by the run carrying it, so the run never sees an accepted change set as proposed.
