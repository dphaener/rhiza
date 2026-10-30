# blocked-session-still-waits-verified

evidence · `01M65GBXWAPS1MS6Y1KYF3N545` · actor [claude-code-8ab0a4dc](../actor/claude-code-8ab0a4dc.md)

## result

PASS

## statement

At 5957cc2, ran go test ./cmd/rhiza -run 'TestScripts/session-message$' -count=1: PASS. The script (unchanged by this change set) opens an AskUserQuestion menu, sends a message, checks the file is written but no send appears in the fake rex log while blocked, then after the Stop hook sees the command finish with exactly one pointer send and one Enter. Also ran go test ./internal/api -run 'TestQueuePointer|TestTypePointer' -count=1: PASS, incl. new TestQueuePointerLeavesNonWorkingSessions (blocked/question, blocked/permission, idle keep no pending pointer). Read the diff: queuePointer in internal/api/session_message.go only queues when the resolved state is working, so a blocked session falls through to typePointer's wait for the prompt.

## Links

### Out

- verifies [blocked-session-still-waits](../capability/rhiza/orchestration/step-messages/blocked-session-still-waits.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
