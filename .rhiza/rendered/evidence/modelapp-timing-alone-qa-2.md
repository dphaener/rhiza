# modelapp-timing-alone-qa-2

evidence · `01M69GZP73YJ7S2F813P4HFB4Z` · actor [claude-code-cc6dc6bd](../actor/claude-code-cc6dc6bd.md)

## result

PASS

## statement

At 594c251d, built go build -o $TMPDIR/rhiza-qa2 ./cmd/rhiza (rhiza version ran). Cycle 1's repro, go test -count=1 -v ./internal/api ./internal/tui/modelapp ./internal/quiettest -run 'TestReadBudget|TestGeneratedModel|...', took 19.5s instead of about 200s: api logged 'other packages: alone after 2.3s', modelapp logged 'alone after 3.1s (last beside api.test quiettest.test)' with a 7.8ms median, and all passed. go test -count=1 -v ./... passed all 30 packages in 5m51s. modelapp logged 'alone after 2m27.3s (last beside rhiza.test)' with a 10.7ms median. api took the lock first and logged 'beside rhiza.test after 3m0s' plus a late note: cmd/rhiza runs for 340s, longer than quietWait, and that predates this change set. The two waiters no longer wait on each other. fmt-check and vet pass. self-check fails only on stale .rhiza/rendered from the review run's model commits. Side effect seen live: the lock is machine-wide ($TMPDIR/rhiza-quiet.lock), so another session's targeted go test run (pid 25571) held its api.test on the lock for about 5 minutes behind this full run.

## Links

### Out

- verifies modelapp-timing-alone (retired)

## History

- `01M696Q0P1JNY00XQ1R312XVRJ` The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run
