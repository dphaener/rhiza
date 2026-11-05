# pr-138-model-app-timing-waits

evidence · `01M69M33FGCX3G8V85TG6QDQ6C` · actor [darin-haener](../actor/darin-haener.md)

## head

bd18af1830531fc596f03ea4c8c83baaf3fdd942

## pr

https://github.com/dphaener/rhiza/pull/138

## result

1 check passed: check SUCCESS

## statement

Pull request #138 (https://github.com/dphaener/rhiza/pull/138) passed CI at bd18af1

## status

verified

## Links

### Out

- verifies modelapp-timing-alone (retired)

## History

- `01M696Q0P1JNY00XQ1R312XVRJ` The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run
