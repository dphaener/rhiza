# workflow-guide-topic-qa-1

evidence · `01M6KK4AT87AEA1T9RA2S5PPN5` · actor [claude-code-a96577be](../actor/claude-code-a96577be.md)

## result

PASS

## statement

QA with a scratch build: rhiza guide ls lists the workflows topic. rhiza guide workflows covers the personal/repo/built-in layers, every step kind and its keys (agent, wait, command, leaves, checks, signoff, arbiter, workflow), outcomes, the reserved targets done/failed, includes as <workflow>/<step>, and the ls/new/describe/check/run start --dry-run loop. I followed that loop by hand in a scratch repo and it worked as written

## Links

### Out

- verifies [workflow-guide-topic](../capability/rhiza/workflows/workflow-authoring/workflow-guide-topic.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
