# reopen-format-three

decision · `01M5WT1NX9GMS5QQJ5ZNS44ZK3`

## statement

The reopen entry is format 3: only reopen entries carry v 3, so older readers fail with ErrNewerFormat and formats 1 and 2 stay byte-identical.

## wording

the implementing agent's choice, accepted by the orchestrator acting as arbiter

## Links

### In

- decided-by [reopen-format](../capability/rhiza/change-sets/reopen/reopen-format.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
