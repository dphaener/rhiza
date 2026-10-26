# acyclic-dependencies-tests

evidence · `01M5T2DH90EFA4F1X2NFW8BPCK` · actor [claude-code-a7bf1047](../actor/claude-code-a7bf1047.md)

## result

make check exit 0 at commit ea4c7f7

## statement

Go tests and testscripts show link refusing a depends-on that closes a loop between change sets, and check and ready reporting one.

## tests

go test: TestCycles, TestWouldCycle (internal/changeset), TestLinkCycle (internal/api), TestRunCycle (internal/check), TestPlanCycle (internal/api); testscript: stacked-cycle.txtar

## Links

### Out

- verifies [acyclic-dependencies](../capability/rhiza/change-sets/acyclic-dependencies.md)

## History

- `01M5SWTT2A55K0REHQ6JNJQZ6W` Writes see the accepted, unlanded change sets they depend on, so a sequence of change sets can build on each other. This is the gap behind the built/intended split and milestone 2's most important item.
