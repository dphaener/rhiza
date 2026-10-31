# modelapp-timing-alone

capability · `01M696Q30C6Z9DX0WGCSNDSGHQ` · parent [model-tree-app](../model-tree-app.md)

## acceptance

TestGeneratedModel in internal/tui/modelapp uses the same quiet wait as TestReadBudget (moved to a shared test helper package if needed), logs an alone after header, and still fails when the read median is 50 ms or more; go test ./... passes with the package run beside others

## statement

The model app's tree and roadmap read timing runs only once go test is running nothing else, as TestReadBudget's does, and logs how long it waited

## Links

### Out

- decided-by [modelapp-timing-waits-for-quiet](../../../decision/modelapp-timing-waits-for-quiet.md)

### In

- verifies [modelapp-timing-alone-verified](../../../evidence/modelapp-timing-alone-verified.md)
- verifies [modelapp-timing-alone-qa-1](../../../evidence/modelapp-timing-alone-qa-1.md)
- verifies [modelapp-timing-alone-verified-2](../../../evidence/modelapp-timing-alone-verified-2.md)
- verifies [modelapp-timing-alone-qa-2](../../../evidence/modelapp-timing-alone-qa-2.md)
- verifies [modelapp-timing-alone-verified-3](../../../evidence/modelapp-timing-alone-verified-3.md)
- verifies [pr-138-model-app-timing-waits](../../../evidence/pr-138-model-app-timing-waits.md)

## History

- `01M696Q0P1JNY00XQ1R312XVRJ` The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run
