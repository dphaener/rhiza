# busy-session-hook-delivery-qa-1

evidence · `01M65GTA6PGPJG6XF35S1J2VEF` · actor [claude-code-311a7882](../actor/claude-code-311a7882.md)

## result

PASS

## statement

At 898cba6, built rhiza to a scratch path and, in a scratch repo, ran `rhiza session start --permission-mode manual` with real Rex and real Claude Code (session bdcf6c60), hooks calling the built binary. While the agent ran a `sleep 4` loop (status working), `rhiza session message` returned in 0.05s with Waited: 0s and wrote 0002.md; the record held the pointer as pending. The agent's transcript shows it arrived once as hook_additional_context on PostToolUse:Bash, never as a typed prompt; the agent read the file, ran `echo GOT-BUSY-MESSAGE` and stopped the loop. Later PostToolUse hooks did not repeat it, and the pending list was empty.

## Links

### Out

- verifies [busy-session-hook-delivery](../capability/rhiza/orchestration/step-messages/busy-session-hook-delivery.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
