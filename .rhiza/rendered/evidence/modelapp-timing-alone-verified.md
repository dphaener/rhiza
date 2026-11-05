# modelapp-timing-alone-verified

evidence · `01M69E8X1NS7DQC38ES2AQ1DMK` · actor [claude-code-0beadaf3](../actor/claude-code-0beadaf3.md)

## result

PASS

## statement

At d60260b0: the diff moves TestReadBudget's quiet wait from internal/api/budget_test.go into internal/quiettest (Wait, Quiet.String, Late; same logic, only renamed and exported, its ps-parsing test moved too) and TestGeneratedModel calls quiettest.Wait() before budgetLock, logs 'other packages: <q>', keeps 'if median >= 50*time.Millisecond { t.Errorf(...) }', and logs quiettest.Late afterwards. go test -count=1 ./... passed every package (modelapp 196s, as it waited for the others). go test -count=1 -v ./internal/tui/modelapp ./internal/index ./internal/slice ./internal/gen logged 'other packages: alone after 23.7s (last beside slice.test)' then 'tree and roadmap read: median 10.12ms' and passed; run alone it logged 'alone after 2.3s', median 12.99ms. make check exited 0.

## Links

### Out

- verifies modelapp-timing-alone (retired)

## History

- `01M696Q0P1JNY00XQ1R312XVRJ` The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run
