# workflow-overrides

capability · `01M5SWT9GPK3F87HZRV8V9XY8D` · parent [workflows](../workflows.md)

## acceptance

Tests with a hermetic HOME/XDG_CONFIG_HOME show a personal file shadowing a repo file and a built-in of the same name, a repo file shadowing a built-in, and XDG_CONFIG_HOME unset falling back to ~/.config; run start and a workflow step naming the shadowed workflow both pick the personal file, and the run journal records that file's path (builtin:<name> for a built-in)

## statement

A bare workflow name resolves to the personal file $XDG_CONFIG_HOME/rhiza/workflows/<name>.toml (~/.config/rhiza/workflows when XDG_CONFIG_HOME is unset), else the repo's .rhiza/workflows/<name>.toml, else the built-in, through one resolver that run start, workflow steps naming another workflow, and every workflow command share

## Links

### Out

- decided-by [built-in-workflows-by-name](../../../decision/built-in-workflows-by-name.md)
- decided-by [workflow-authoring-shape](../../../decision/workflow-authoring-shape.md)

### In

- depends-on [workflow-ls](workflow-authoring/workflow-ls.md)
- depends-on [workflow-describe](workflow-authoring/workflow-describe.md)
- depends-on [workflow-check](workflow-authoring/workflow-check.md)
- verifies [workflow-overrides-verified](../../../evidence/workflow-overrides-verified.md)
- verifies [workflow-overrides-qa-1](../../../evidence/workflow-overrides-qa-1.md)
- verifies [pr-249-workflow-authoring](../../../evidence/pr-249-workflow-authoring.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
