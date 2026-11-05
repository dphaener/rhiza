# workflow-ls-qa-1

evidence · `01M6KK1WRXMNMDEE4K6XTC67MC` · actor [claude-code-a96577be](../actor/claude-code-a96577be.md)

## result

PASS

## statement

QA with a scratch build: workflow ls in a scratch repo listed the 4 built-ins as builtin:<name>. After adding repo and personal files it showed each name once with its winning layer and path, plus 'shadows repo <path>' and 'shadows builtin builtin:review' lines for review (personal > repo > builtin). workflow ls --json printed the same bytes as rhiza api call workflow.ls '{}' (cmp)

## Links

### Out

- verifies [workflow-ls](../capability/rhiza/workflows/workflow-authoring/workflow-ls.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
