# t8-d2-budget-stored-views

decision · `01M5SWSEBQVG3HGBN807ZVKR8M`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Budget as-is/to-be only (Recommended): Use stored views for as-is and to-be (accepted change sets live there), compute what-if on demand, and put the 50 ms budget test on as-is/to-be node and change slices. Fixing what-if/catch-up refolding stays the follow-up task 5 named. Keeps this PR off task 7's index work.

## statement

as-is and to-be read the stored views through the index; what-if uses `index.WhatIf` and slices the resulting `model.Model`. The 50 ms budget test covers only as-is/to-be. What-if speed stays a follow-up.

## wording

agent of milestone-1 task 8 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-8](../source/plan-task-8.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [slice](../capability/genesis/views/slice.md)

## History

- `01M5SWSE0Z4QGWQGN2FC9WMVAW` Task 8's plan decisions (PR #9, merge a951e38), accepted by the owner
