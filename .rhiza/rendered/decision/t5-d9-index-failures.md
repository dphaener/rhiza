# t5-d9-index-failures

decision · `01M5SWS2EJMPBDNMC529NZWBY4`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Rebuild; wait on locks (Recommended): Schema version in PRAGMA user_version: mismatch or a corrupt/unopenable file → delete and rebuild from the log. Catch-up runs in BEGIN IMMEDIATE with WAL and a busy timeout (~5 s), re-checking freshness once it holds the lock; a timeout is an error, never a silent half-update.

## statement

Failures: schema version in `PRAGMA user_version`; mismatch or corrupt/unopenable → delete and rebuild. Catch-up in `BEGIN IMMEDIATE`, WAL, busy timeout 5 s, freshness re-checked under the lock; a timeout is an error.

## wording

agent of milestone-1 task 5 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-5](../source/plan-task-5.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [fast-reads](../capability/rhiza/views/fast-reads.md)

## History

- `01M5SWS22WXZMST7XXR9ZKY4EE` Task 5's plan decisions (PR #7, merge ed8d292), accepted by the owner
