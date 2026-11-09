# message-says-why-waiting-qa-1

evidence · `01M701V7CG9R26AN330PRNHZFG` · actor [claude-code-1b20f401](../actor/claude-code-1b20f401.md)

## result

PASS

## statement

QA cycle 1, built binary: rhiza session message <sid> 'Use blue' --timeout 3s to a session whose hook reported a permission prompt printed one stderr line 'waiting for session <sid>: blocked on a permission: Claude needs your permission to use Bash; answer it with rhiza inbox open (or rhiza session show <sid>); waiting up to 3s' and after 3.1s failed (exit 3) with 'did not reach its prompt within 3s (last status: blocked on a permission: ...)', next rhiza session show. Scoped api TestSessionMessage passes.

## Links

### Out

- verifies [message-says-why-waiting](../capability/rhiza/orchestration/step-messages/message-says-why-waiting.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
