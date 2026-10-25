# dependent-writes

capability · `01M5SZHVFR2SPC0TXN7KBE9EAK` · parent [change-sets](../change-sets.md)

## statement

A write in a change set sees the change sets it depends on, directly or not, as well as the landed model and its own operations.

## Links

### Out

- decided-by [writes-see-dependencies](../../../decision/writes-see-dependencies.md)

### In

- verifies [dependent-writes-tests](../../../evidence/dependent-writes-tests.md)

## History

- `01M5SWTT2A55K0REHQ6JNJQZ6W` Writes see the accepted, unlanded change sets they depend on, so a sequence of change sets can build on each other. This is the gap behind the built/intended split and milestone 2's most important item.
