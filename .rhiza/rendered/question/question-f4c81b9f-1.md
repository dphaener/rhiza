# question-f4c81b9f-1

question · `01M7PZYYT99GN2EXD6AVFB427X`

## session

f4c81b9f-379e-416c-9e48-01a1edfbc5f3

## statement

Change set 01M7PZQ4 (throwaway projections leak): the /tmp leak's cause is found (workspace.projection() falls back to a bare Ctx with no Getenv, so the throwaway index goes to the real /tmp untracked by the test guard, and the test helper ctxWorkspace never closes it; reproduced on current code: 29 folders from 4 tests). The intent also says 'look at why rhiza ship open-pr needs 2.2 GB RSS'. What I found: open-pr's pre-push check runs rhiza check in-process over a fresh checkout (ship_prepush.go prePushRhizaCheck); rhiza check alone on this repo (7,298 log entries) peaks at 832 MB RSS. Cutting that needs profiling of the check/fold path, unrelated to the leak's files. Where should the open-pr memory work go?

## Links

### In

- answers [open-pr-memory-separate](../decision/open-pr-memory-separate.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
