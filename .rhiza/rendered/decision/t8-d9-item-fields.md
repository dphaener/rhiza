# t8-d9-item-fields

decision · `01M5SWSFRVKFM7BP5547PKDBJN`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: statement + evidence status (Recommended): One-line items use `statement` (falling back to the handle). An evidence node's `status` field reads `verified` or `failing`; anything else, or no field, is `unverified`. A node's evidence status is failing > verified > unverified across its `verifies` links, `none` with no evidence. Focus shows every field.

## statement

Items read `statement`, falling back to the handle. An evidence node's `status` field is the JSON string `verified` or `failing`; anything else is `unverified`. A node's evidence status across live `verifies` links is failing > verified > unverified, and `none` with no evidence. Focus shows every field.

## wording

agent of milestone-1 task 8 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-8](../source/plan-task-8.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [slice](../capability/rhiza/views/slice.md)

## History

- `01M5SWSE0Z4QGWQGN2FC9WMVAW` Task 8's plan decisions (PR #9, merge a951e38), accepted by the owner
