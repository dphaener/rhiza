# t2-d9-atomic-write

decision · `01M5SWRT8TBGWEHS6DCVVESSXV`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Temp + fsync + link (Recommended): Readers never see a partial file, and an existing entry can't be clobbered even by a racing writer; the reader ignores dot-prefixed temp files.

## statement

Atomic, never-overwrite write: temp file in the change set folder, fsync, hard-link to the final name (fails if it exists), remove temp, fsync the folder

## wording

agent of milestone-1 task 2 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-2](../source/plan-task-2.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [operation-log](../capability/rhiza/model/operation-log.md)

## History

- `01M5SWRT0HG6ATCA3BBJR35JZ5` Task 2's plan decisions (PR #2, merge 396505b), accepted by the owner
