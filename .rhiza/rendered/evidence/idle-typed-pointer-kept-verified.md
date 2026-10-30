# idle-typed-pointer-kept-verified

evidence · `01M65GDX306CRDY2Y6M4EEANPT` · actor [claude-code-2cda637b](../actor/claude-code-2cda637b.md)

## result

PASS

## statement

At b98ed19, git diff be12967..HEAD shows cmd/rhiza/testdata/script/session-message.txtar unchanged, and go test ./cmd/rhiza -run 'TestScripts/(session-message|session-message-busy|session-message-idle-fallback|session-hook-message)$' -count=5 passed. session-message.txtar asserts that an idle session gets the pointer typed once, then one Enter, with exactly 2 sends. The new session-message-idle-fallback.txtar checks three cases: a message queued while the session is working is typed once with Enter when Stop arrives with no tool call between; a later Stop types nothing; and a message already carried by a PostToolUse hook (0002.md in stdout) is not typed again (sends stay at 2). The diff adds typeIdlePending in internal/api/session_run.go. Under one store Update it takes the pending pointers only when the session resolves to idle, and clears them before typing, so none is typed twice. go test ./... passed on two reruns (one earlier run failed on the timing-sensitive read-path test). make check's self-check fails only on rendered actor files for sessions recorded after this commit, not on this leaf's code.

## Links

### Out

- verifies [idle-typed-pointer-kept](../capability/rhiza/orchestration/step-messages/idle-typed-pointer-kept.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
