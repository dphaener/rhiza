# inbox-across-projects-verified

evidence · `01M6SR90VKQ4HB59351A9M4QNC` · actor [claude-code-e025624b](../actor/claude-code-e025624b.md)

## result

PASS

## statement

At d82b495e ran go test -trimpath ./cmd/rhiza -run TestScripts/inbox-projects -count=1 (PASS): a root of two projects alpha and beta, each with one fixture run journal; rhiza inbox --json at the root lists both items with a project field each (alpha question, beta signoff), human output shows a project: line; inbox open of alpha's question from the root (and with --all from inside beta) focuses its tab through the fake rex; inbox approve of beta's sign-off from the root reaches beta's run SIGNOFFRUN01; from inside alpha, or with --project alpha/beta, only that project's item and no project field; unregistered --project and --project with --all are usage. Also passed every TestScripts/inbox* and tui-inbox script, describe-golden, and go test ./internal/api -run Inbox|inbox. Read the diff: inboxScopes picks the projects, mergeInbox tags and dedupes machine-wide items, findInboxItem returns the item's project ctx and approve/request-changes run in it.

## Links

### Out

- verifies [inbox-across-projects](../capability/rhiza/orchestration/inbox-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
