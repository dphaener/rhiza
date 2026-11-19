# run-status-poll-waits-for-hold

decision · `01M7RHG4HB1E1GFVFDYCS6A394`

## reasoning

CI (PR #657, run 44423965760 attempt 1) showed 'Current step: first' under 'running' at run-status.txtar:38. Changing step first's run to 'sleep 0.5' reproduced it every time. Polling for hold alone failed at line 37 because latestRun in internal/api/run.go returns the latest run when the effective change set has none yet, which is the stopped run at hold. Requiring both lines in one read passed twice under the delay.

## statement

Fix the run-status flake in its script, not by quarantine and not in Go: the live-run poll waits until one bare `rhiza run status` read shows both the run running and its current step hold, because the run is legitimately running while step first is current, and before the live run's journal exists bare status falls back to the earlier stopped run, which is also at hold

## wording

Fix the script's poll; do not quarantine; no product change. Planner's ruling in an autonomous plan run, from the reproduction

## Links

### Out

- derived-from [run-status-flake-plan](../source/run-status-flake-plan.md)

### In

- decided-by [run-status-live-poll-waits-for-hold](../capability/run-status-live-poll-waits-for-hold.md)

## History

- `01M7RHC6VMMBXVWS3BKKFF6F8N` Fix flaky test TestScripts/run-status: failed then passed on a rerun in PR #657 at head f5b2f2b99b740b4a9a18311f074ec9e4e2df94f2; reproduce it locally, repeatedly, then fix it or quarantine it
