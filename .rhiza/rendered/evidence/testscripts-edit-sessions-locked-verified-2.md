# testscripts-edit-sessions-locked-verified-2

evidence · `01M78JFMWT9BXSTYRYEK2TE2PW` · actor [claude-code-185c9b52](../actor/claude-code-185c9b52.md)

## result

PASS

## statement

Verifier, at a60935df9 (attempt 2). Delay injected (not committed) in scratch worktrees per the leaf's note: in internal/api/session.go, before session start's tab Update, poll st.Read(sid) until Status.Hook is set (5ms, up to 2s), and sleep 300ms inside the Update closure before r.Tab = tab. Before the change (6e5cf35bf~1), go test -trimpath ./cmd/rhiza -run 'TestScripts/inbox-stalled$' -count=1 failed 3/3 at inbox-stalled.txtar:36 (first stalled check, the unlocked hook-at edit dropped). At HEAD it passed 3 consecutive runs, and more: 25 of 29 runs passed in all, with at most 7 passes in a row. The 4 failures were at :43 with inbox items [] and no 'session started' line from the run: inbox ran while the injected sleep still held session start, before the run journal named the session. No edit was dropped. That is a separate wait gap the 300ms injection widens. inbox-stalled and run-status-stalled wait only for the hook, not for the run to record its session. The leaf does not claim that gap; a follow-up could wait for block:b1 as inbox-permission now does. With the same delay, inbox-permission, daemon-notify-hook and run-status-stalled passed -count=3, so the quality finding is fixed: inbox-permission and daemon-notify-hook now wait for block:b1 before editing. Without the delay, on the leaf checkout, go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-stalled|inbox-permission|run-status-stalled|daemon-notify-hook)$' -count=3 passed (ok 39.3s). No script in cmd/rhiza/testdata/script rewrites .git/rhiza/sessions/*.json with perl (grep). internal/session is unchanged. The edit-sessions command in cmd/rhiza/main_test.go applies each edit through session.Store.Update, under the store's lock.

## Links

### Out

- verifies [testscripts-edit-sessions-locked](../capability/rhiza/integrity/self-check/testscripts-edit-sessions-locked.md)

## History

- `01M78ETR0NZZX29FXHDTZRYPKS` Fix flaky test scripts (4/4) / UNKNOWN STEP: failed then passed on a rerun in PR #498 at head d1e350c60c5ad977683e83207632bbe5ffa12be6; reproduce it locally, repeatedly, then fix it or quarantine it
