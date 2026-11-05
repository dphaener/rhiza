# test-temp-folders-swept-qa-1

evidence · `01M6MK86SK8Y2R7QC6QV2AWCT5` · actor [claude-code-d4b1f57a](../actor/claude-code-d4b1f57a.md)

## result

PASS

## statement

QA at 78269c0d, hands-on with the built test binaries: in a scratch TMPDIR made, for each of rhiza-api-budget-, rhiza-index-budget-, rhiza-slice-budget- and rhiza-quiet-, a prefixed folder with content back-dated 2h (touch -d) and a fresh prefixed one, plus an old other-old folder. Built go test -c binaries for internal/api, internal/index and internal/slice and ran each with TMPDIR set and -test.run '^$' (TestMain only): all exit 0, and afterwards the three old budget folders were gone while their fresh ones and other-old stayed. Built internal/quiettest's binary and ran TestTakeTurnOneAtATime (calls takeTurn) with the same TMPDIR: PASS, rhiza-quiet-old gone, rhiza-quiet-fresh and other-old kept, and the run's own rhiza-quiet-<ppid> removed on exit. go test -trimpath -count=1 ./internal/tmpsweep/ ./internal/quiettest/ ok.

## Links

### Out

- verifies [test-temp-folders-swept](../capability/rhiza/integrity/self-check/test-temp-folders-swept.md)

## History

- `01M6M1M10N7AP6F0DHGANKXEA1` Fix Test temp folders leak and fill the disk: TestMain in internal/api/budget_test.go removes the generated rhiza-api-budget-* model (~146MB) only when the test binary exits normally, so every go test timeout (common under load) leaves one behind; quiettest leaves rhiza-quiet-* the same way, and slice/index budgets their own. On the VM 124 budget folders (~18GB) and 369 quiet folders piled up in TMPDIR by 2026-11-04 and the disk fell to 19GB free. Have each package remove its own stale folders (older than an hour, by prefix) when it starts, or put them under one per-run directory make test deletes, and test it
