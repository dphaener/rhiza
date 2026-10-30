# idle-typed-pointer-kept

capability · `01M65074BASR8PYQHQSC7B7WJ7` · parent [step-messages](../step-messages.md)

## acceptance

The existing session-message testscript still passes unchanged: an idle session gets the pointer typed once and then Enter, and a message delivered by hook context is not typed again

## statement

A message to a session that goes idle with no further tool call still has its pointer typed at the prompt, once, then Enter, as it does today

## Links

### Out

- decided-by [busy-session-message-via-tool-hook](../../../../decision/busy-session-message-via-tool-hook.md)

### In

- verifies [idle-typed-pointer-kept-verified](../../../../evidence/idle-typed-pointer-kept-verified.md)
- verifies [idle-typed-pointer-kept-qa-1](../../../../evidence/idle-typed-pointer-kept-qa-1.md)
- verifies [pr-105-message-a-busy-session](../../../../evidence/pr-105-message-a-busy-session.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
