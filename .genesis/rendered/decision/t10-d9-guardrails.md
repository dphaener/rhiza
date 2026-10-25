# t10-d9-guardrails

decision · `01M5SWST79AB77JFKDRMQ01N7Y`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Yes, as listed (Recommended): No changes to internal/{op,changeset,model,entry,store,canon,id,key}: source add builds on api's runWrite. A small, isolated change to cmd/genesis/ctx.go (add a Home field), since task 5 may touch it too. Tests never touch the real ~/.claude, the network, or run `claude`. No AGENTS.md install, no genesis-reconcile skill, no evals. milestone-1.md is not edited.

## statement

Guardrails: no changes to `internal/{op,changeset,model,entry,store,canon,id,key}`; a small isolated `ctx.go` change; tests never touch the real `~/.claude` or the network, and never run `claude`; no `AGENTS.md`, no reconcile skill, no evals; `milestone-1.md` is untouched

## wording

agent of milestone-1 task 10 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-10](../source/plan-task-10.md)
- derived-from [task-interviews](../source/task-interviews.md)

## History

- `01M5SWSRAVRYJ1AZH6RT1GF122` Task 10's plan decisions (PR #8, merge 3e17d7b), accepted by the owner
