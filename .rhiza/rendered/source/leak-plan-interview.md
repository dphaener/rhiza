# leak-plan-interview

source · `01M7Q01C8EJ93NW8ZEJ7CF9VWX`

## quotes

```json
[]
```

## raw_sha256

7754b4f193495ecb577c8695f44a763d532c68f7f89c69fd56a89cbb00aee101

## session

f4c81b9f-379e-416c-9e48-01a1edfbc5f3

## summary

Autonomous plan run, no owner. Research reproduced the leak on current code: running TestWritesThroughIndexMatchFold, TestAddKeysAndParent, TestCarveKeepsIDs and TestReviewCheckQuestionInNodeChangeSet with TMPDIR set to a scratch folder left 29 rhiza-projection-* folders in the real /tmp, each with an un-checkpointed index.db-wal (never closed). Cause: workspace.projection() in internal/api/projection.go, when a workspace was opened with openWorkspace and given no projection, opens one over a bare &Ctx{Root, Now, IDs}: no Getenv, so userdir.Temp gives /tmp, and no liveProjections, so guardProjections neither sees it open nor finds its folder in the test's TMPDIR. The test helper ctxWorkspace (change_test.go) and six direct openWorkspace calls in tests open such workspaces and never close them; 23 test files use ctxWorkspace. Every production caller of openWorkspace (status.go, run_model_cache.go readModel, openCtxWorkspace, openRootWorkspace) gives a projection first, so the fallback serves only tests. openWriteIndex already removes a throwaway folder on close and on a failed open. Snapshot builds and testCtxNoIndex use testCache as Getenv, which answers no TMPDIR, so their throwaway indexes also go to /tmp (closed, but outside any guard). The leak is from change set 01M7GFTKCA (one read path, #606), which served m39-one-read-path. On open-pr memory: its pre-push check runs rhiza check in-process over a fresh checkout (ship_prepush.go prePushRhizaCheck), and rhiza check alone on this repo (7,298 log entries) peaked at 832 MB RSS. The planner asked whether that memory work belongs here; the arbiter ruled (Decision open-pr-memory-separate) to keep this change set on the leak and put the memory work in its own discovered-from draft change set.

## title

Plan of the throwaway projections leak fix (autonomous)

## Links

### In

- derived-from [no-bare-ctx-projection](../decision/no-bare-ctx-projection.md)
- derived-from [package-tmp-leak-guard](../decision/package-tmp-leak-guard.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
