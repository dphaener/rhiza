# workflow-guide-topic

capability · `01M6HNBNSSMC77JWXF4Q21GM0G` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

rhiza guide ls lists the workflows topic; a guide test checks that every rhiza command the topic names resolves and that it names every step kind the workflow package defines

## statement

rhiza guide workflows is a guide topic that tells an agent how to write a custom workflow: the resolution layers, the step kinds and their keys, outcomes and reserved targets, includes, and the loop of workflow new, describe, check and run start --dry-run

## Links

### Out

- depends-on [workflow-ls](workflow-ls.md)
- depends-on [workflow-describe](workflow-describe.md)
- depends-on [workflow-new](workflow-new.md)
- depends-on [workflow-check](workflow-check.md)
- decided-by [explain-custom-workflows](../../../../decision/explain-custom-workflows.md)

### In

- verifies [workflow-guide-topic-verified](../../../../evidence/workflow-guide-topic-verified.md)
- verifies [workflow-guide-topic-qa-1](../../../../evidence/workflow-guide-topic-qa-1.md)
- verifies [pr-249-workflow-authoring](../../../../evidence/pr-249-workflow-authoring.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
