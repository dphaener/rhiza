# permission-block-clears-qa-1

evidence · `01M701V7FYJ9W10Z9AXXMH30DP` · actor [claude-code-1b20f401](../actor/claude-code-1b20f401.md)

## result

PASS

## statement

QA cycle 1, built binary: after a hook permission block, rhiza session run with a fake rex events stream reporting block:b1 idle 12s later made rhiza session show give 'Status: idle ... (from rex)' with the hook block kept underneath, and rhiza session message then delivered at once (Waited 0s): exactly one rex send of the pointer and one send-key enter. Scoped internal/session tests and session-message-permission.txtar pass.

## Links

### Out

- verifies [permission-block-clears](../capability/rhiza/runtime/agent-status/permission-block-clears.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
