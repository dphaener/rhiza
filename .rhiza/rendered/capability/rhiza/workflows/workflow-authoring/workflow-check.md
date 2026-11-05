# workflow-check

capability · `01M6HNBKTRQ2H5SNEPMXYSQVNX` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

A testscript runs workflow check on each built-in and on the internal/workflow testdata files: valid ones exit 0, bad-outcome, missing-target and unknown-kind each exit non-zero naming the step and problem, a missing include is reported, and --json prints the same bytes as rhiza api call workflow.check

## statement

rhiza workflow check <name|path> parses and validates a workflow with includes expanded, as run start would, and reports every problem it finds with the step it is in, or that the workflow is valid

## Links

### Out

- depends-on [workflow-overrides](../workflow-overrides.md)
- decided-by [workflow-writing-tools](../../../../decision/workflow-writing-tools.md)

### In

- depends-on [workflow-new](workflow-new.md)
- depends-on [workflow-guide-topic](workflow-guide-topic.md)
- verifies [workflow-check-verified](../../../../evidence/workflow-check-verified.md)
- verifies [workflow-check-qa-1](../../../../evidence/workflow-check-qa-1.md)
- verifies [workflow-check-verified-2](../../../../evidence/workflow-check-verified-2.md)
- verifies [pr-249-workflow-authoring](../../../../evidence/pr-249-workflow-authoring.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
