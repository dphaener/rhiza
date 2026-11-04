# workflow-describe-qa-1

evidence · `01M6KK2F5VB8Q8566WB5YKKD9T` · actor [claude-code-a96577be](../actor/claude-code-a96577be.md)

## result

PASS

## statement

QA with a scratch build: workflow describe ship-change showed source builtin:ship-change, start build/build, includes build, review, and steps as <workflow>/<step> with kinds, on.* targets and roles with model ladders and prompts. A personal review shadowing repo and builtin showed both shadows lines. Editing the repo review.toml changed what ship-change's included review resolved to (shared resolver). describe left the file count under .git/rhiza and .rhiza unchanged (18 -> 18). --json was byte-identical to api call workflow.describe '{"workflow":"ship-change"}'

## Links

### Out

- verifies [workflow-describe](../capability/rhiza/workflows/workflow-authoring/workflow-describe.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
