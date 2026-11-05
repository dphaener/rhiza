# text-ref-display-verified-2

evidence · `01M6KTG417F7HV1V2MQBBS058X` · actor [claude-code-f4009b95](../actor/claude-code-f4009b95.md)

## result

PASS

## statement

Verifier attempt 2 at a2992808 (fix for review finding f3). Read the diff: text-refs-display.txtar now hand-writes a pre-feature log entry holding '@word and @@{ID}', asserts show/show --json print '@@word and @@@{ID}', writes that back with rhiza set and asserts the log stores 'Uses @@word and @@{ID}' again (2 entries) and still shows @@@; hand-writes a reference to a node only an unlanded change set has and asserts show prints @@{<26-char ID>}, not @@elsewhere; asserts the retired node shows exactly @@{<gone's ID>}; asserts inbox and inbox --json print 'Is @@renamed ready?' from a blocking question; asserts ls, tree, feed and roadmap (text and --json) carry no @@{ (they print no node statements); show, why, slice, search '@renamed', change show/diff (text and --json), render and tui app show/slice/why --once print @@renamed; shown text written back keeps the same IDs and literal @@@. Ran go test -trimpath ./cmd/rhiza -run TestScripts/text-refs -count=1: ok. make fmt-check vet pass; make test passes (one earlier run failed with the output lost to tail, then a full rerun passed: EXIT 0). make self-check fails only on text-ref-check-verified holding @@{01ZZ...} (qa finding f2 on text-ref-syntax, another leaf) and on stale .rhiza/rendered from run bookkeeping, not on this leaf's diff.

## Links

### Out

- verifies [text-ref-display](../capability/rhiza/model/text-references/text-ref-display.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
