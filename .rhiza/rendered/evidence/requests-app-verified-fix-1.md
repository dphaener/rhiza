# requests-app-verified-fix-1

evidence · `01M6DF3YN9V6KVJCX6SY8DBCW7` · actor [claude-code-f4386038](../actor/claude-code-f4386038.md)

## result

PASS

## statement

Verified at 99c050d8 (review fixer, quality f1): Submit now builds argv with RequestArgs = [--kind, kind, --, text], so dash-leading text is filed as text. Ran `go test -count=1 ./internal/tui/requestsapp/ -v`: all pass, including TestFilesBugAndFeatureFromTypedText (bug and feature filed from typed text, shown queued on a fake runner) and the new TestDashLeadingTextIsFiledAsText; the fake now rejects any other argv shape. Ran `go test -count=1 ./cmd/rhiza -run 'TestScripts/(tui-requests-app|rex-install|describe-golden)$|TestRequestsAppArgs|TestRequestFromArgs'`: tui-requests-app.txtar (--once JSON equals request ls), rex-install.txtar (rhiza_app_requests on a free key, no key moved) and describe-golden pass, and TestRequestsAppArgsKeepDashLeadingText runs RequestArgs(bug, --json output is wrong) through requestFromArgs for the real request op and gets Text=--json output is wrong, Kind=bug. Read the diff: 3 Go files only, no .rhiza change. `make check`: fmt, vet and all go tests pass; self-check reports the same 8 stale-render failures (actor files for review sessions plus index.md) at this commit as at its parent dcf804b9, coming from the run's own chore(rhiza) record commits, not from this fix.

## Links

### Out

- verifies [requests-app](../capability/rhiza/orchestration/daemon/plan-requests/requests-app.md)

## History

- `01M6CRZRVGZBRB51YFPXHKZA75` A requests app files features, bugs and milestones and shows each request's progress
