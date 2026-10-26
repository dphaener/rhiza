# t7-d2-default-view

decision · `01M5SWSA5C6GVZ9X2ZSGZS54QP`

## accepted_by

owner

## revised

revised by the orchestrator under the owner's authorization: one default-view rule shared by the reads (task 7) and slice (task 8): the effective change set draft or proposed -> what-if of it; accepted -> to-be; none or landed -> as-is. Work happens on accepted change sets, so the main loop stays on the fast indexed to-be path; what-if stays outside the 50 ms budget. (The agent's accepted option had been: to-be in a change set, as-is otherwise.)

## settled_by

Agent's option, accepted by the owner: Literal: to-be or as-is (Recommended): to-be (landed + accepted) when a change set is current or given by --change / GENESIS_CHANGE, as-is otherwise. what-if:<cs> only when asked for, and exempt from the 50 ms budget. Fast; your draft's edits need --view what-if:<cs>.

## statement

Default view, one rule shared with task 8's `slice`: the effective change set (`--change`, `GENESIS_CHANGE`, or the current one) draft or proposed → what-if of it; accepted → to-be; none, landed or otherwise terminal → as-is (except D17). what-if stays exempt from the budget. (Revised by the owner's session after approval; the original rule was to-be in a change set, as-is otherwise.)

## wording

agent of milestone-1 task 7 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-7](../source/plan-task-7.md)
- derived-from [task-interviews](../source/task-interviews.md)
- derived-from [founding-conversation](../source/founding-conversation.md)

### In

- decided-by [overlays](../capability/rhiza/views/overlays.md)

## History

- `01M5SWS9YBW0KGX13FGMDHMKTG` Task 7's plan decisions (PR #10, merge 29e811d; PR #11, merge 3b5d1da), accepted by the owner
