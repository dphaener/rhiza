# live-app-focus-verified

evidence · `01M698KMENB5ZM3HA6XV2ECBB9` · actor [claude-code-86b50697](../actor/claude-code-86b50697.md)

## result

PASS

## statement

At 85f19c7 ran go test ./internal/tui/liveapp -run TestEnter -v -count=1: TestEnterFocusesLiveRunTab and TestEnterOnFinishedRunCallsNothing PASS. With a fake runner, enter on live row C2 calls run.open once with [C2]; opened:true shows 'Focused Rex tab block:b1 of run R2'; tab_gone:true shows 'The tab of run R2 is gone'; enter on the finished run R1 makes no run.open call and shows 'is not live'. Read the diff: focus() in internal/tui/liveapp/liveapp.go returns before calling the runner when the row is not Live, and decodes api.RunOpenResponse; cmd/rhiza/tui.go wires the runner to the real api call. go test ./..., gofmt check and go vet pass (one earlier make test run failed only the read-timing budget under load, passing on re-run); make self-check reports only stale rendered actor files and index from the run, not this leaf.

## Links

### Out

- verifies [live-app-focus](../capability/rhiza/live-app-focus.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
