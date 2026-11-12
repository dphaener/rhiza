# nested-step-paths-shown-qa-1

evidence · `01M76B17EN22JCK7AT9A3JYPPE` · actor [claude-code-de15fd75](../actor/claude-code-de15fd75.md)

## result

PASS

## statement

QA at fb85ef0e with the built binary on the same scratch runs: rhiza run status listed steps and leaves at ship-change/build/build and wrap/ship-change/build/build whole; the run journal recorded the full paths; Rex tab titles (fake rex argv) carried 'ship-change/build/build · builder only #1' and 'wrap/ship-change/build/build · builder only #1'; rhiza inbox showed the failed run at step wrap/ship-change/review/review with 'rhiza run resume --retry' next; rhiza workflow describe bugfix in this repo lists ship-change/build/build as its start.

## Links

### Out

- verifies [nested-step-paths-shown](../capability/rhiza/workflows/workflow-runs/nested-step-paths-shown.md)

## History

- `01M75Z51562NX7KRJH5WJ2M3HN` Fix: a workflow that wraps another workflow (kind = "workflow") can't start an agent session more than one level down, so Rhiza's new bugfix workflow (.rhiza/workflows/bugfix.toml, #435/#439, which wraps the built-in ship-change) fails every bug at its first build. Evidence: run 01M75WR1PQ4H1DM53KK6F2KAVP (bugfix, rule ready.bug, change set 01M75S3B1YNJWFDGXCN2BEEACQ) ended at 2026-11-11T18:08Z: 'ship-change/build/build (visit 1): error, then end: leaf watchdog-holds-skip-abandoned: start its builder: step "ship-change/build/build": want lower-case words joined by single dashes, or workflow/step, at most 61 characters'. Cause: stepPattern (internal/api/session.go:74) allows one '/' only (workflow/step), and startStep (session.go:493) rejects deeper nested step paths; maxStep is 61 (session.go:80). The same applies to plan-bug (wraps plan) for its agent steps, and Nebari's bugfix wraps its own ship-change the same way. Fix: accept any depth of nested workflow step paths (each segment lower-case words joined by dashes, within a length limit large enough for real nesting, e.g. 120), and check every place that validates or parses a step path (session start, run journal, inbox, watchdog, rex tab labels, rhiza run status) handles 2+ slashes; a testscript that runs a workflow wrapping ship-change through to its builder session start. Filed as a feature on purpose: as a bug it would run through the broken bugfix workflow. After it lands, the orchestrator retries run 01M75WR1 (rhiza run resume --retry).
