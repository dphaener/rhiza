# anchor-fields-verified

evidence · `01M6J3FK45J8A5SJ6Q0YNGBPRW` · actor [claude-code-af895225](../actor/claude-code-af895225.md)

## result

PASS

## statement

At 041cee52 (RHIZA_* env unset): go test ./internal/anchor passes; it covers a func, a value-receiver and a pointer-receiver method (T.V, T.P, generic G.M), a type, a var and a const inside grouped declarations, each hashed as exactly the go/ast node span without its doc comment, and reports a missing path, a missing symbol and a symbol on a .md file. TestScripts/anchor-fields, field-schema and describe-golden pass; describe-golden -update leaves no diff, so no schema changed. The whole cmd/rhiza script suite passes with -timeout 60m, and make fmt-check, vet and self-check pass. By hand, rhiza apply adding an anchor with path=a.txt stored sha256:98ea6e4f..., which matches sha256sum. Read the diff: anchor.Fingerprint, anchorFingerprint wired into buildAdd and buildSet, and the knownFields anchor entry. Timing budgets were flaky under load (4 cores, load 6-27): TestReadBudget also fails at base 5e327e0e, and modelapp TestGeneratedModel missed 50ms by 2.7ms run alone. The leaf touches only the write path.

## Links

### Out

- verifies [anchor-fields](../capability/rhiza/integrity/anchors/anchor-fields.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
