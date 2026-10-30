# busy-session-hook-delivery-verified

evidence · `01M65G4DMMNN1RDEXMW6JNTYAV` · actor [claude-code-7d681652](../actor/claude-code-7d681652.md)

## result

PASS

## statement

At f858406 ran go test ./cmd/rhiza -run 'TestScripts/session-message' -count=1: ok. session-message-busy.txtar drives a session to working via a PostToolUse hook, then session message --timeout 1s returns at once (Waited: 0s), writes .rhiza/local/messages/<sid>/0001.md with the text, and the fake rex log has no send line; the next PostToolUse hook prints hookSpecificOutput {hookEventName PostToolUse, additionalContext 'rhiza: message from test: read .rhiza/local/messages/<sid>/0001.md and act on it.'}, and the following PostToolUse hook prints nothing. Read the diff: session.message queues the pointer in Record.Pending only when the resolved state is working, else types it as before; sessionHook's deliverPending drains Pending once on an event the adapter's ContextReply accepts (PostToolUse for Claude Code). gofmt, vet and go test ./... pass apart from TestReadBudget, a timing check that failed under full-suite load and passes alone (also passes on the parent commit); make self-check fails only on stale rendered actor files for live sessions outside the commit.

## Links

### Out

- verifies [busy-session-hook-delivery](../capability/rhiza/orchestration/step-messages/busy-session-hook-delivery.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
