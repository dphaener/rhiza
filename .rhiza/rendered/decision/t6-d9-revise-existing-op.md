# t6-d9-revise-existing-op

decision · `01M5SWS743JS3TYR6Q6K869T27`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Revise the existing op (Recommended): One op per node per change set where possible: `set` on a node added in this change set revises its node.create; a second `set` merges into the existing node.update; a second `mv` replaces the move; `retire` of a node added here drops its create. Uses the revision-by-op-ID mechanism already in changeset.

## statement

When the current change set already touches the node, the write revises the existing op: `set` on a node added here revises its `node.create`; a second `set` merges into the existing `node.update`; a second `mv` replaces the move; `retire` of a node added here drops its create

## wording

agent of milestone-1 task 6 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-6](../source/plan-task-6.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [revisions](../capability/genesis/change-sets/revisions.md)

## History

- `01M5SWS5QKWPQ3FS357EAV60D2` Task 6's plan decisions (PR #5, merge 6d5c61e), accepted by the owner
