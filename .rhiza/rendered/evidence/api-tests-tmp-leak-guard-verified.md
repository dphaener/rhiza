# api-tests-tmp-leak-guard-verified

evidence · `01M7Q1WHWQXQKQ6F4AE1TXYSNR` · actor [claude-code-b50a0b63](../actor/claude-code-b50a0b63.md)

## result

PASS

## statement

Verified at b2a1225f9 in a detached worktree. Acceptance: TMPDIR=<fresh scratch> go test -trimpath -count=1 -p 2 ./internal/api -run 'TestProjection|TestSnapshot|TestShipCheckModelConflictClosesProjection|TestRunStatusGatedClosesProjection|TestWorktreeWorkspaceClosesReplacedProjection|TestSessionWorkspaceMissClosesProjection|TestTmpLeak' exit 0; /tmp/rhiza-projection-* 96 before and 96 after, none new; scratch TMPDIR empty afterwards (no rhiza-api-tmp- left). -v shows TestTmpLeakReportsLeftProjection (rhiza-projection-x folder named on output, tmpLeakCode non-zero; a file and rhiza-run-1 ignored), TestTmpLeakEmptyFolderPasses (0, no output, failed code kept) and TestTmpLeakTestCacheAnswersPackageTmp (TMPDIR -> testTmp, XDG_CACHE_HOME unchanged, testCtxNoIndex's Getenv TMPDIR -> testTmp) passing. Checks block re-run: full TMPDIR=<scratch> go test -trimpath -count=1 -p 2 ./internal/api exit 0 (175.6s), no new /tmp projection, no rhiza-api-tmp- left; go vet ./internal/api exit 0; rhiza ship generate and rhiza render exit 0 (render drift only the conductor's uncommitted actor record, not this leaf). Diff: only internal/api/budget_test.go, init_test.go, projection_leak_test.go plus model records; no production code, guardProjections and its per-test TMPDIR unchanged, TestMain adds tmpsweep.StaleTemp(rhiza-api-tmp-), no TODOs.

## Links

### Out

- verifies [api-tests-tmp-leak-guard](../capability/rhiza/views/fast-reads/api-tests-tmp-leak-guard.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
