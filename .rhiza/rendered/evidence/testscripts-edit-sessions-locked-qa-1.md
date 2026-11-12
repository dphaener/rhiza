# testscripts-edit-sessions-locked-qa-1

evidence · `01M78HYQPTC1G7ZE73F5RVW8Z8` · actor [claude-code-884e44e2](../actor/claude-code-884e44e2.md)

## result

PASS

## statement

QA at 0a5ffad3a. Scratch detached worktrees at the merge-base 78c8b343b (before) and HEAD (after), each with an uncommitted injection in internal/api/session.go before session start's tab Update: poll st.Read(sid) every 5ms (up to 2s) until Status.Hook is set, then time.Sleep(300ms) inside the Update closure before r.Tab = tab. `go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/inbox-stalled$'` failed 3/3 before (inbox-stalled.txtar:36: no match for "kind": "stalled", the hook-at edit dropped) and passed 3/3 after. Without the delay, at HEAD, `go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-stalled|inbox-permission|run-status-stalled|daemon-notify-hook)$' -count=3` passed; go vet ./cmd/rhiza clean, gofmt clean. No script under cmd/rhiza/testdata/script rewrites .git/rhiza/sessions/*.json with perl (the remaining perl hits are inside fake Rex JSON payloads). The only Go change is cmd/rhiza/main_test.go, so the shipped binary is unchanged; built it with go build ./cmd/rhiza and ran `rhiza check` with it in a clean detached checkout of HEAD (run env unset): ok. session-run-blocked-notice.txtar still edits a session record with sed+mv, outside this leaf's four scripts.

## Links

### Out

- verifies [testscripts-edit-sessions-locked](../capability/rhiza/integrity/self-check/testscripts-edit-sessions-locked.md)

## History

- `01M78ETR0NZZX29FXHDTZRYPKS` Fix flaky test scripts (4/4) / UNKNOWN STEP: failed then passed on a rerun in PR #498 at head d1e350c60c5ad977683e83207632bbe5ffa12be6; reproduce it locally, repeatedly, then fix it or quarantine it
