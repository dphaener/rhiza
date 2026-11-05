# workflow-new-verified-2

evidence · `01M6KPHA44QYV9KAV4KJ11H8VW` · actor [claude-code-21d15751](../actor/claude-code-21d15751.md)

## result

PASS

## statement

At d30ad3cf (attempt 2): go test -trimpath ./cmd/rhiza -run TestScripts/workflow- passed (workflow-new.txtar forks ship-change with --from, name line aside a fork of the fork is byte-identical and passes workflow check; the skeleton passes workflow check and run start --dry-run; an existing file gives exit 3 and stays byte-for-byte unchanged, including with --from and --json conflict; sub/x and x.toml give exit 2; --dry-run writes nothing; --json matches api call). TestWithNameDiffersOnlyInName and TestWithNameAddsMissing pass. Read the fix diff: the skeleton drops agent = "claude-code" and permission_mode, so nothing in internal/api/workflow_new*.go names an agent. A scratch build's skeleton passed workflow check, and like the built-in build workflow it asks for --permission-mode at run start.

## Links

### Out

- verifies [workflow-new](../capability/rhiza/workflows/workflow-authoring/workflow-new.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
