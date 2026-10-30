# blocked-session-still-waits-qa-1

evidence · `01M65GTBVZHSDKS5F9TZJ96JWW` · actor [claude-code-311a7882](../actor/claude-code-311a7882.md)

## result

PASS

## statement

At 898cba6, built rhiza to a scratch path and, in a scratch repo, ran `rhiza session start --permission-mode manual` with real Rex and real Claude Code (session bdcf6c60), hooks calling the built binary. With the agent showing an AskUserQuestion menu (status blocked question), `rhiza session message ... --timeout 6s` waited, then failed with exit 3 ('did not reach its prompt within 6s', message in 0006.md). Nothing was queued as pending, and the text never appeared in the terminal: the menu stayed untouched.

## Links

### Out

- verifies [blocked-session-still-waits](../capability/rhiza/orchestration/step-messages/blocked-session-still-waits.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
