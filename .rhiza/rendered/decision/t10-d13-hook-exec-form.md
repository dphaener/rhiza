# t10-d13-hook-exec-form

decision · `01M5SWSV1ZAMDFN9EZBG2TFZ53`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Exec form `genesis status --hook` via PATH (Recommended): {command: "genesis", args: ["status","--hook"], timeout: 10}, no matcher (so it fires on startup, resume, clear and compact). Needs `genesis` on PATH, which `make install` provides, and survives rebuilds.

## statement

Hook: `{"type":"command","command":"genesis","args":["status","--hook"],"timeout":10}`, no matcher, resolved through PATH

## wording

agent of milestone-1 task 10 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-10](../source/plan-task-10.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [session-start](../capability/rhiza/guide/session-start.md)

## History

- `01M5SWSRAVRYJ1AZH6RT1GF122` Task 10's plan decisions (PR #8, merge 3e17d7b), accepted by the owner
