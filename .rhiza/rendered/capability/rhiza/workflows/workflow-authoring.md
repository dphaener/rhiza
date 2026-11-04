# workflow-authoring

capability · `01M5SWT9JXJVYF0CXW2QWBP2JY` · parent [workflows](../workflows.md)

## statement

Genesis tells an agent how to write a custom workflow: a guide topic, workflow describe, new and check.

## Links

### Out

- decided-by [explain-custom-workflows](../../../decision/explain-custom-workflows.md)
- decided-by [workflow-writing-tools](../../../decision/workflow-writing-tools.md)

### In

- verifies [pr-249-workflow-authoring](../../../evidence/pr-249-workflow-authoring.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
