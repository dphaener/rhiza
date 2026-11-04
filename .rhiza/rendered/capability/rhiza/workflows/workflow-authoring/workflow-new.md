# workflow-new

capability · `01M6HNBHR1P4N0M34NBQD7V10H` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

A testscript forks the built-in ship-change with --from (the file differs from the built-in only in name and passes workflow check), writes a skeleton that passes workflow check and runs under run start --dry-run, exits 3 leaving an existing file byte-for-byte unchanged, exits 2 on a name with a / or .toml, and --dry-run writes nothing

## statement

rhiza workflow new <name> [--from <workflow>] writes .rhiza/workflows/<name>.toml: a copy of the resolved workflow's source with its name set to <name>, or without --from a minimal commented skeleton; it takes a bare name only, refuses an existing file with a conflict, validates the result as workflow check does and prints the path it wrote

## Links

### Out

- depends-on [workflow-check](workflow-check.md)
- decided-by [workflow-authoring-shape](../../../../decision/workflow-authoring-shape.md)

### In

- depends-on [workflow-guide-topic](workflow-guide-topic.md)
- verifies [workflow-new-verified](../../../../evidence/workflow-new-verified.md)
- verifies [workflow-new-qa-1](../../../../evidence/workflow-new-qa-1.md)
- verifies [workflow-new-verified-2](../../../../evidence/workflow-new-verified-2.md)
- verifies [pr-249-workflow-authoring](../../../../evidence/pr-249-workflow-authoring.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
