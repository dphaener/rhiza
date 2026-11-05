# workflow-describe

capability · `01M6HNBFTE249W98NA8YV0PKCP` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

A testscript describes a built-in (source builtin:<name>), a repo override that shadows it, and a file that includes another workflow (its steps shown as <workflow>/<step>); --json prints the same bytes as rhiza api call workflow.describe; no journal or run is written

## statement

rhiza workflow describe <name|path> resolves the workflow as run start does and shows, without starting anything, where it resolved from and the layers it shadows, its start step and each step in order with its kind, on.<outcome> targets, roles with their model ladders and prompts, and the workflows it includes

## Links

### Out

- depends-on [workflow-overrides](../workflow-overrides.md)
- decided-by [workflow-authoring-shape](../../../../decision/workflow-authoring-shape.md)

### In

- depends-on [workflow-guide-topic](workflow-guide-topic.md)
- verifies [workflow-describe-verified](../../../../evidence/workflow-describe-verified.md)
- verifies [workflow-describe-qa-1](../../../../evidence/workflow-describe-qa-1.md)
- verifies [workflow-describe-verified-2](../../../../evidence/workflow-describe-verified-2.md)
- verifies [pr-249-workflow-authoring](../../../../evidence/pr-249-workflow-authoring.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
