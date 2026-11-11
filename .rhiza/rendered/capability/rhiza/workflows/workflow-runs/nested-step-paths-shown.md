# nested-step-paths-shown

capability · `01M75ZJM0M4319FSK1RE1XPXGE` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Tests of the packages touched (internal/api, internal/daemon and any other place that splits a step on '/') feed a step path ship-change/build/build and see it kept whole: run status and the inbox show it, resume continues at it, the watchdog names it, and the session's tab title carries it; go test -trimpath on those packages and the nested testscript pass

## statement

Every place that reads, parses or shows a run's step handles step paths with two or more '/': rhiza run status, the run journal and resume, the inbox, the daemon watchdog, session messages and Rex tab titles name the full nested step and never cut it at the first '/'

## Links

### Out

- decided-by [step-path-any-depth](../../../../decision/step-path-any-depth.md)
- decided-by [nested-step-paths-serve-m4](../../../../decision/nested-step-paths-serve-m4.md)

### In

- verifies [nested-step-paths-shown-verified](../../../../evidence/nested-step-paths-shown-verified.md)
- verifies [nested-step-paths-shown-qa-1](../../../../evidence/nested-step-paths-shown-qa-1.md)
- verifies [pr-445-fix-01m75z51562n](../../../../evidence/pr-445-fix-01m75z51562n.md)

## History

- `01M75Z51562NX7KRJH5WJ2M3HN` Fix: a workflow that wraps another workflow (kind = "workflow") can't start an agent session more than one level down, so Rhiza's new bugfix workflow (.rhiza/workflows/bugfix.toml, #435/#439, which wraps the built-in ship-change) fails every bug at its first build. Evidence: run 01M75WR1PQ4H1DM53KK6F2KAVP (bugfix, rule ready.bug, change set 01M75S3B1YNJWFDGXCN2BEEACQ) ended at 2026-11-11T18:08Z: 'ship-change/build/build (visit 1): error, then end: leaf watchdog-holds-skip-abandoned: start its builder: step "ship-change/build/build": want lower-case words joined by single dashes, or workflow/step, at most 61 characters'. Cause: stepPattern (internal/api/session.go:74) allows one '/' only (workflow/step), and startStep (session.go:493) rejects deeper nested step paths; maxStep is 61 (session.go:80). The same applies to plan-bug (wraps plan) for its agent steps, and Nebari's bugfix wraps its own ship-change the same way. Fix: accept any depth of nested workflow step paths (each segment lower-case words joined by dashes, within a length limit large enough for real nesting, e.g. 120), and check every place that validates or parses a step path (session start, run journal, inbox, watchdog, rex tab labels, rhiza run status) handles 2+ slashes; a testscript that runs a workflow wrapping ship-change through to its builder session start. Filed as a feature on purpose: as a bug it would run through the broken bugfix workflow. After it lands, the orchestrator retries run 01M75WR1 (rhiza run resume --retry).
