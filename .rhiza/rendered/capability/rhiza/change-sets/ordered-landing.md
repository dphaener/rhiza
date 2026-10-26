# ordered-landing

capability · `01M5SZHVMY4CWQJFFW14XVXTXT` · parent [change-sets](../change-sets.md)

## statement

A change set lands only after every change set it depends on has landed.

## Links

### Out

- decided-by [land-in-dependency-order](../../../decision/land-in-dependency-order.md)

### In

- verifies [ordered-landing-tests](../../../evidence/ordered-landing-tests.md)

## History

- `01M5SWTT2A55K0REHQ6JNJQZ6W` Writes see the accepted, unlanded change sets they depend on, so a sequence of change sets can build on each other. This is the gap behind the built/intended split and milestone 2's most important item.
