# step-path-any-depth

decision · `01M75ZJKF4ZKAZ5MK4RZKRDCMD`

## statement

A session's step is a path of one or more segments joined by '/', one per wrapping workflow level, each segment lower-case words joined by single dashes, at most 120 characters in all; no depth limit beyond the length

## wording

(assumed by the planner in an autonomous run, no owner to ask)

## Links

### Out

- derived-from [nested-step-paths-plan](../source/nested-step-paths-plan.md)

### In

- decided-by [nested-step-session-start](../capability/rhiza/workflows/workflow-runs/nested-step-session-start.md)
- decided-by [nested-step-paths-shown](../capability/rhiza/workflows/workflow-runs/nested-step-paths-shown.md)

## History

- `01M75Z51562NX7KRJH5WJ2M3HN` Fix: a workflow that wraps another workflow (kind = "workflow") can't start an agent session more than one level down, so Rhiza's new bugfix workflow (.rhiza/workflows/bugfix.toml, #435/#439, which wraps the built-in ship-change) fails every bug at its first build. Evidence: run 01M75WR1PQ4H1DM53KK6F2KAVP (bugfix, rule ready.bug, change set 01M75S3B1YNJWFDGXCN2BEEACQ) ended at 2026-11-11T18:08Z: 'ship-change/build/build (visit 1): error, then end: leaf watchdog-holds-skip-abandoned: start its builder: step "ship-change/build/build": want lower-case words joined by single dashes, or workflow/step, at most 61 characters'. Cause: stepPattern (internal/api/session.go:74) allows one '/' only (workflow/step), and startStep (session.go:493) rejects deeper nested step paths; maxStep is 61 (session.go:80). The same applies to plan-bug (wraps plan) for its agent steps, and Nebari's bugfix wraps its own ship-change the same way. Fix: accept any depth of nested workflow step paths (each segment lower-case words joined by dashes, within a length limit large enough for real nesting, e.g. 120), and check every place that validates or parses a step path (session start, run journal, inbox, watchdog, rex tab labels, rhiza run status) handles 2+ slashes; a testscript that runs a workflow wrapping ship-change through to its builder session start. Filed as a feature on purpose: as a bug it would run through the broken bugfix workflow. After it lands, the orchestrator retries run 01M75WR1 (rhiza run resume --retry).
