# dependent-writes-tests

evidence · `01M5T2DGQSTX12JC9N7P1M38ES`

## result

make check exit 0 at commit ea4c7f7

## statement

Go tests and testscripts show a write in a change set seeing the change sets it depends on, directly or not, with the landed model and its own operations.

## tests

go test: TestStackedWrites (internal/api), TestModel/Dependencies (internal/api), TestWhatIfDependencies (internal/model), TestDependencies (internal/changeset); testscript: stacked.txtar

## Links

### Out

- verifies [dependent-writes](../capability/rhiza/change-sets/dependent-writes.md)

## History

- `01M5SWTT2A55K0REHQ6JNJQZ6W` Writes see the accepted, unlanded change sets they depend on, so a sequence of change sets can build on each other. This is the gap behind the built/intended split and milestone 2's most important item.
