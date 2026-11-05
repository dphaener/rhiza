# modelapp-timing-alone-verified-3

evidence · `01M69JB896CXS70A892067XFX3` · actor [claude-code-2ec9fdd2](../actor/claude-code-2ec9fdd2.md)

## result

PASS

## statement

At 34037ad2 (cycle 3 fix: waiter dir and lock named by the go test parent pid, so only one run's waiters take turns), read the diff: TestGeneratedModel calls quiettest.Wait (shared with api's TestReadBudget) before budgetLock, logs 'other packages: <quiet>' as its header, keeps the 'median >= 50ms' t.Errorf, and logs quiettest.Late afterwards. Ran `go test -count=1 -v ./internal/api ./internal/tui/modelapp -run 'TestReadBudget|TestGeneratedModel'`: api logged 'alone after 3.2s (last beside modelapp.test)' and modelapp 'alone after 2.7s (last beside api.test)', median 9.7ms, both ok, so the two waiters take turns instead of waiting out 3m. Ran it again with ./internal/quiettest: all PASS. Ran `go test ./...` and fmt-check + vet: all ok. In a full -v ./... run, modelapp logged 'beside rhiza.test after 3m0s' plus the late note (cmd/rhiza's scripts outlast quietWait, the same cap TestReadBudget has) and still passed at median 9.7ms. make self-check fails only on stale .rhiza/rendered actor files from this run's own sessions, not on the leaf's code.

## Links

### Out

- verifies modelapp-timing-alone (retired)

## History

- `01M696Q0P1JNY00XQ1R312XVRJ` The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run
