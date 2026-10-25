# t9-d4-check-base

decision · `01M5SWSKSYB1DRJ05F3KQ4PD1B`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: --base, default origin/HEAD (Recommended): `--base <rev>`; default the merge-base of HEAD with origin/HEAD, else main. Diff base→working tree on .genesis/log: only added files allowed. No git or no base: the rule is skipped with a warning, not a failure.

## statement

Base: `--base <rev>`; default merge-base of `HEAD` with the first of `origin/HEAD`, `origin/main`, `main` that resolves; diff base → working tree on `.genesis/log`, only added files allowed; no git or no base → a warning and the rule is skipped

## wording

agent of milestone-1 task 9 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-9](../source/plan-task-9.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [check](../capability/genesis/integrity/check.md)

## History

- `01M5SWSK1HH2VMB5RJ89DJ6YBX` Task 9's plan decisions (PR #6, merge 8355823), accepted by the owner
