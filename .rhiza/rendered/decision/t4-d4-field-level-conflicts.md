# t4-d4-field-level-conflicts

decision · `01M5SWRZR4YT228YN64QN8PCBA`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Field-level overlap (Recommended): Matches "updates name fields; edits to different fields are not a real conflict"; retire treated as touching every field.

## statement

Conflict = stale base and overlap: update vs an earlier write to the same field or the handle; move vs an earlier move; a retire, and any op after a retire, overlaps everything; link.retire vs link.retire. Non-overlapping stale ops merge silently

## wording

agent of milestone-1 task 4 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-4](../source/plan-task-4.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [conflicts](../capability/rhiza/change-sets/conflicts.md)

## History

- `01M5SWRZJR0EXNHS4GAQDAJPHZ` Task 4's plan decisions (PR #4, merge 380f2f2), accepted by the owner
