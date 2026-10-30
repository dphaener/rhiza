# checks-and-fixes-anchored-verified

evidence · `01M66HTP49YKAF2A9M57P2DJEY` · actor [claude-code-93a7a813](../actor/claude-code-93a7a813.md)

## result

PASS

## statement

At c6c2527 I ran go test ./internal/api -run 'TestStartCheck$|TestStartFix$' -count=1 -v: both PASS. Each starts a session (a check via StartCheck, a change set fixer via StartFix) with Anchor S-first, whose record holds block:a1 in window:w9, against the fake rex, and asserts rex got split -b block:a1 and no window new call. I read the diff (one new test file) and StartCheck/StartFix in internal/api/run_checks.go: both set req.anchor = r.paneAnchor(spec.Anchor), which returns the anchored record's tab. go test ./... passes. make check fails only in self-check, on stale .rhiza/rendered files from the earlier evidence commit 2a76e54 (actor claude-code-bc1374d8), which this leaf's commit does not touch.

## Links

### Out

- verifies [checks-and-fixes-anchored](../capability/rhiza/runtime/rex-sessions/checks-and-fixes-anchored.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
