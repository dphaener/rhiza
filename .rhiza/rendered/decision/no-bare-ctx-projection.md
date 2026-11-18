# no-bare-ctx-projection

decision · `01M7Q02864Z0214NHD6B6BT6W0`

## reasoning

Every production caller of openWorkspace already gives a projection first, so the bare-Ctx fallback serves only test helpers, and it is the leak: no Getenv sends the throwaway index to the real /tmp and no liveProjections hides it from guardProjections. Failing loudly turns any future miss into a test failure rather than a silent folder in /tmp.

## statement

A workspace given no projection (openWorkspace without useProjection or setProjection) opens none: its reads fail with an error saying it has no projection, instead of opening a throwaway index over a bare Ctx with no Getenv in the process's temp folder; every caller, tests included, gives the workspace its Ctx's projection and closes it

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [leak-plan-interview](../source/leak-plan-interview.md)

### In

- decided-by [workspace-projection-no-bare-ctx](../capability/rhiza/views/fast-reads/workspace-projection-no-bare-ctx.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
