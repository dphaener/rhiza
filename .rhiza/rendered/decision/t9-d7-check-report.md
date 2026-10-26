# t9-d7-check-report

decision · `01M5SWSMFJH43SX8T28VMX4MCS`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Report + exit 1 (Recommended): The response is always the full report (findings with severity, code, path/IDs, message, next), human or --json; exit 0 with only warnings, 1 with any failure. The frame gains a small hook so a response can set a non-zero exit while still printing.

## statement

`check` always returns the full report (findings with severity, rule, path/IDs, message, next) as human text or `--json`; exit 0 when only warnings, 1 with any failure; the frame gains a hook for a response-set exit code

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
