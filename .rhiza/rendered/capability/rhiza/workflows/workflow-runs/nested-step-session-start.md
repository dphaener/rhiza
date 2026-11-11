# nested-step-session-start

capability · `01M75ZJKQBWS387RP6CCVZHM4Y` · parent [workflow-runs](../workflow-runs.md)

## acceptance

internal/api tests accept step paths of one, two and three segments up to 120 characters and reject an empty segment, a leading or trailing '/', upper case, double dashes and 121 characters; a testscript in cmd/rhiza/testdata/script runs a repo workflow whose workflow step wraps ship-change, with a fake agent, through to the builder's session start for step ship-change/build/build, and the run journal records that step

## statement

An agent step of a workflow nested two or more levels deep (a workflow step wrapping a workflow that wraps another, such as bugfix wrapping ship-change wrapping build) starts its session: session.start accepts its step path, such as ship-change/build/build

## Links

### Out

- decided-by [step-path-any-depth](../../../../decision/step-path-any-depth.md)
- decided-by [nested-step-paths-serve-m4](../../../../decision/nested-step-paths-serve-m4.md)

### In

- verifies [nested-step-session-start-verified](../../../../evidence/nested-step-session-start-verified.md)
- verifies [nested-step-session-start-qa-1](../../../../evidence/nested-step-session-start-qa-1.md)
- verifies [pr-445-fix-01m75z51562n](../../../../evidence/pr-445-fix-01m75z51562n.md)

## History

- `01M75Z51562NX7KRJH5WJ2M3HN` Fix: a workflow that wraps another workflow (kind = "workflow") can't start an agent session more than one level down, so Rhiza's new bugfix workflow (.rhiza/workflows/bugfix.toml, #435/#439, which wraps the built-in ship-change) fails every bug at its first build. Evidence: run 01M75WR1PQ4H1DM53KK6F2KAVP (bugfix, rule ready.bug, change set 01M75S3B1YNJWFDGXCN2BEEACQ) ended at 2026-11-11T18:08Z: 'ship-change/build/build (visit 1): error, then end: leaf watchdog-holds-skip-abandoned: start its builder: step "ship-change/build/build": want lower-case words joined by single dashes, or workflow/step, at most 61 characters'. Cause: stepPattern (internal/api/session.go:74) allows one '/' only (workflow/step), and startStep (session.go:493) rejects deeper nested step paths; maxStep is 61 (session.go:80). The same applies to plan-bug (wraps plan) for its agent steps, and Nebari's bugfix wraps its own ship-change the same way. Fix: accept any depth of nested workflow step paths (each segment lower-case words joined by dashes, within a length limit large enough for real nesting, e.g. 120), and check every place that validates or parses a step path (session start, run journal, inbox, watchdog, rex tab labels, rhiza run status) handles 2+ slashes; a testscript that runs a workflow wrapping ship-change through to its builder session start. Filed as a feature on purpose: as a bug it would run through the broken bugfix workflow. After it lands, the orchestrator retries run 01M75WR1 (rhiza run resume --retry).
