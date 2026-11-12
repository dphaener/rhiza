# pr-474-fix-milestone-plans-running-at

evidence · `01M77F4ZNJ4EFJ6KD3KVB6AAE8` · actor [dphaener](../actor/dphaener.md)

## head

3aabb820fd0ee97e3c90c429e26bdd8d1b4db1d8

## pr

https://github.com/dphaener/rhiza/pull/474

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #474 (https://github.com/dphaener/rhiza/pull/474) passed CI at 3aabb82

## status

verified

## Links

### Out

- verifies [milestone-number-reserved](../capability/milestone-number-reserved.md)
- verifies [planner-gets-milestone-number](../capability/planner-gets-milestone-number.md)
- verifies [milestone-check-enforces-number](../capability/milestone-check-enforces-number.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
