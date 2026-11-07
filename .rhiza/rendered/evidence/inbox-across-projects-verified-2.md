# inbox-across-projects-verified-2

evidence · `01M6SWQH9592J9MB6C4A3DHPSG` · actor [claude-code-d6355a26](../actor/claude-code-d6355a26.md)

## result

PASS

## statement

Verified at 2b09cb8c (attempt 2). Read the diff: inboxScopes now returns the caller's own scope when the wd holds an initialised model (projectInitialised(ctx.Wd)) unless --project/--all, matching roadmap/health/daemon status (the review's quality f1). Ran go test -trimpath ./cmd/rhiza -run TestScripts/inbox (ok): inbox-projects.txtar shows at a two-project root inbox --json merged with project alpha and beta on each item, open focusing alpha's question from the root and not from inside beta without --all, approve reaching beta's sign-off run from the root, --project and an in-project wd giving only that project's items with no project field, and the new case where a root with its own model reads its own empty inbox while --all still covers alpha and beta. go test -trimpath ./internal/api -run Inbox: ok.

## Links

### Out

- verifies [inbox-across-projects](../capability/rhiza/orchestration/inbox-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
