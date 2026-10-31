# modelapp-timing-alone-qa-1

evidence · `01M69F3Z9Q0YTDCFXSY10T62TD` · actor [claude-code-ef8faf9d](../actor/claude-code-ef8faf9d.md)

## result

FAIL

## statement

At 4df26617, built go build ./cmd/rhiza (rhiza version ran), then go test -count=1 -v ./internal/... ./cmd/...: TestGeneratedModel logged 'other packages: beside api.test rhiza.test after 3m0s' and 'note: api.test rhiza.test ran beside tree and roadmap as-is', and TestReadBudget logged 'beside modelapp.test rhiza.test after 3m0s'. Both passed, at a 10.3 ms median. go test of modelapp, api, quiettest, index, slice and gen together did the same: each waited the full 3 minutes for the other and then timed beside it. The two quiet waiters see each other as running, so neither runs alone in the normal run and each adds 3 minutes. quiettest unit tests pass.

## Links

### Out

- verifies [modelapp-timing-alone](../capability/rhiza/model-tree-app/modelapp-timing-alone.md)

## History

- `01M696Q0P1JNY00XQ1R312XVRJ` The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run
