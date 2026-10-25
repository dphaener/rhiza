# t1-d2-cmd-internal-layout

decision · `01M5SWRRBWGNYS453JGZENTBSS`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: cmd/genesis + internal/ (Recommended): `cmd/genesis/main.go` stays thin and domain packages go under `internal/`, which matches 'a pure domain with a thin CLI'. Installs with `go install .../cmd/genesis`.

## statement

`cmd/genesis/main.go` is the thin entry point. Domain packages go under `internal/` (none created now).

## wording

agent of milestone-1 task 1 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-1](../source/plan-task-1.md)
- derived-from [task-interviews](../source/task-interviews.md)

## History

- `01M5SWRR9E9G71HBT66132P6Z8` Task 1's plan decisions (PR #1, merge 5a4b716), accepted by the owner
