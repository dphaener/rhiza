# autonomy-app-toggles-verified-2

evidence · `01M6NBJ3GTWA4W7FV9G211H3SE` · actor [claude-code-3e0de970](../actor/claude-code-3e0de970.md)

## result

PASS

## statement

Verifier, attempt 2, at aa6c9ded: go test -trimpath -count=1 -v ./internal/tui/autonomyapp/... passes all 8 tests with a fake runner: TestStartsAtDaemonSettings (the four Select fields Plan/Accept/Ship/Sign-off start at the fake's daemon.settings values), TestSubmitWritesChangedOnly (two changed toggles give exactly two daemon.set calls, shell closes, nil error), TestSubmitUnchangedWritesNothing (0 sets, closes), TestEscCallsNothing (0 sets, closes), TestSetErrorIsResult, TestResultIsDaemonSettings, TestRegistered, and the new TestLoadErrorShownAndNotWritten (a failing daemon.settings shows the error in the rendered form, submit writes nothing and closes with that error, fixing review finding f1). go test ./cmd/rhiza -run TestScripts/(tui-autonomy-app|describe-golden)$ passes: tui app autonomy --once cmp-matches rhiza daemon settings --json before and after a daemon set. Read the diff (autonomyapp.go, tests, cmd/rhiza/tui.go): one rexuikit ui.Form of four auto/owner Selects over rhiza.autonomy.{plan,accept,ship,signoff}, submit diffs against the loaded values and writes via daemon.set -- key value, then modalshell.Done. go vet and gofmt clean.

## Links

### Out

- verifies [autonomy-app-toggles](../capability/rhiza/autonomy-app-toggles.md)

## History

- `01M6MYPA335J1S3952RJAGZJN9` Autonomy app: the plan, accept, ship and sign-off toggles (rhiza.autonomy.*), rebuilt on the modal shell over daemon.settings and daemon.set, replacing internal/tui/autonomyapp and retiring the capabilities of the old autonomy app it replaces
