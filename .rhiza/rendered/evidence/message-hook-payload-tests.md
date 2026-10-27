# message-hook-payload-tests

evidence · `01M5XB6YEMHF5PG7Z7EGT39G7E` · actor [claude-code-d3710d3d](../actor/claude-code-d3710d3d.md)

## result

`go test ./internal/adapter/claudecode -run 'HookMessage|HookStatusUnchangedBySendMessage|SettingsOnePostToolUseGroupWithoutMatcher' -count=1` passed at 1466643

## statement

rhiza session hook sends a message payload for a tool's send-message call and leaves status unchanged.

## status

verified

## tests

TestHookMessage, TestHookMessageNoMessage, TestHookStatusUnchangedBySendMessage, TestSettingsOnePostToolUseGroupWithoutMatcher

## Links

### Out

- verifies [message-hook-payload](../capability/rhiza/runtime/session-journal/message-hook-payload.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
