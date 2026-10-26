# t6-d13-lifecycle-verbs

decision · `01M5SWS7TVE183FWM8KNG0EE99`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Current by default, rule=owner (Recommended): propose/accept/reject/withdraw/land act on --change or the current change set; accept/reject take --rule (default "owner") and --reason; an illegal move fails with exit 3 and next. A terminal move on the current change set clears it. --record opens, applies and lands in one entry chain and never changes the current one.

## statement

Lifecycle verbs act on `--change`/current. `accept`/`reject` take `--rule` (default `owner`) and `--reason`. An illegal move exits 3 with `next`. A terminal move on the current change set clears it. `--record` never changes the current change set

## wording

agent of milestone-1 task 6 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-6](../source/plan-task-6.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [change-lifecycle](../capability/rhiza/change-sets/change-lifecycle.md)
- decided-by [record-shortcut](../capability/rhiza/change-sets/record-shortcut.md)

## History

- `01M5SWS5QKWPQ3FS357EAV60D2` Task 6's plan decisions (PR #5, merge 6d5c61e), accepted by the owner
