# workflow-check-qa-1

evidence · `01M6KK379M4TKYNH0NQMZ3DQFT` · actor [claude-code-a96577be](../actor/claude-code-a96577be.md)

## result

PASS

## statement

QA with a scratch build: workflow check on build, plan, review and ship-change and on testdata/valid.toml exited 0 ('is valid'). bad-outcome, missing-target and unknown-kind each exited 1 naming the step and the problem. A file with two bad steps reported both, in the same bytes across 8 runs. A workflow step naming a missing workflow 'ghost' was reported. --json was byte-identical to api call workflow.check

## Links

### Out

- verifies [workflow-check](../capability/rhiza/workflows/workflow-authoring/workflow-check.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
