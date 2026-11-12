# nested-step-paths-shown-verified

evidence · `01M76AVNF6R18D143H6K0B5E04` · actor [claude-code-d3d1e585](../actor/claude-code-d3d1e585.md)

## result

PASS

## statement

Verifier at 845b4c92: read the leaf's diffs (1137b1cc, 845b4c92) and the session-start pattern change (40c4123c); grep found no remaining code cutting a run's step at its first '/' (remaining strings.Cut on '/' are workflow-definition targets that keep the whole remainder, or unrelated paths). Ran go test -trimpath -count=1 -v: TestWatchdogNestedStep (watchdog hold and its inbox item keep ship-change/build/build whole, detail names it), TestTabTitleNestedStep (TabTitle and RunTabTitle carry ship-change/build/build), TestNestedWorkflowStepResumesAtFullPath (failure at mid/inner/two is named whole and retry resumes there), TestStartStepNestedPaths: all PASS. go test -trimpath ./cmd/rhiza -run TestScripts/run-nested-workflow-session: ok (run status text and --json show ship-change/build/build and ship-change/review/review). Package runs: internal/run, internal/session, internal/daemon ok; internal/api's tests all PASS but its enclosing-root guard flagged ~/code/.rhiza-root/usage.live.json changed, which the live rhiza-root.service rewrote at 22:14:58 after the run, not the tests (the leaf adds only test code).

## Links

### Out

- verifies [nested-step-paths-shown](../capability/rhiza/workflows/workflow-runs/nested-step-paths-shown.md)

## History

- `01M75Z51562NX7KRJH5WJ2M3HN` Fix: a workflow that wraps another workflow (kind = "workflow") can't start an agent session more than one level down, so Rhiza's new bugfix workflow (.rhiza/workflows/bugfix.toml, #435/#439, which wraps the built-in ship-change) fails every bug at its first build. Evidence: run 01M75WR1PQ4H1DM53KK6F2KAVP (bugfix, rule ready.bug, change set 01M75S3B1YNJWFDGXCN2BEEACQ) ended at 2026-11-11T18:08Z: 'ship-change/build/build (visit 1): error, then end: leaf watchdog-holds-skip-abandoned: start its builder: step "ship-change/build/build": want lower-case words joined by single dashes, or workflow/step, at most 61 characters'. Cause: stepPattern (internal/api/session.go:74) allows one '/' only (workflow/step), and startStep (session.go:493) rejects deeper nested step paths; maxStep is 61 (session.go:80). The same applies to plan-bug (wraps plan) for its agent steps, and Nebari's bugfix wraps its own ship-change the same way. Fix: accept any depth of nested workflow step paths (each segment lower-case words joined by dashes, within a length limit large enough for real nesting, e.g. 120), and check every place that validates or parses a step path (session start, run journal, inbox, watchdog, rex tab labels, rhiza run status) handles 2+ slashes; a testscript that runs a workflow wrapping ship-change through to its builder session start. Filed as a feature on purpose: as a bug it would run through the broken bugfix workflow. After it lands, the orchestrator retries run 01M75WR1 (rhiza run resume --retry).
