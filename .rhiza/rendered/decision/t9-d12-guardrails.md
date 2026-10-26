# t9-d12-guardrails

decision · `01M5SWSNPG60B1RXYJWQFSTPVJ`

## accepted_by

owner

## settled_by

Agent's options, accepted by the owner (several picked): Core pkgs untouched (Recommended): No edits to id/canon/entry/store/op/key/changeset/model; new code goes in new packages (e.g. internal/render, internal/check) plus new api ops. / Minimal api/frame edits (Recommended): workspace.go and write code untouched (task 5 runs in parallel); the only frame change is the exit-code hook. / No new dependencies (Recommended): Standard library and existing modules only; git through exec in cmd/genesis, passed in via Ctx. / Nothing from milestone 2 (Recommended): No anchors, fingerprints, trailers, drift, retired-endpoint flags, relates-to watching, auto-render, status, or reads.

## statement

Guardrails: core packages untouched; `workspace.go` and write code untouched; only frame change is the exit hook; no new dependencies; nothing from milestone 2

## wording

agent of milestone-1 task 9 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-9](../source/plan-task-9.md)
- derived-from [task-interviews](../source/task-interviews.md)

## History

- `01M5SWSK1HH2VMB5RJ89DJ6YBX` Task 9's plan decisions (PR #6, merge 8355823), accepted by the owner
