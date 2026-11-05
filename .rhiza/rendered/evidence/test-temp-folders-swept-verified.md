# test-temp-folders-swept-verified

evidence · `01M6MK3XBTZ5PGYPWD47YC0YSW` · actor [claude-code-40556548](../actor/claude-code-40556548.md)

## result

PASS

## statement

Verifier at a0d02ef9: read the diff. internal/tmpsweep adds one helper Stale(root,prefix)/StaleTemp(prefix) removing only directories under root with the prefix whose mtime is over an hour old (MaxAge). TestStaleRemovesOnlyOldPrefixed makes, under t.TempDir(), a prefixed folder back-dated 2h with os.Chtimes, a fresh prefixed one and an old one with another prefix, runs Stale and asserts only the stale prefixed one is gone; go test -trimpath -count=1 -v ./internal/tmpsweep/ PASS (both tests). internal/api, internal/index and internal/slice TestMain call tmpsweep.StaleTemp with rhiza-api-budget-, rhiza-index-budget-, rhiza-slice-budget- before m.Run (so before genDir is made); quiettest takeTurn calls StaleTemp(rhiza-quiet-) before MkdirAll(waiterDir). Exercised end to end: with TMPDIR set to a scratch folder holding <prefix>old (touch -d '2 hours ago') and <prefix>new for each budget prefix plus other-old, go test -run '^$' ./internal/api ./internal/index ./internal/slice left only the three -new folders and other-old; with rhiza-quiet-99999 back-dated 2h, a fresh rhiza-quiet-88888 and an old rhiza-api-budget-old, go test ./internal/quiettest/ passed and left only rhiza-quiet-88888 and rhiza-api-budget-old. gofmt -l clean and go vet clean on the five packages. make check not run here, per the run's rule that the PR's CI runs it before merge

## Links

### Out

- verifies [test-temp-folders-swept](../capability/rhiza/integrity/self-check/test-temp-folders-swept.md)

## History

- `01M6M1M10N7AP6F0DHGANKXEA1` Fix Test temp folders leak and fill the disk: TestMain in internal/api/budget_test.go removes the generated rhiza-api-budget-* model (~146MB) only when the test binary exits normally, so every go test timeout (common under load) leaves one behind; quiettest leaves rhiza-quiet-* the same way, and slice/index budgets their own. On the VM 124 budget folders (~18GB) and 369 quiet folders piled up in TMPDIR by 2026-11-04 and the disk fell to 19GB free. Have each package remove its own stale folders (older than an hour, by prefix) when it starts, or put them under one per-run directory make test deletes, and test it
