# autonomy-app-toggles-verified

evidence · `01M6NAH6PHJDFP03V7TZZ0G0C2` · actor [claude-code-460f36a0](../actor/claude-code-460f36a0.md)

## result

PASS

## statement

Verifier at 3b55185c: go test -trimpath -count=1 -v ./internal/tui/autonomyapp/... passed all 7 tests with a fake runner: TestStartsAtDaemonSettings (the four Select fields Plan/Accept/Ship/Sign-off start at the daemon.settings values owner/auto/auto/owner), TestSubmitWritesChangedOnly (changing plan and ship then submitting makes exactly 2 daemon.set calls, values written, shell closed), TestSubmitUnchangedWritesNothing (0 sets, closed), TestEscCallsNothing (after changing a field, Esc makes 0 sets and closes), TestSetErrorIsResult, TestResultIsDaemonSettings. go test -trimpath ./cmd/rhiza -run 'TestScripts/(tui-autonomy-app|tui-health|describe-golden)$' passed: tui app autonomy --once cmp-matches rhiza daemon settings --json before and after a set. Read the diff: submit skips unchanged toggles and calls daemon.set with [--, key, value], which requestFromArgs in cmd/rhiza/argv.go maps to the key/value positional args; cmd/rhiza/tui.go runs the app in its modal shell and keeps --once on Result. go vet and gofmt clean.

## Links

### Out

- verifies [autonomy-app-toggles](../capability/rhiza/autonomy-app-toggles.md)

## History

- `01M6MYPA335J1S3952RJAGZJN9` Autonomy app: the plan, accept, ship and sign-off toggles (rhiza.autonomy.*), rebuilt on the modal shell over daemon.settings and daemon.set, replacing internal/tui/autonomyapp and retiring the capabilities of the old autonomy app it replaces
