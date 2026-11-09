# message-says-why-waiting-verified

evidence · `01M6ZSKW2D1RNN8ZK4WNP3H6A9` · actor [claude-code-f7fb8eb8](../actor/claude-code-f7fb8eb8.md)

## result

PASS

## statement

Verifier (Opus) at c7fc3245: go test -trimpath -count=1 ./internal/api -run 'TestSessionMessage|TestTypePointer|TestQueuePointer' passes, including TestSessionMessageSaysWhyWaiting (a permission-blocked session gets exactly one stderr line with 'blocked on a permission: Bash' and 'rhiza inbox open'; the timeout error names 'last status: blocked on a permission: Bash') and TestSessionMessageSaysWhenStatusChanges (blocked then working mid-wait gives exactly two lines, the second ': working;' with 'waiting up to 1s'; the timeout names 'last status: working'). With GOTMPDIR in a scratch dir and RHIZA_ROOT/CHANGE/RUN unset, go test ./cmd/rhiza -run TestScripts/(session-message|session-message-busy|session-message-idle-fallback|session-hook-message|run-status-messages) passes; the leaf's commits touch no testscript or golden, so stdout and --json are unchanged. Diff read: typePointer prints only when the resolved status changes, to ctx.Stderr, and defaultMessageTimeout is still 10 * time.Minute.

## Links

### Out

- verifies [message-says-why-waiting](../capability/rhiza/orchestration/step-messages/message-says-why-waiting.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
