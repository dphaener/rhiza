# t1-d9-version-stamping

decision · `01M5SWRRH4RFHAJ90KPR4S8JJ2`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: ldflags, fallback to buildinfo (Recommended): Makefile stamps `git describe --always --dirty` via -ldflags. Unstamped builds (`go install`, tests) fall back to debug.ReadBuildInfo, then "dev". Tests assert on the shape, not the exact value.

## statement

The version is stamped by the Makefile via `-ldflags "-X main.version=$(git describe --always --dirty)"`. Unstamped builds fall back to `debug.ReadBuildInfo` and then `"dev"`. Tests assert shape, not value.

## wording

agent of milestone-1 task 1 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-1](../source/plan-task-1.md)
- derived-from [task-interviews](../source/task-interviews.md)

## History

- `01M5SWRR9E9G71HBT66132P6Z8` Task 1's plan decisions (PR #1, merge 5a4b716), accepted by the owner
