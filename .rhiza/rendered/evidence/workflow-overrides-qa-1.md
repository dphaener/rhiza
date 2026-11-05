# workflow-overrides-qa-1

evidence · `01M6KK1A8DJWSDNXVM5P77JGEA` · actor [claude-code-a96577be](../actor/claude-code-a96577be.md)

## result

PASS

## statement

QA with a scratch build of this branch in a scratch git repo, hermetic HOME: a personal myflow.toml under ~/.config/rhiza/workflows (XDG_CONFIG_HOME unset) shadowed the repo's .rhiza/workflows/myflow.toml in workflow ls; with XDG_CONFIG_HOME set, the $XDG_CONFIG_HOME/rhiza/workflows file won. Gave cmdonly a personal and a repo file, each a command step writing its layer name. run start --workflow cmdonly --foreground and a repo workflow 'outer' (a workflow step naming cmdonly) both ran the personal file (said.txt = layer-personal). Each run journal recorded workflow_path as the personal file's path (outer.toml for outer)

## Links

### Out

- verifies [workflow-overrides](../capability/rhiza/workflows/workflow-overrides.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
