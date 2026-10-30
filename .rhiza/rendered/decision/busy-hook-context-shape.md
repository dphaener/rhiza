# busy-hook-context-shape

decision · `01M65F4R2H97J7YHPZAHQ14T3H`

## statement

An agent adds the context a PostToolUse hook returns to its context at the next tool call; the hook prints JSON {"hookSpecificOutput":{"hookEventName":"PostToolUse","additionalContext":"..."}}. The build checks this with a fake hook before it ships

## wording

Arbiter answer: the documented Claude Code hook shape; the build verifies it, since the arbiter asked for a check before build

## Links

### Out

- answers [hook-context-injected](../question/hook-context-injected.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
