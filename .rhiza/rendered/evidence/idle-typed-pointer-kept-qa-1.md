# idle-typed-pointer-kept-qa-1

evidence · `01M65GTB19MEY77B26DH563S9R` · actor [claude-code-311a7882](../actor/claude-code-311a7882.md)

## result

PASS

## statement

At 898cba6, built rhiza to a scratch path and, in a scratch repo, ran `rhiza session start --permission-mode manual` with real Rex and real Claude Code (session bdcf6c60), hooks calling the built binary. A message to the idle session (0003) was typed at the prompt at once (Waited: 0s; a typed user prompt in the transcript). A message sent while the agent was writing its final reply after its last tool call (0004) was queued as pending, then typed at the prompt once when the Stop hook fired (one typed prompt, no hook context for it). The agent then ran `echo GOT-IDLE-FALLBACK`. `go test ./cmd/rhiza -run TestScripts/session-message -count=1`: ok.

## Links

### Out

- verifies [idle-typed-pointer-kept](../capability/rhiza/orchestration/step-messages/idle-typed-pointer-kept.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
