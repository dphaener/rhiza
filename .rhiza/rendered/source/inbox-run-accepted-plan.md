# inbox-run-accepted-plan

source · `01M7ANSXHE8VTQFW47RR2AVVA8`

## quotes

```json
[]
```

## raw_sha256

c1f75440e3a57b82a080ffd0503298727f81cb4199264b50ccb617669137703b

## session

f54a56fc-1d63-4d1c-b35a-01a1d55af362

## summary

Autonomous plan-bug run, no owner present. Bug: rhiza inbox raised kind 'accept' items ('waiting for acceptance', next rhiza accept --rule owner) for change sets 01M79KGRYE, 01M79KGHYF and 01M79KGZ51 whose autonomous plan runs had already journaled 'accepted' at their gate, because the inbox reads main's model, where the change set stays proposed until the plan PR merges. Research: inboxItems (internal/api/inbox.go) raises an accept item for each change set in st.Gated that main's model says is proposed; the run's own accept is in its journal as a TypeAccepted entry (st.Accepts, written by autonomousGate.journal in internal/api/lifecycle.go when the gate accepts under rule autonomous or autonomy). The bug offers two remedies: raise no accept item, or label it 'accepted on branch, PR #N merging'. The planner chose to raise no item: the change set needs nothing from the owner, the run already carries it to merge, and a label would still list a non-actionable item and need a PR lookup the inbox does not do. Scope: the gated accept item only; the change set serves m4-orchestrated-by-rhiza (the owner accepts plans and answers the inbox), as the sibling inbox fix 01M7A1PX does.

## title

Plan of the inbox accept-item bug, autonomous plan-bug run

## Links

### In

- derived-from [run-accepted-no-inbox-item](../decision/run-accepted-no-inbox-item.md)

## History

- `01M7ANMWSVKNBSNPYCZXPK41GC` Fix The inbox asks the owner to accept a change set that an autonomous plan run already accepted. On 2026-11-13, plan runs 01M7AH6Q9G (change set 01M79KGRYE) and 01M7AH6RHK (01M79KGHYF) journaled 'gated' then 'accepted' at 14:20Z (log entries on the plan branch, e.g. 01M7AMGNTAEAQSSDF5YFMHXPYD), opened plan PRs #550/#551 and moved on to wait-rebased; meanwhile rhiza inbox listed both as kind 'accept', 'waiting for acceptance', next 'rhiza accept --rule owner --change ...' (items e512e92, 400de15), because main's model still says proposed until the plan PR merges. The same happened for 01M79KGZ51 (run 01M7AH6R61) before its PR merged. Nothing needed doing: the items clear when the PR merges. Wanted: the inbox does not raise an accept item for a change set whose gated run journaled 'accepted' (or labels it 'accepted on branch, PR #N merging'), so the owner is never asked to do what the run already did.
