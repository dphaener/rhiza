# t5-d5-stored-views

decision · `01M5SWS29TSS62V86818WHZBT3`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: As-is + to-be stored (Recommended): Materialize both as rows (nodes, links, problems per view); what-if for one change set is computed on demand in memory from the stored folded change sets via model.WhatIf, and its cost is benchmarked and reported.

## statement

As-is and to-be materialized; what-if computed on demand from stored entries via `model.WhatIf`, its cost benchmarked.

## wording

agent of milestone-1 task 5 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-5](../source/plan-task-5.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- relates-to [whatif-folds-inflight-over-asis](whatif-folds-inflight-over-asis.md) — supersedes it for the what-if part: a what-if no longer runs model.WhatIf over every stored folded change set; it seeds model.WhatIfFrom from the stored as-is rows and decodes only the change set, its in-flight closure and its landed dependencies. As-is and to-be stay materialized as t5-d5 says
- decided-by [fast-reads](../capability/rhiza/views/fast-reads.md)

## History

- `01M5SWS22WXZMST7XXR9ZKY4EE` Task 5's plan decisions (PR #7, merge ed8d292), accepted by the owner
