# dropped-dependencies-tests

evidence · `01M5T2DH30XHKFKE4CTYGR5NQ9` · actor [claude-code-a7bf1047](../actor/claude-code-a7bf1047.md)

## result

make check exit 0 at commit ea4c7f7

## statement

Go tests and testscripts show check reporting operations that rest on a withdrawn or rejected dependency, and ready giving it as the reason.

## tests

go test: TestRestsOn, TestRestsOnSelf (internal/changeset), TestRunRestsOn (internal/check), TestPlanRestsOnDropped, TestReadyRestsOnEndToEnd (internal/api); testscript: stacked-dropped.txtar, stacked-revised-drop.txtar

## Links

### Out

- verifies [dropped-dependencies](../capability/rhiza/change-sets/dropped-dependencies.md)

## History

- `01M5SWTT2A55K0REHQ6JNJQZ6W` Writes see the accepted, unlanded change sets they depend on, so a sequence of change sets can build on each other. This is the gap behind the built/intended split and milestone 2's most important item.
