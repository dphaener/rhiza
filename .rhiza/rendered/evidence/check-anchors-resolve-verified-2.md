# check-anchors-resolve-verified-2

evidence · `01M6JD6AH24SZ3KEEZ76JVZYQH` · actor [claude-code-423ff437](../actor/claude-code-423ff437.md)

## result

PASS

## statement

Verified attempt 2 at fc05f9e7. The fix commit adds anchor (between evidence-actor and render) to CheckFinding.Rule's enum and the Findings desc in internal/api/check.go, and regenerates describe-golden.txtar to match (lines 1492, 1534, 10636, 10678), so the quality finding f1 is fixed. Re-ran the acceptance: check.Input has Tree fs.FS, filled with os.DirFS(ctx.Root) in internal/api/check.go:130, and check.Run reads only its Input. go test ./internal/check -run TestAnchor -v passes: fstest.MapFS cases give one fail finding each for a missing path and a missing symbol, naming 'anchor pin' and 'implements thing', with Report.OK false; a resolving anchor and a retired one give none. TestScripts/check-anchors and TestScripts/describe-golden pass: deleting the function gives 'fail anchor pkg/p.go: anchor pin (implements thing): symbol Keep of pkg/p.go is gone' with a non-zero exit, and deleting the file gives the path finding. The check package doc lists anchor after evidence-actor, in report order. gofmt -l is clean and go vet passes. go test passes for internal/check, internal/anchor and cmd/rhiza, and for internal/api with TestReadBudget skipped. The only failures were timing budgets under load (TestReadBudget, internal/slice timing, modelapp TestGeneratedModel), none in code this leaf touches. make self-check fails only on stale .rhiza/rendered actor files from the run's own chore commits, not on this leaf.

## Links

### Out

- verifies [check-anchors-resolve](../capability/rhiza/integrity/anchors/check-anchors-resolve.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
