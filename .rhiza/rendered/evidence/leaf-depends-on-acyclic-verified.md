# leaf-depends-on-acyclic-verified

evidence · `01M5X0Z2F7S7V6N7SYS3J1WYX7` · actor [claude-code-b6232fc2-verifier](../actor/claude-code-b6232fc2-verifier.md)

## result

PASS under make check at commit a345b25 on branch cs/leaf-dependencies-and, run by an independent verifier (claude-code-b6232fc2-verifier)

## statement

rhiza link and rhiza apply refuse a depends-on that would close a loop among capabilities a, b, c (exit 3, the loop named, the log entry count unchanged), and rhiza check folds a loop made by two concurrent change sets into one fail finding on rule depends-on naming a, b and c.

## tests

testscript leaf-cycle; go test internal/api: TestLinkCapabilityLoop (Direct, Two, ApplyBatch, Unlinked, NoLoop); internal/check: TestRunCapabilityLoop (Concurrent, Retired, Chain), TestLoops

## Links

### Out

- verifies [leaf-depends-on-acyclic](../capability/rhiza/model/leaf-dependencies/leaf-depends-on-acyclic.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
