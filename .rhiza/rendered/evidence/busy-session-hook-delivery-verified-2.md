# busy-session-hook-delivery-verified-2

evidence · `01M65HEAQXPJR137ATJ0BDQTCF` · actor [claude-code-6cbe1f78](../actor/claude-code-6cbe1f78.md)

## result

PASS

## statement

At commit 323a663 (fix attempt 2, which records the Decision hook-context-posttooluse-shape answering hook-context-injected; no code change since f858406 for this leaf's path), I ran go test ./cmd/rhiza -run 'TestScripts/(session-message-busy|session-hook-message|session-message-idle-fallback)$' -count=1: all PASS. session-message-busy.txtar drives a session to status working via a PostToolUse hook, then rhiza session message returns at once (Waited: 0s), writes .rhiza/local/messages/<sid>/0001.md with the text, logs no rex send; the next PostToolUse hook prints hookSpecificOutput with hookEventName PostToolUse and additionalContext naming that file; a second hook prints nothing and still no send. Read the diff of internal/api/session_message.go (queuePointer before typePointer) and session_run.go (deliverPending clears Pending once). make check exits 0. The Decision is linked decided-by from the leaf and answers the question.

## Links

### Out

- verifies [busy-session-hook-delivery](../capability/rhiza/orchestration/step-messages/busy-session-hook-delivery.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
