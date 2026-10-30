# busy-session-hook-delivery

capability · `01M65072NYFJ9RDZPJMXRKSWZ3` · parent [step-messages](../step-messages.md)

## acceptance

A testscript sends a message to a session whose status is working: the command returns at once, writes the message file, types nothing into the Rex tab, and the session's next tool-call hook output carries the pointer line naming that file; a second tool-call hook output for the same session does not repeat it

## statement

A message to a session that is busy working is delivered mid-task: its pointer is returned to the agent as hook context at its next tool call, and the send command does not wait

## Links

### Out

- decided-by [busy-session-message-via-tool-hook](../../../../decision/busy-session-message-via-tool-hook.md)
- decided-by [hook-context-posttooluse-shape](../../../../decision/hook-context-posttooluse-shape.md)

### In

- verifies [busy-session-hook-delivery-verified](../../../../evidence/busy-session-hook-delivery-verified.md)
- verifies [busy-session-hook-delivery-qa-1](../../../../evidence/busy-session-hook-delivery-qa-1.md)
- verifies [busy-session-hook-delivery-verified-2](../../../../evidence/busy-session-hook-delivery-verified-2.md)
- verifies [pr-105-message-a-busy-session](../../../../evidence/pr-105-message-a-busy-session.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
