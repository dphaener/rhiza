# testscripts-edit-sessions-locked

capability · `01M78F38SGXG1XXVX4V8BKTX2K` · parent [self-check](../self-check.md)

## acceptance

With a temporary delay injected inside session start's tab Update (between its read and its write, not committed), go test -trimpath ./cmd/rhiza -run 'TestScripts/inbox-stalled$' fails before the change and passes after it on 3 consecutive runs; without the delay, go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-stalled|inbox-permission|run-status-stalled|daemon-notify-hook)$' -count=3 passes; no script under cmd/rhiza/testdata/script still rewrites .git/rhiza/sessions/*.json with an unlocked perl and mv; go test -trimpath ./internal/session/... passes if internal/session changed

## statement

cmd/rhiza's testscripts inbox-stalled, inbox-permission, run-status-stalled and daemon-notify-hook make every hand edit of a live run's session record under the session store's lock, through a testscript command registered in cmd/rhiza/main_test.go, so a concurrent locked writer (session start recording the tab, a hook report) can no longer drop the edit and flake the script

## Links

### Out

- decided-by [session-edits-under-lock](../../../../decision/session-edits-under-lock.md)

### In

- verifies [testscripts-edit-sessions-locked-verified](../../../../evidence/testscripts-edit-sessions-locked-verified.md)
- verifies [testscripts-edit-sessions-locked-qa-1](../../../../evidence/testscripts-edit-sessions-locked-qa-1.md)
- verifies [testscripts-edit-sessions-locked-verified-2](../../../../evidence/testscripts-edit-sessions-locked-verified-2.md)
- verifies [pr-503-fix-flaky-test-scripts](../../../../evidence/pr-503-fix-flaky-test-scripts.md)
- discusses [note-01m78hd4w8mbj2](../../../../note/note-01m78hd4w8mbj2.md)

## History

- `01M78ETR0NZZX29FXHDTZRYPKS` Fix flaky test scripts (4/4) / UNKNOWN STEP: failed then passed on a rerun in PR #498 at head d1e350c60c5ad977683e83207632bbe5ffa12be6; reproduce it locally, repeatedly, then fix it or quarantine it
