# t9-d3-check-fail-warn

decision · `01M5SWSKJNX58GF1ZQ06GVM920`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Fail landed, warn in-flight (Recommended): FAIL: log read problems, fold problems except ErrConcurrent, any as-is problem (landed conflicts, skipped ops), endpoint violations on live as-is links, stale render, append-only violations. WARN: ErrConcurrent (an honest merge, unrepairable), problems only in to-be/what-if of in-flight change sets.

## statement

`check` **fails** on: log read problems, fold problems except `ErrConcurrent`, any as-is problem, endpoint violations on live as-is links, stale render, append-only violations. **Warns** on: `ErrConcurrent`, problems that appear only in to-be or what-if of in-flight change sets

## wording

agent of milestone-1 task 9 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-9](../source/plan-task-9.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [check](../capability/rhiza/integrity/check.md)

## History

- `01M5SWSK1HH2VMB5RJ89DJ6YBX` Task 9's plan decisions (PR #6, merge 8355823), accepted by the owner
