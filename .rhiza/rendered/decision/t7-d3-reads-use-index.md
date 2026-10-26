# t7-d3-reads-use-index

decision · `01M5SWSACMZXFX3G4TZ7PRYEAY`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: The index only (Recommended): Each read opens the worktree's SQLite index (catching up), and new index queries are added where needed (list by kind, subtree). The writes keep the interim api.Model; swapping that is not this task.

## statement

Reads use the index only (open + catch-up); new index queries where needed. Writes and the interim `api.Model` are untouched.

## wording

agent of milestone-1 task 7 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-7](../source/plan-task-7.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [fast-reads](../capability/rhiza/views/fast-reads.md)

## History

- `01M5SWS9YBW0KGX13FGMDHMKTG` Task 7's plan decisions (PR #10, merge 29e811d; PR #11, merge 3b5d1da), accepted by the owner
