# t8-d1-slice-shape

decision · `01M5SWSE3S9EMJG4FZQ6TZV06Q`

## accepted_by

owner

## revised

revised by the orchestrator under the owner's authorization: change D1's default view so reads and slice share one rule across tasks 7 and 8. Effective change set draft or proposed -> what-if of it; accepted -> to-be; none or landed -> as-is. Work happens on accepted change sets, so the main loop stays on the fast indexed to-be path; what-if stays out of the 50 ms budget.

## settled_by

Agent's option, accepted by the owner: slice [ref] + --view (Recommended): `slice <ref>` slices one node; with no ref, slices the effective change set (--change / GENESIS_CHANGE / current). --view as-is|to-be|what-if defaults to what-if when that change set is in flight, else as-is. Matches the guide's two forms and reuses the global --change.

## statement

`slice [ref]`: a ref gives a node slice; no ref gives a change slice of the effective change set (`--change` / `GENESIS_CHANGE` / current; none is a usage error with next `genesis slice <ref>`). `--view as-is|to-be|what-if`. Default view (amended by the owner's orchestrator on 2026-10-24, shared with task 7's reads): the effective change set is draft or proposed → what-if of it; accepted → to-be; none or landed (or any other terminal state) → as-is. A node slice with what-if needs a draft or proposed change set.

## wording

agent of milestone-1 task 8 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-8](../source/plan-task-8.md)
- derived-from [task-interviews](../source/task-interviews.md)
- derived-from [founding-conversation](../source/founding-conversation.md)

### In

- decided-by [slice](../capability/genesis/views/slice.md)

## History

- `01M5SWSE0Z4QGWQGN2FC9WMVAW` Task 8's plan decisions (PR #9, merge a951e38), accepted by the owner
