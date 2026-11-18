# package-tmp-leak-guard

decision · `01M7Q02888E9MAC9FTS0KVXS5D`

## reasoning

guardProjections covers only Ctxs from testCtx and snapshot.copy; testCache answers no TMPDIR, so snapshot builds and testCtxNoIndex Ctxs fall back to /tmp. One package folder checked at the end covers them without reading or cleaning the host's shared /tmp, which other runs write to concurrently.

## statement

internal/api's TestMain gives every test Ctx a TMPDIR of its own (a rhiza-api-tmp- folder made once, answered by testCache, swept when stale like rhiza-api-cache-) and fails the test binary when any rhiza-projection-* folder is left in it after m.Run, so a projection opened outside guardProjections (snapshot builds, testCtxNoIndex) still lands under the test run and is still checked

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [leak-plan-interview](../source/leak-plan-interview.md)

### In

- decided-by [api-tests-tmp-leak-guard](../capability/rhiza/views/fast-reads/api-tests-tmp-leak-guard.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
