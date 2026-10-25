# t5-d4-input-incremental

decision · `01M5SWS28N1KWTM4F6XWMR7MJ4`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Input-incremental (Recommended): Only new entry files are read and decoded, only touched change sets are refolded, and a view (as-is, to-be) is recomputed whole in memory only when its set of change sets or their ops changed, then written by diff; uses internal/model unchanged.

## statement

Input-incremental: read and decode only new entry files, refold touched change sets, recompute a view whole in memory with `internal/model` only when its inputs changed, write its rows by diff.

## wording

agent of milestone-1 task 5 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-5](../source/plan-task-5.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [fast-reads](../capability/genesis/views/fast-reads.md)

## History

- `01M5SWS22WXZMST7XXR9ZKY4EE` Task 5's plan decisions (PR #7, merge ed8d292), accepted by the owner
