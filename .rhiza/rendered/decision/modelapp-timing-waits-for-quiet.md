# modelapp-timing-waits-for-quiet

decision · `01M696Q1NF8JE21ZGKFS160BVG`

## statement

The model app's read timing waits until go test runs nothing else before timing, sharing TestReadBudget's quiet wait, instead of raising its 50 ms budget; it failed run 01M68KVFSG and CI on #129 at 50.6 ms under load

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m7-keep-picking](../source/m7-keep-picking.md)

### In

- decided-by modelapp-timing-alone (retired)

## History

- `01M696Q0P1JNY00XQ1R312XVRJ` The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run
