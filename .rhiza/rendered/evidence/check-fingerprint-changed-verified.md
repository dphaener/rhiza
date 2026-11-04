# check-fingerprint-changed-verified

evidence · `01M6JA2DD49TJ0681HF3JYS72R` · actor [claude-code-4b7e64f4](../actor/claude-code-4b7e64f4.md)

## result

PASS

## statement

Verifier at b4f08e38. go test ./internal/check -run Anchor -count=1: TestAnchorFingerprintChanged passes its four cases. A changed declaration and a changed whole file each give exactly one warn finding naming anchor pin and implements thing, with Next rhiza anchor refresh pin and Report.OK true. A change outside the declaration and an unchanged file give no anchor finding. go test ./cmd/rhiza -run TestScripts/check-fingerprint$ passes. After editing only bytes outside Keep, check prints ok with no anchor line. After editing Keep's body, check exits 0 and prints 'warn anchor pkg/p.go: anchor pin (implements thing): symbol Keep of pkg/p.go changed since it was fingerprinted' with rhiza anchor refresh pin, and check --json reports ok true. make fmt-check vet self-check pass (self-check prints ok). In the full go test ./... run on a host at load average ~10 on 4 cores, only timing tests failed (TestReadBudget in internal/api, TestBudget in internal/index, TestBudget timing_test.go, medians just over 50ms) plus a plan-workflow script timeout. plan-workflow passes when run alone. None of these is in the diff, which touches only internal/check and the new txtar. The diff reuses the resolved FingerprintFS result, warns only when a stored fingerprint exists and differs, and moves the implements wording into implementsOf, which clones the slice before sorting.

## Links

### Out

- verifies [check-fingerprint-changed](../capability/rhiza/integrity/fingerprints/check-fingerprint-changed.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
