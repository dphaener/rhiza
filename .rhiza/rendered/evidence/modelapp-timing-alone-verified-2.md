# modelapp-timing-alone-verified-2

evidence · `01M69GD2ATSENMJG3B6XCA3K5J` · actor [claude-code-da1ab0fe](../actor/claude-code-da1ab0fe.md)

## result

PASS

## statement

Verifier, attempt 2, at da7a266c. Read the diff: TestGeneratedModel (internal/tui/modelapp) and TestReadBudget both call the shared quiettest.Wait(t), which now marks the process as a waiter, takes a shared flock, and leaves marked waiters out of what it counts (TestTestProcesses covers the skip set). The test logs 'other packages: <Quiet>' and still errors when the median is 50ms or more (modelapp_test.go:177). Ran the review finding's repro `go test -count=1 -v ./internal/api ./internal/tui/modelapp ./internal/quiettest -run 'TestReadBudget|TestGeneratedModel|TestTestProcesses'`: api logged 'alone after 2.7s (last beside modelapp.test)', modelapp 'alone after 2.6s (last beside api.test)' with a median of 8.4ms, and the run took 20s where it had taken about 200s. Ran `go test -count=1 -v ./...`: 30 packages ok, no FAIL; modelapp logged 'alone after 2m12.4s (last beside rhiza.test)' with a median of 8.4ms. make fmt-check and vet pass. make self-check fails only on stale .rhiza/rendered files, which come from the review run's model commit 0a3bfc2f, not from this Go-only commit.

## Links

### Out

- verifies modelapp-timing-alone (retired)

## History

- `01M696Q0P1JNY00XQ1R312XVRJ` The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run
