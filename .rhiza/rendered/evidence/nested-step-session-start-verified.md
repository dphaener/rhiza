# nested-step-session-start-verified

evidence · `01M76AE1D2VX7ZXZ7P6M4E1VQK` · actor [claude-code-331cc548](../actor/claude-code-331cc548.md)

## result

PASS

## statement

Verifier at 40c4123c: read the diff (stepPattern now repeats /segment any number of times, maxStep 120). go test -trimpath -count=1 ./internal/api -run TestStartStepNestedPaths passes: it accepts build, ship-change/build, ship-change/build/build, a/b/c/d and a 120-char three-segment path, and rejects a//b, /a, a/, a/b/, upper case, double dashes, trailing/leading dash and a 121-char path. All internal/api tests passed; the package reported FAIL only because the test run changed ~/code/.rhiza-root/usage.live.json, a known leak from running under ~/code that this change does not touch. go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/run-nested-workflow-session passes: a repo workflow wrap with a workflow step wrapping ship-change runs with fake rex/claude/gh, the fake builder starts (role=builder leaf=only), no step-pattern usage error appears, and both rhiza run status and the run journal under .git/rhiza/runs record step ship-change/build/build.

## Links

### Out

- verifies [nested-step-session-start](../capability/rhiza/workflows/workflow-runs/nested-step-session-start.md)

## History

- `01M75Z51562NX7KRJH5WJ2M3HN` Fix: a workflow that wraps another workflow (kind = "workflow") can't start an agent session more than one level down, so Rhiza's new bugfix workflow (.rhiza/workflows/bugfix.toml, #435/#439, which wraps the built-in ship-change) fails every bug at its first build. Evidence: run 01M75WR1PQ4H1DM53KK6F2KAVP (bugfix, rule ready.bug, change set 01M75S3B1YNJWFDGXCN2BEEACQ) ended at 2026-11-11T18:08Z: 'ship-change/build/build (visit 1): error, then end: leaf watchdog-holds-skip-abandoned: start its builder: step "ship-change/build/build": want lower-case words joined by single dashes, or workflow/step, at most 61 characters'. Cause: stepPattern (internal/api/session.go:74) allows one '/' only (workflow/step), and startStep (session.go:493) rejects deeper nested step paths; maxStep is 61 (session.go:80). The same applies to plan-bug (wraps plan) for its agent steps, and Nebari's bugfix wraps its own ship-change the same way. Fix: accept any depth of nested workflow step paths (each segment lower-case words joined by dashes, within a length limit large enough for real nesting, e.g. 120), and check every place that validates or parses a step path (session start, run journal, inbox, watchdog, rex tab labels, rhiza run status) handles 2+ slashes; a testscript that runs a workflow wrapping ship-change through to its builder session start. Filed as a feature on purpose: as a bug it would run through the broken bugfix workflow. After it lands, the orchestrator retries run 01M75WR1 (rhiza run resume --retry).
