# leaf-depends-on-guide-verified

evidence · `01M5X0Z2Z4FS6B6SVJT6XNF1WD` · actor [claude-code-b6232fc2-verifier](../actor/claude-code-b6232fc2-verifier.md)

## result

PASS under make check at commit a345b25 on branch cs/leaf-dependencies-and, run by an independent verifier (claude-code-b6232fc2-verifier)

## statement

rhiza guide apply contains a recipe linking export-csv depends-on report-query and says a link that would close a loop is refused; the apply topic is 663 tokens, within its cap of 800.

## tests

testscript guide; go test internal/guide: TestCaps

## Links

### Out

- verifies [leaf-depends-on-guide](../capability/rhiza/model/leaf-dependencies/leaf-depends-on-guide.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
