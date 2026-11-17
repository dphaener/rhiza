# m40-run-tmp

intent · `01M7NVQ7ZQNEKVZEP01124NGNM` · parent [milestone-40](../milestone-40.md)

## statement

When a run ends, whatever its outcome (merged, failed, abandoned), nothing its steps wrote to the temporary folder is left behind. Today the run's own {run_tmp} is already removed when it ends (internal/run/engine.go, and on abandon internal/api/run.go), but only what a workflow's [env] points there lands in it: the agents and commands of a workflow that sets no TMPDIR = "{run_tmp}" (this repo's workflows set only GOTMPDIR) write to the shared tmp, where nothing removes it. The daemon gives every process a run's steps start TMPDIR = {run_tmp} unless the workflow's [env] sets TMPDIR itself. The /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's. Measure: a test runs a workflow whose [env] names no {run_tmp} and whose step writes a file under $TMPDIR, ends the run each way, and finds nothing of it left under the tmp root.

## Links

### In

- serves change set `01M7NWFRZTGPHG4DHNN8B6NWSP` (When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md
- `01M7NW2XQ9TRY5N3A7MGRAH022` apply /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza-the-daemon-runs-any-workflow/d2542435-eaf2-4a09-ab5d-01a1ebc04a37/scratchpad/fix.json
