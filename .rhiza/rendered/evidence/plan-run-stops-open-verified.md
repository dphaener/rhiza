# plan-run-stops-open-verified

evidence · `01M69DZJJ1QFZ2CMK5G8FQ5JE7` · actor [claude-code-21f16464](../actor/claude-code-21f16464.md)

## result

PASS

## statement

At d56666c1: go test ./cmd/rhiza -run 'TestScripts/plan-workflow-open$' -count=1 -v passes. Read the script: with the fake rex/claude (plus gh) of the plan workflow scripts, the planner parks question park-disk; the run output is not 'succeeded', gh.log does not exist (no pr create), run status --json shows stopped at gate with 'open questions: park-disk' and 'rhiza run resume', inbox --json has kind open naming park-disk with next 'rhiza run resume <run>'; after a Decision answers park-disk, run resume exits 0, prints succeeded, gh.log has exactly one 'pr create', and the open inbox item is gone. go test ./internal/run -run TestCommandStopExitStopsRun passes (exit 75 stops the run at the gate without ending it; resume reruns it and the run succeeds). Diff: run gate returns ExitCode 75 with a 'stop:' line when questions are open; the step runner captures it as the stop detail; the engine stops instead of failing; inbox gains kind open. make fmt-check, vet and test all pass; make self-check passes with RHIZA_ROOT/CHANGE/RUN/ACTOR unset (with them set it only flags this verifier session's own unrendered actor file).

## Links

### Out

- verifies [plan-run-stops-open](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-stops-open.md)

## History

- `01M625W0AJZ6AMCCWGFXQW19FN` A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.
