# busy-session-arbiter-ruling

source · `01M6506NXFBSVGWHHHJHENJ8J1`

## quotes

```json
[]
```

## raw_sha256

b03b667fc24882ba2751c30902adbc9294e93dff73cf94246c25555c83959a1c

## session

c315a867-1d73-4456-b5c5-01a18a012b41

## summary

Autonomous plan run on change set 01M61GW821FWSX05A9FC231N65. The owner was not present, so the interview question on how a message reaches a session that is busy working was raised to the arbiter. The arbiter ruled: a message to a busy session is written to its message file as today, and its pointer is returned to the agent as hook context at the next tool-call boundary (PostToolUse). The typed pointer at the prompt stays as the fallback when the session goes idle with no further tool call. No owner words were given, so no quotes.

## title

Message a busy session: arbiter ruling on delivery

## Links

### In

- derived-from [busy-session-message-via-tool-hook](../decision/busy-session-message-via-tool-hook.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
