# leaf-order-verified

evidence · `01M5X0Z2Q3M9X64873P5M7AM9M` · actor [claude-code-b6232fc2-verifier](../actor/claude-code-b6232fc2-verifier.md)

## result

PASS under make check at commit a345b25 on branch cs/leaf-dependencies-and, run by an independent verifier (claude-code-b6232fc2-verifier)

## statement

rhiza change show --json lists leaves c, b, a (added in that order, c depends-on b, b depends-on a) as a, b, c with b needs [a] and c needs [b]; leaves with no links come out in tree order; parents and retired capabilities are not leaves; the describe golden carries the new leaves field.

## tests

testscript leaf-order, describe-golden; go test internal/api: TestChangeShowLeaves, TestChangeShowLeavesLanded

## Links

### Out

- verifies [leaf-order](../capability/rhiza/model/leaf-dependencies/leaf-order.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
