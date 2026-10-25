# t6-d10-apply-file-format

decision · `01M5SWS79XT7Q0P33V4WE96472`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: A list of write requests (Recommended): JSON: {"writes":[{"op":"add",...},{"op":"link",...}]}, each item the same request shape as the matching command's operation. Refs may name nodes added earlier in the file by handle. All of it becomes one ops entry, or nothing is written. '-' reads stdin.

## statement

`apply` reads `{"writes":[{"op":"add",...},...]}` made of the command request shapes. Refs may name nodes added earlier in the file. The whole file becomes one ops entry, or nothing is written. `-` reads stdin

## wording

agent of milestone-1 task 6 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-6](../source/plan-task-6.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [apply-files](../capability/genesis/change-sets/apply-files.md)

## History

- `01M5SWS5QKWPQ3FS357EAV60D2` Task 6's plan decisions (PR #5, merge 6d5c61e), accepted by the owner
