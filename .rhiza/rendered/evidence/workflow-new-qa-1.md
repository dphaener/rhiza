# workflow-new-qa-1

evidence · `01M6KK3RV9XQ5TM5QJ3T1WDHN5` · actor [claude-code-a96577be](../actor/claude-code-a96577be.md)

## result

PASS

## statement

QA with a scratch build: workflow new sc --from ship-change wrote .rhiza/workflows/sc.toml, which diffs from internal/workflow/builtin/ship-change.toml only in the name line and passes workflow check. workflow new myflow wrote a commented skeleton that passes check and runs under run start --workflow myflow --dry-run ('would start', nothing written). Re-running new on an existing name exited 3 with the file unchanged; a/b and x.toml exited 2; --dry-run printed 'would write' and created no file

## Links

### Out

- verifies [workflow-new](../capability/rhiza/workflows/workflow-authoring/workflow-new.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
