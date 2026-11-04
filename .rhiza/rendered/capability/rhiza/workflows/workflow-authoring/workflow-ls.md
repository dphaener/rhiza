# workflow-ls

capability · `01M6HNBDN6QV991BT9C284SA3H` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

A testscript with a personal file, a repo file and the built-ins shows each name once with its winning layer and path, marks the shadowed layers, and --json prints the same bytes as rhiza api call workflow.ls

## statement

rhiza workflow ls lists every workflow a bare name can reach: personal, repo and built-in, each with the layer and path it comes from and which ones it shadows

## Links

### Out

- depends-on [workflow-overrides](../workflow-overrides.md)
- decided-by [workflow-authoring-shape](../../../../decision/workflow-authoring-shape.md)

### In

- depends-on [workflow-guide-topic](workflow-guide-topic.md)
- verifies [workflow-ls-verified](../../../../evidence/workflow-ls-verified.md)
- verifies [workflow-ls-qa-1](../../../../evidence/workflow-ls-qa-1.md)
- verifies [pr-249-workflow-authoring](../../../../evidence/pr-249-workflow-authoring.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
