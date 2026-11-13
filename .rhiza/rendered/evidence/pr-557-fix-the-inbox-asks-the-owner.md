# pr-557-fix-the-inbox-asks-the-owner

evidence · `01M7AX61AE8Z9640DMA1BYH4ZY` · actor [dphaener](../actor/dphaener.md)

## head

68e36d5f61e580679bb0e3d656c06ae15a263ea5

## pr

https://github.com/dphaener/rhiza/pull/557

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #557 (https://github.com/dphaener/rhiza/pull/557) passed CI at 68e36d5

## status

verified

## Links

### Out

- verifies [inbox-skips-run-accepted-gated](../capability/inbox-skips-run-accepted-gated.md)

## History

- `01M7ANMWSVKNBSNPYCZXPK41GC` Fix The inbox asks the owner to accept a change set that an autonomous plan run already accepted. On 2026-11-13, plan runs 01M7AH6Q9G (change set 01M79KGRYE) and 01M7AH6RHK (01M79KGHYF) journaled 'gated' then 'accepted' at 14:20Z (log entries on the plan branch, e.g. 01M7AMGNTAEAQSSDF5YFMHXPYD), opened plan PRs #550/#551 and moved on to wait-rebased; meanwhile rhiza inbox listed both as kind 'accept', 'waiting for acceptance', next 'rhiza accept --rule owner --change ...' (items e512e92, 400de15), because main's model still says proposed until the plan PR merges. The same happened for 01M79KGZ51 (run 01M7AH6R61) before its PR merged. Nothing needed doing: the items clear when the PR merges. Wanted: the inbox does not raise an accept item for a change set whose gated run journaled 'accepted' (or labels it 'accepted on branch, PR #N merging'), so the owner is never asked to do what the run already did.
