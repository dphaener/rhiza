# unique-finding-ids-qa-1

evidence · `01M64TR90B32HHH28R4YK7NW7G` · actor [claude-code-2d141294](../actor/claude-code-2d141294.md)

## result

PASS

## statement

QA at db28257 with a binary built from the branch (go build -o $TMPDIR/rhiza-qa ./cmd/rhiza), in scratch git repos with RHIZA_* unset. Plan review round path: a run held at a wait step, 16 concurrent plain `rhiza run findings f.json --json` calls (3 findings each) all exited 0, responses named rounds 1..16 once each, the journal held 16 findings entries each numbered f1,f2,f3, and run status listed Round 1..16 with 3 findings each. Leaf and check paths: a real `rhiza run start` of a builder/verifier + checks workflow (stub rex/claude only launch the session; the findings calls are the real binary), the verifier and check alpha each firing 10 concurrent findings calls: all 20 exited 0, the leaf's findings were f1..f10 once each, and check alpha had rounds 1..10 (cycle 1) each with f1,f2. go test ./cmd/rhiza -run TestScripts/findings-concurrent -count=3 passed.

## Links

### Out

- verifies [unique-finding-ids](../capability/unique-finding-ids.md)

## History

- `01M5ZNV40RB8VYBB0A3CPXCNXZ` rhiza run findings numbers findings under the run's lock, so concurrent calls in one session never journal two findings with the same ID (leaf findings and plan review findings alike).
