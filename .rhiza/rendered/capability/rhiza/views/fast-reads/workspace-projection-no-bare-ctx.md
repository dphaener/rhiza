# workspace-projection-no-bare-ctx

capability · `01M7Q028AA08ZP6HAFHHVHHF15` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run 'TestProjection|TestWorkspace|TestWritesThroughIndexMatchFold$|TestCarve$|TestLogCarve$|TestReviewCheckQuestionInNodeChangeSet$|TestOtherRoot' passes. A new test in projection_test.go opens a workspace with openWorkspace and no projection under a ctx whose TMPDIR is a fresh t.TempDir(): (failure case) its changeSets() answers nothing and projectionFailure() is an error naming the missing projection, and no rhiza-projection-* folder is created in that TMPDIR or by the call at all. Run with TMPDIR set to a fresh scratch folder, the five named leaking tests TestWritesThroughIndexMatchFold, TestCarve, TestLogCarve, TestReviewCheckQuestionInNodeChangeSet (through ctxWorkspace) and TestWorkspace (through testWorkspace) leave no rhiza-projection-* folder in /tmp whose index.db-wal names that scratch folder (on current code, each run alone, they leave 25, 21, 1, 4 and 7). ctxWorkspace closes its projection when its caller is done, at the latest in a t.Cleanup registered after guardProjections's, so guardProjections still reports a projection left open.

## files

internal/api/projection.go:120-140, internal/api/projection_test.go, internal/api/change_test.go (ctxWorkspace and its openWorkspace calls), internal/api/workspace_test.go, internal/api/other_root_index_test.go

## must_not

Change openWriteIndex's throwaway folder location or its removal on close and on a failed open; change the production callers of openWorkspace (status.go, run_model_cache.go, openCtxWorkspace, openRootWorkspace) beyond what the fallback's removal needs; change cmd/rhiza; delete or skip any test

## statement

A workspace opened with no projection never opens a throwaway index over a bare Ctx: workspace.projection() with none given answers a failed projection naming the missing projection, and the test helpers that open workspaces give it a projection and close it, so no test leaves a rhiza-projection-* folder in the real /tmp. ctxWorkspace and the direct openWorkspace calls in tests that open their Ctx's own root (ctx.Root) give it their Ctx's projection (useProjection(ctx)); a test that opens a workspace on a root other than its Ctx's (other_root_index_test.go:181, TestOtherRootGateAnswer's openWorkspace(o.root, ...)) opens it with openRootWorkspace(caller, o.root, false) instead, so its projection is o.root's own index (a Ctx copy with Root and IndexPath moved to that root), and closes it with closeProjection; testWorkspace, which has no Ctx, builds one over its root, clock and generator whose Getenv answers a fresh t.TempDir() for TMPDIR, gives the workspace that Ctx's projection, and closes it in a t.Cleanup

## Links

### Out

- decided-by [no-bare-ctx-projection](../../../../decision/no-bare-ctx-projection.md)

### In

- verifies [workspace-projection-no-bare-ctx-verified](../../../../evidence/workspace-projection-no-bare-ctx-verified.md)
- verifies [workspace-projection-no-bare-ctx-qa-1](../../../../evidence/workspace-projection-no-bare-ctx-qa-1.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
