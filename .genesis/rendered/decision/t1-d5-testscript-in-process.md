# t1-d5-testscript-in-process

decision · `01M5SWRRE5H7M2JZ83M05NC94Z`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: In-process, hermetic git (Recommended): `testscript.Main` registers `genesis` in the test binary, so it is never the one on $PATH. Setup isolates HOME and the git config with a fixed identity, and scripts live in `cmd/genesis/testdata/script/*.txtar`.

## statement

testscript runs in-process via `testscript.Main`, with a hermetic HOME and git config (fixed identity). Scripts live in `cmd/genesis/testdata/script/*.txtar`.

## wording

agent of milestone-1 task 1 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-1](../source/plan-task-1.md)
- derived-from [task-interviews](../source/task-interviews.md)

## History

- `01M5SWRR9E9G71HBT66132P6Z8` Task 1's plan decisions (PR #1, merge 5a4b716), accepted by the owner
