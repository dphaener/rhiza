# m39-one-read-path

intent · `01M7GFT04Y5R6MNQED0VHBWX11` · parent [milestone-39](../milestone-39.md)

## statement

The interim model is deleted and there is one read path: internal/api has no newModel or newModelSets, no branch on whether IndexPath is set, and no readLog hook on Ctx. Measure: grep -nE 'func newModel|newModelSets|IndexPath (!=|==) ""|readLog' internal/api/*.go prints nothing.

## Links

### Out

- decided-by [projection-targets](../../decision/projection-targets.md)

### In

- serves change set `01M7GFTKCAK0GY3PPDRTPVYF56` (The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).)
- serves change set `01M7PZQ416EJR49M7DRVEBJ8WH` (Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
