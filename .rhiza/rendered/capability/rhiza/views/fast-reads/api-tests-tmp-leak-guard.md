# api-tests-tmp-leak-guard

capability · `01M7Q028CFJFVMEDZAKVXVRPNF` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run 'TestProjection|TestSnapshot|TestShipCheckModelConflictClosesProjection|TestRunStatusGatedClosesProjection|TestWorktreeWorkspaceClosesReplacedProjection|TestSessionWorkspaceMissClosesProjection|TestTmpLeak' passes, and afterwards no rhiza-api-tmp- folder of that run is left. The end-of-run check is a function over a folder (leftProjections or similar) covered by a test: (failure case) a folder holding a rhiza-projection-x subfolder is reported with its name and makes TestMain exit non-zero; an empty one is not. testCache answers the package folder for TMPDIR and still answers XDG_CACHE_HOME as today, and the snapshot builds (testCtxNoIndex) put their throwaway indexes there, not in /tmp. tmpsweep.StaleTemp sweeps rhiza-api-tmp- like rhiza-api-cache-.

## files

internal/api/budget_test.go (TestMain), internal/api/init_test.go (testCache, testCtxNoIndex), internal/api/projection_leak_test.go

## must_not

Touch production code; read, glob or clean the host's /tmp beyond the run's own rhiza-api-tmp- folder; replace guardProjections or its per-test TMPDIR; change any other test package

## statement

internal/api's test binary gives every test Ctx a TMPDIR of its own (testCache answers a rhiza-api-tmp- folder TestMain makes once and sweeps when stale) and fails, naming them, when any rhiza-projection-* folder is left there after m.Run, then removes the folder

## Links

### Out

- decided-by [package-tmp-leak-guard](../../../../decision/package-tmp-leak-guard.md)

### In

- verifies [api-tests-tmp-leak-guard-verified](../../../../evidence/api-tests-tmp-leak-guard-verified.md)
- verifies [api-tests-tmp-leak-guard-qa-1](../../../../evidence/api-tests-tmp-leak-guard-qa-1.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
