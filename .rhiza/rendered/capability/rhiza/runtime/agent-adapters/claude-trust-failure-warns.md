# claude-trust-failure-warns

capability · `01M6FWFF8A7KTJAS9GXSK8CSHK` · parent [agent-adapters](../agent-adapters.md)

## acceptance

A test with a malformed .claude.json shows that session run starts the agent, reports a warning naming the file, and leaves the file byte-for-byte unchanged; the write goes through a temp file and a rename

## statement

When .claude.json can't be read, parsed or written, the session still starts, the failure appears in session run's warnings, and the file is never overwritten or left half-written

## Links

### Out

- derived-from [trust-fix-plan](../../../../source/trust-fix-plan.md)

### In

- verifies [claude-trust-failure-warns-verified](../../../../evidence/claude-trust-failure-warns-verified.md)
- verifies [claude-trust-failure-warns-qa-1](../../../../evidence/claude-trust-failure-warns-qa-1.md)
- verifies [pr-205-fix](../../../../evidence/pr-205-fix.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
