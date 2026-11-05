# requests-app-verified

evidence · `01M6DEJ369B8M1GV8KDX9CNR72` · actor [claude-code-c2c33fbd](../actor/claude-code-c2c33fbd.md)

## result

PASS

## statement

At 5201981a (verifier, attempt 2): the fix commit only drops the stray rendered actor page and index line flagged in f1. `env -u RHIZA_ROOT make check` exits 0 (a first run failed once in go test beside the read-budget timings, a re-run of go test ./... and make check both passed). `go test ./cmd/rhiza -run 'TestScripts/(tui-requests-app|rex-install|describe-golden)$'` passes: tui-requests-app shows `rhiza tui app requests --once` byte-equal to `rhiza request ls --json` both empty and after filing a bug and a feature (both queued); rex-install greps the rhiza_app_requests action, ["requests"] = "q", and roadmap/ready/live/inbox on their old keys. internal/tui/requestsapp unit tests pass on a fake runner: TestFilesBugAndFeatureFromTypedText files 'login fails' as kind bug and 'dark mode' as a feature and the view shows them queued; TestEnterFocusesPlanRunTab covers enter on a planning request. Launcher key q is registered by no other app, and no launcher code changed.

## Links

### Out

- verifies requests-app (retired)

## History

- `01M6CRZRVGZBRB51YFPXHKZA75` A requests app files features, bugs and milestones and shows each request's progress
