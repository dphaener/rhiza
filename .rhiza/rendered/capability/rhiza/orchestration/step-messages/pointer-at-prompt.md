# pointer-at-prompt

capability · `01M5TYN3RQJHQJKFATJ8K4ZMQW` · parent [step-messages](../step-messages.md)

## acceptance

With a fake status of "menu open" no keys are sent, and they are sent once the status is idle.

## statement

A pointer that is typed goes in only while the session waits at its prompt, never over an open menu or question: a blocked session is waited on until it reaches its prompt, while a working session's pointer goes to its next tool call as hook context instead (busy-session-hook-delivery).

## Links

### Out

- decided-by [status-from-hooks-first](../../../../decision/status-from-hooks-first.md)

### In

- verifies [pointer-at-prompt-script](../../../../evidence/pointer-at-prompt-script.md)
- verifies [pr-105-message-a-busy-session](../../../../evidence/pr-105-message-a-busy-session.md)

## History

- `01M5SWTQ11MR7CP75QWPYFBQFQ` Genesis amends a running step and sends it a message.
- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
