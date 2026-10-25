# t6-d14-guardrails

decision · `01M5SWS81R5ZBV76HFHB38BPPE`

## accepted_by

owner

## settled_by

Agent's options, accepted by the owner (several picked): changeset/op/store: additive only: No behaviour changes in existing packages; at most small exported helpers (task 4 is editing nearby). Recommended. / No reads, no index, no --view: Nothing from tasks 5, 7, 8, 9: no show/tree/ls/why, no SQLite, no render/check. Recommended. / No new dependencies: Standard library plus what go.mod already has. Recommended. / Hermetic tests: No network, no real HOME/git config, no ~/bin binary; clock and IDs fixed through test-only env vars. Recommended.

## statement

Guardrails: `changeset`/`op`/`store` additive only; no reads, index, `--view`, render or check; no new dependencies; hermetic tests

## wording

agent of milestone-1 task 6 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-6](../source/plan-task-6.md)
- derived-from [task-interviews](../source/task-interviews.md)

## History

- `01M5SWS5QKWPQ3FS357EAV60D2` Task 6's plan decisions (PR #5, merge 6d5c61e), accepted by the owner
