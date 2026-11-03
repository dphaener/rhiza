# claude-folder-trust

capability · `01M6FWF2YX1F6WAGBT9JDYEZT5` · parent [agent-adapters](../agent-adapters.md)

## acceptance

Tests with a hermetic HOME show that session run, before it starts the agent, leaves both paths with hasTrustDialogAccepted=true and every other key in .claude.json unchanged; that it honours CLAUDE_CONFIG_DIR and creates the file when it is missing; that a resumed start trusts the paths again; that --dry-run writes nothing; and that Adapter.Command still writes no file

## statement

Before every Claude Code start, the claudecode adapter marks the main checkout root and the session's working directory as trusted (hasTrustDialogAccepted=true) in ~/.claude.json, or $CLAUDE_CONFIG_DIR/.claude.json when set, through an optional adapter interface that session run calls, so a session in a new repo never stops at the folder-trust prompt

## Links

### Out

- derived-from [trust-fix-plan](../../../../source/trust-fix-plan.md)

### In

- verifies [claude-folder-trust-verified](../../../../evidence/claude-folder-trust-verified.md)
- verifies [claude-folder-trust-qa-1](../../../../evidence/claude-folder-trust-qa-1.md)
- verifies [pr-205-fix](../../../../evidence/pr-205-fix.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
