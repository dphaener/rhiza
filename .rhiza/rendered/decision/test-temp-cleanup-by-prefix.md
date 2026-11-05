# test-temp-cleanup-by-prefix

decision · `01M6M1QJ6SSGRK54NZ1DP1KY00`

## reasoning

The change set has no intent; the earlier test-suite fix (test-suite-shared, under self-check) served milestone-25, and this is the same kind of backlog test-hygiene work, so it serves the same intent. For cleanup, the leak comes from go test timeouts killing the binary before TestMain's removal runs, so only cleanup at the next start is reliable. Sweeping by prefix and an hour's age works for bare go test, CI and make test alike, and leaves alone folders a concurrent run in another worktree is still using (make test runs many worktrees on one machine, and timing budgets finish well within an hour). A per-run directory deleted by make test would not help ad-hoc go test runs, and a run killed with make would leak it the same way. The hour threshold and prefixes come from the change set's own intent.

## statement

The change set serves milestone-25 (the backlog intent), as test-suite-shared did. Each package that makes a fixed-prefix temp folder (internal/api rhiza-api-budget-*, internal/index rhiza-index-budget-*, internal/slice rhiza-slice-budget-*, internal/quiettest rhiza-quiet-*) removes its own stale folders when it starts: those under os.TempDir() with its prefix whose modification time is older than an hour, through one small shared helper, with a test that a stale folder goes and a fresh one stays. No per-run directory that make test deletes.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-45b188d1-1](../question/question-45b188d1-1.md)
- derived-from [test-temp-cleanup-plan](../source/test-temp-cleanup-plan.md)

### In

- decided-by [test-temp-folders-swept](../capability/rhiza/integrity/self-check/test-temp-folders-swept.md)

## History

- `01M6M1M10N7AP6F0DHGANKXEA1` Fix Test temp folders leak and fill the disk: TestMain in internal/api/budget_test.go removes the generated rhiza-api-budget-* model (~146MB) only when the test binary exits normally, so every go test timeout (common under load) leaves one behind; quiettest leaves rhiza-quiet-* the same way, and slice/index budgets their own. On the VM 124 budget folders (~18GB) and 369 quiet folders piled up in TMPDIR by 2026-11-04 and the disk fell to 19GB free. Have each package remove its own stale folders (older than an hour, by prefix) when it starts, or put them under one per-run directory make test deletes, and test it
