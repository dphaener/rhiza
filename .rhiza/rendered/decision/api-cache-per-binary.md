# api-cache-per-binary

decision · `01M7EWKKW3H9XAN03K0XQ2ERQB`

## reasoning

Go panics when a test that called t.Setenv calls t.Parallel, so with the per-fixture Setenv the leaf's largest group of tests could never run in parallel; one per-binary folder keeps the host's ~/.cache untouched as the fixture's comment requires, and keeping run-folder tests serial avoids two tests with the same run ID sharing and removing one folder without changing non-test code.

## statement

internal/api's TestMain points XDG_CACHE_HOME at one temporary folder for the whole test binary (removed when it exits), and newSessionFixture stops calling t.Setenv, so the 65 session-fixture tests may call t.Parallel; a test that starts a run step (which creates os.UserCacheDir()/rhiza/run-tmp/<run id>), sweeps or inspects run temp folders, or sets its own XDG_CACHE_HOME stays serial, because testCtx's fixed clock and seed give every test the same run IDs

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [api-parallel-plan-review](../source/api-parallel-plan-review.md)

### In

- decided-by [api-tests-parallel-session](../capability/rhiza/integrity/self-check/api-tests-parallel-session.md)
- decided-by [api-tests-parallel-runtime](../capability/rhiza/integrity/self-check/api-tests-parallel-runtime.md)
- decided-by [api-tests-parallel-rest](../capability/rhiza/integrity/self-check/api-tests-parallel-rest.md)
- decided-by [api-tests-parallel-model](../capability/rhiza/integrity/self-check/api-tests-parallel-model.md)
- decided-by [api-tests-parallel-finish](../capability/rhiza/integrity/self-check/api-tests-parallel-finish.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
