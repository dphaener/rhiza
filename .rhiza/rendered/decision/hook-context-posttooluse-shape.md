# hook-context-posttooluse-shape

decision · `01M65GY4HQ517P0S6XVJJBS0Q4`

## options

```json
[
  "PostToolUse hookSpecificOutput.additionalContext",
  "PreToolUse hook context",
  "type the pointer into the busy terminal"
]
```

## reasoning

The implementing agent's choice, accepted by the orchestrator acting as arbiter in review. Code: Adapter.ContextReply in internal/adapter/claudecode/claudecode.go prints exactly this shape for PostToolUse and returns nothing for any other event. Checked against the real agent in busy-session-hook-delivery-qa-1: with real Rex and real Claude Code, the pointer returned this way showed up once in the transcript as hook_additional_context on PostToolUse:Bash, and the agent read the message file and acted on it mid-task. The testscript evidence (busy-session-hook-delivery-verified) only drives fake hooks and checks the printed bytes.

## statement

Yes: Claude Code adds a PostToolUse hook's additional context to the agent's context before its next step. The hook prints {"hookSpecificOutput":{"hookEventName":"PostToolUse","additionalContext":"<text>"}} on stdout, and only for the PostToolUse event; other events return nothing.

## Links

### Out

- answers [hook-context-injected](../question/hook-context-injected.md)

### In

- decided-by [busy-session-hook-delivery](../capability/rhiza/orchestration/step-messages/busy-session-hook-delivery.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
