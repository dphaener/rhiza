# test-temp-folders-swept

capability · `01M6M1S9S4QNBMCMDWPDM62T6M` · parent [self-check](../self-check.md)

## acceptance

A unit test of the shared helper creates, under a temp root, a folder with the prefix back-dated over an hour (os.Chtimes), a fresh one with the prefix and an old one with another prefix, runs the sweep and asserts only the stale prefixed folder is gone; each of internal/api, internal/index and internal/slice TestMain and internal/quiettest's turn-taking call it for their own prefix before making a folder; make check passes

## statement

A test run that is killed before its cleanup (a go test timeout) leaves none of the four fixed-prefix temp folders behind for good: each package that makes one (internal/api rhiza-api-budget-*, internal/index rhiza-index-budget-*, internal/slice rhiza-slice-budget-*, internal/quiettest rhiza-quiet-*) removes, when it starts, the folders under os.TempDir() with its prefix whose modification time is older than an hour, through one shared helper. Folders made by t.TempDir() (such as internal/tui/modelapp TestGeneratedModel's) are outside this leaf

## Links

### Out

- decided-by [test-temp-cleanup-by-prefix](../../../../decision/test-temp-cleanup-by-prefix.md)

### In

- verifies [test-temp-folders-swept-verified](../../../../evidence/test-temp-folders-swept-verified.md)
- verifies [test-temp-folders-swept-qa-1](../../../../evidence/test-temp-folders-swept-qa-1.md)
- verifies [pr-271-fix-test-temp-folders-leak-and](../../../../evidence/pr-271-fix-test-temp-folders-leak-and.md)

## History

- `01M6M1M10N7AP6F0DHGANKXEA1` Fix Test temp folders leak and fill the disk: TestMain in internal/api/budget_test.go removes the generated rhiza-api-budget-* model (~146MB) only when the test binary exits normally, so every go test timeout (common under load) leaves one behind; quiettest leaves rhiza-quiet-* the same way, and slice/index budgets their own. On the VM 124 budget folders (~18GB) and 369 quiet folders piled up in TMPDIR by 2026-11-04 and the disk fell to 19GB free. Have each package remove its own stale folders (older than an hour, by prefix) when it starts, or put them under one per-run directory make test deletes, and test it
