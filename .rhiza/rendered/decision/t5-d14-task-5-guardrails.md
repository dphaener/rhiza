# t5-d14-task-5-guardrails

decision · `01M5SWS2MK868EN02QN3SQZ6TE`

## accepted_by

owner

## settled_by

Agent's options, accepted by the owner: (1) No edits to api/cmd (Recommended): internal/api and cmd/genesis stay untouched (task 9 works there in parallel; describe golden unchanged); only go.mod/go.sum change outside the new packages. (2) internal/model unchanged (Recommended): The projection engine and its semantics are reused as-is; the index never reimplements apply, order or conflict rules. (3) Never write .genesis (Recommended): The index only reads the log; it writes nothing under .genesis, and no cgo anywhere (CGO_ENABLED=0 build must work).

## statement

Guardrails: no edits to `internal/api` or `cmd/genesis`; `internal/model` unchanged; the index never writes under `.genesis`; no cgo.

## wording

agent of milestone-1 task 5 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-5](../source/plan-task-5.md)
- derived-from [task-interviews](../source/task-interviews.md)

## History

- `01M5SWS22WXZMST7XXR9ZKY4EE` Task 5's plan decisions (PR #7, merge ed8d292), accepted by the owner
