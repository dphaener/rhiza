# run-accepted-no-inbox-item

decision · `01M7ANT76HV7Z963R14R7GWF3E`

## statement

The inbox raises no accept item for a gated change set whose run journaled its own accept, rather than labelling the item 'accepted on branch, PR #N merging': the run already carries the change set to merge, so the owner has nothing to do, and a label would keep a non-actionable item and need a PR lookup the inbox does not do

## wording

planner's pick of the two remedies the bug request offers; autonomous plan-bug run, no owner present

## Links

### Out

- derived-from [inbox-run-accepted-plan](../source/inbox-run-accepted-plan.md)

### In

- decided-by [inbox-skips-run-accepted-gated](../capability/inbox-skips-run-accepted-gated.md)

## History

- `01M7ANMWSVKNBSNPYCZXPK41GC` Fix The inbox asks the owner to accept a change set that an autonomous plan run already accepted. On 2026-11-13, plan runs 01M7AH6Q9G (change set 01M79KGRYE) and 01M7AH6RHK (01M79KGHYF) journaled 'gated' then 'accepted' at 14:20Z (log entries on the plan branch, e.g. 01M7AMGNTAEAQSSDF5YFMHXPYD), opened plan PRs #550/#551 and moved on to wait-rebased; meanwhile rhiza inbox listed both as kind 'accept', 'waiting for acceptance', next 'rhiza accept --rule owner --change ...' (items e512e92, 400de15), because main's model still says proposed until the plan PR merges. The same happened for 01M79KGZ51 (run 01M7AH6R61) before its PR merged. Nothing needed doing: the items clear when the PR merges. Wanted: the inbox does not raise an accept item for a change set whose gated run journaled 'accepted' (or labels it 'accepted on branch, PR #N merging'), so the owner is never asked to do what the run already did.
