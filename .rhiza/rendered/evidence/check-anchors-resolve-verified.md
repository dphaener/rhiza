# check-anchors-resolve-verified

evidence · `01M6J616DF4R6BG0T1893ZDK73` · actor [claude-code-a7f3a5e3](../actor/claude-code-a7f3a5e3.md)

## result

PASS

## statement

Verified at 8c477d25 by reading the diff and re-running the acceptance. check.Run stays pure: internal/check imports io/fs only (no os), Input.Tree fs.FS is filled with os.DirFS(ctx.Root) in internal/api/check.go, and anchor.FingerprintFS backs the resolution. go test ./internal/check -run 'TestAnchor|TestRunOrder' -v passes: fstest.MapFS cases give one fail finding for a missing path ('path gone.go is gone') and one for a missing symbol ('symbol Gone of p.go is gone'), each naming 'anchor pin' and 'implements thing', with Report.OK false; a resolving anchor, a file-only anchor and a retired anchor give no finding; with no Tree the rule is skipped; TestRunOrder puts anchor between evidence-actor and render. go test ./cmd/rhiza -run TestScripts/check-anchors passes: after Keep is removed rhiza check exits non-zero with 'fail anchor pkg/p.go: anchor pin (implements thing): symbol Keep of pkg/p.go is gone', after rm a 'path pkg/p.go is gone' finding, and once restored it is clean. The check package doc lists anchor in report order, and Rules() matches. make fmt-check vet self-check pass. A full go test ./... failed only the timing budgets TestReadBudget (api), TestBudget (slice) and TestGeneratedModel (modelapp) at load average ~15 beside other test binaries; these read paths are untouched by the diff, and TestBudget passes when run alone.

## Links

### Out

- verifies [check-anchors-resolve](../capability/rhiza/integrity/anchors/check-anchors-resolve.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
