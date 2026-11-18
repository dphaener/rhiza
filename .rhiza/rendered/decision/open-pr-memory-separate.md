# open-pr-memory-separate

decision · `01M7PZZNYEYG25TCDK5CG81MKZ`

## reasoning

The intent's expected outcomes are all about the leak; the open-pr clause asks only to 'look at why', which the planner has done. The OOM kills came mostly from the leak filling a RAM-backed /tmp (7.1 GB), which this change set fixes; open-pr's RSS is a separate cost in the check/fold path, touching files unrelated to the leak and needing profiling work whose size is unknown. Folding it in would grow a bug fix's scope and delay the leak fix, which is the urgent part. A discovered-from draft change set keeps the finding traceable and lets it be planned and sized on its own.

## statement

Keep change set 01M7PZQ4 scoped to the /tmp projection leak (workspace.projection() fallback without Getenv, ctxWorkspace never closing its projection, every throwaway folder removed on close, tests under t.TempDir(), the leak guard checking no folder is left). The open-pr memory work goes into its own draft change set, discovered from 01M7PZQ4, carrying the finding: open-pr's pre-push check runs rhiza check in-process over a fresh checkout (ship_prepush.go prePushRhizaCheck), and rhiza check alone on this repo (7,298 log entries) peaks at 832 MB RSS, so the fix needs profiling of the check/fold path. This change set's 'look at why' clause is met by that investigation and the follow-up draft; it builds no memory reduction itself.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-f4c81b9f-1](../question/question-f4c81b9f-1.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
