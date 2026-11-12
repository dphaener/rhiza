# testscripts-edit-sessions-locked-verified

evidence · `01M78HHSN5MBWE2S8BXQ66BSNV` · actor [claude-code-9ae23b7d](../actor/claude-code-9ae23b7d.md)

## result

PASS

## statement

Verified at 4c8f5ef67. Scratch detached worktrees at 78c8b343b (before) and 4c8f5ef67 (after), with an uncommitted injection in internal/api/session.go just before session start's tab Update: poll st.Read(sid) every 5ms (up to 2s) until Status.Hook is set, then time.Sleep(300ms) inside the Update closure, between its locked read and write, before r.Tab = tab. `go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/inbox-stalled$'` failed 3/3 before (inbox-stalled.txtar:36, no stalled item; a debug cat showed the tab recorded but hook.at still the real 2026-10-24T12:00:06Z, so the unlocked perl+mv edit was dropped) and passed 3/3 after. Without the delay, `go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-stalled|inbox-permission|run-status-stalled|daemon-notify-hook)$' -count=3` passed. grep finds no script under cmd/rhiza/testdata/script that rewrites .git/rhiza/sessions/*.json with perl; the four scripts use the new edit-sessions command (cmd/rhiza/main_test.go), which applies each edit through session.Store.Update under the store's lock. internal/session is unchanged, so its tests were not needed. Not in scope: session-run-blocked-notice.txtar and printout-gap-names-reason.txtar still edit session records with sed+mv.

## Links

### Out

- verifies [testscripts-edit-sessions-locked](../capability/rhiza/integrity/self-check/testscripts-edit-sessions-locked.md)

## History

- `01M78ETR0NZZX29FXHDTZRYPKS` Fix flaky test scripts (4/4) / UNKNOWN STEP: failed then passed on a rerun in PR #498 at head d1e350c60c5ad977683e83207632bbe5ffa12be6; reproduce it locally, repeatedly, then fix it or quarantine it
