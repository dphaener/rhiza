# t6-d12-error-shape

decision · `01M5SWS7NABA97JA9ZGPZK8F4M`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: JSON on stdout with --json (Recommended): With --json, stdout gets {"error":{"code":"not_found","message":"...","next":"genesis change new \"...\""}} and nothing else, so an agent parses one stream. Without it, stderr gets `genesis: message` plus a `next: ...` line. Exit codes 0/1/2/3/4 either way.

## statement

Errors with `--json`: stdout gets `{"error":{"code","message","next"}}` and nothing else. Without `--json`: stderr gets `genesis: <message>` and a `next: <command>` line. Exit codes are 0/1/2/3/4 either way

## wording

agent of milestone-1 task 6 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-6](../source/plan-task-6.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [command-schemas](../capability/rhiza/guide/command-schemas.md)

## History

- `01M5SWS5QKWPQ3FS357EAV60D2` Task 6's plan decisions (PR #5, merge 6d5c61e), accepted by the owner
