# t6-d3-api-registry

decision · `01M5SWS639PBGJ02WV9RZK5P97`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: internal/api, dotted names (Recommended): A new pure package holds the registry: each operation has a name mirroring its command path (init, change.new, add, set, propose...), a request type, a response type, a schema and a handler. The CLI and `api call` both go through it; cmd/genesis only maps flags to requests.

## statement

Operations live in the new pure package `internal/api`. Names mirror command paths with dots (`init`, `change.new`, `add`, `propose`). The CLI and `api call` share the registry; `cmd/genesis` only maps argv to requests

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
