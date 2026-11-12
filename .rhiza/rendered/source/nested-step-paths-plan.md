# nested-step-paths-plan

source · `01M75ZJ0W5C0YTD9C3REHNY38R`

## quotes

```json
[]
```

## raw_sha256

29f540e7e86a3168088f26d805f11d06253a1a5e9c8e7d01154bb753075ad865

## session

c37521a1-cb5f-48f3-b341-01a1cbf4a60c

## summary

Autonomous plan run (RHIZA_AUTONOMOUS=1), no owner present. The change set's intent names the bug: session.start's stepPattern (internal/api/session.go) accepts one '/' only and maxStep is 61, so a workflow wrapping ship-change (bugfix) fails at its builder with step "ship-change/build/build". The planner read session.go, internal/workflow/include.go (flattened inner steps are named <workflow>/<step>, so each wrapping level adds a segment) and the precedent change set 01M75RDB, whose bug workflows serve m4-orchestrated-by-rhiza. Planned two leaves under workflow-runs: session start accepts step paths of any depth (segments dashed lower-case words, at most 120 characters), proven by a testscript running a repo workflow wrapping ship-change to its builder session; and run status, inbox, watchdog, journal and Rex tab labels handle 2+ slash step paths. Choices taken as assumptions: serve m4-orchestrated-by-rhiza; 120-character limit as the intent suggests.

## title

Nested workflow step paths fix, planned autonomously

## Links

### In

- derived-from [nested-step-paths-serve-m4](../decision/nested-step-paths-serve-m4.md)
- derived-from [step-path-any-depth](../decision/step-path-any-depth.md)

## History

- `01M75Z51562NX7KRJH5WJ2M3HN` Fix: a workflow that wraps another workflow (kind = "workflow") can't start an agent session more than one level down, so Rhiza's new bugfix workflow (.rhiza/workflows/bugfix.toml, #435/#439, which wraps the built-in ship-change) fails every bug at its first build. Evidence: run 01M75WR1PQ4H1DM53KK6F2KAVP (bugfix, rule ready.bug, change set 01M75S3B1YNJWFDGXCN2BEEACQ) ended at 2026-11-11T18:08Z: 'ship-change/build/build (visit 1): error, then end: leaf watchdog-holds-skip-abandoned: start its builder: step "ship-change/build/build": want lower-case words joined by single dashes, or workflow/step, at most 61 characters'. Cause: stepPattern (internal/api/session.go:74) allows one '/' only (workflow/step), and startStep (session.go:493) rejects deeper nested step paths; maxStep is 61 (session.go:80). The same applies to plan-bug (wraps plan) for its agent steps, and Nebari's bugfix wraps its own ship-change the same way. Fix: accept any depth of nested workflow step paths (each segment lower-case words joined by dashes, within a length limit large enough for real nesting, e.g. 120), and check every place that validates or parses a step path (session start, run journal, inbox, watchdog, rex tab labels, rhiza run status) handles 2+ slashes; a testscript that runs a workflow wrapping ship-change through to its builder session start. Filed as a feature on purpose: as a bug it would run through the broken bugfix workflow. After it lands, the orchestrator retries run 01M75WR1 (rhiza run resume --retry).
