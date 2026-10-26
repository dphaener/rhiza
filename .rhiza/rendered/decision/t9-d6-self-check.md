# t9-d6-self-check

decision · `01M5SWSM8G4KA5W45E6XFBSYK0`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: make check runs it (Recommended): `genesis init` here and commit .genesis/ (repo.json, .gitignore, rendered/); `make check` builds a scratch binary in $TMPDIR and runs `check` with the same-commit build; ci.yml gets fetch-depth 0. `ci install` is not run on this repo (its install step would fetch a released binary, not this commit).

## statement

Genesis checks itself: `.genesis/` initialised and committed here (repo.json, .gitignore, rendered/); `make check` builds a scratch binary in a temp dir and runs `check`; `ci.yml` checkout uses `fetch-depth: 0`; `ci install` is not run on this repo

## wording

agent of milestone-1 task 9 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-9](../source/plan-task-9.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [self-check](../capability/rhiza/integrity/self-check.md)

## History

- `01M5SWSK1HH2VMB5RJ89DJ6YBX` Task 9's plan decisions (PR #6, merge 8355823), accepted by the owner
