# workspace-projection-no-bare-ctx-verified

evidence · `01M7Q19SM4KBTBDET7BM7B44KZ` · actor [claude-code-7bcc26c0](../actor/claude-code-7bcc26c0.md)

## result

PASS

## statement

Verifier at 9da2b3af8: with TMPDIR set to a fresh scratch folder, go test -trimpath -count=1 -p 2 ./internal/api -run 'TestProjection|TestWorkspace|TestWritesThroughIndexMatchFold$|TestCarve$|TestLogCarve$|TestReviewCheckQuestionInNodeChangeSet$|TestOtherRoot' exits 0 (incl. new TestProjectionNoneGiven: no changeSets, projectionFailure names the missing projection, no rhiza-projection-* in its TMPDIR). /tmp/rhiza-projection-* count 616 before and 616 after the run, none of them with an index.db-wal naming the scratch TMPDIR. Whole go test ./internal/api exits 0 (174s), and afterwards no /tmp/rhiza-projection-* names the scratch folder. go vet ./internal/api exits 0. Diff: projection() with none given answers failedProjection{errNoProjection} and opens no index. ctxWorkspace adds useProjection(ctx) and registers t.Cleanup(closeProjection) after guardProjections's, so its projection closes before guardProjections checks for open ones. landNodeOn uses its Ctx's projection and closes it. testWorkspace uses a Ctx whose TMPDIR is t.TempDir(). TestOtherRootGateAnswer uses openRootWorkspace and defers closeProjection. openWriteIndex, the production openWorkspace callers (which already set a projection) and cmd/rhiza are unchanged. No test is deleted or skipped.

## Links

### Out

- verifies [workspace-projection-no-bare-ctx](../capability/rhiza/views/fast-reads/workspace-projection-no-bare-ctx.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
