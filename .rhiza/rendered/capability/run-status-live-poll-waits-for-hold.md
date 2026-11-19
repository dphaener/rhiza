# run-status-live-poll-waits-for-hold

capability · `01M7RHG4P5E62HN684GFW7NYHQ`

## acceptance

Every go test below runs inside systemd-run --user --scope --quiet -p MemoryMax=6G -p TasksMax=300 with -p 2 and RHIZA_*/SHIP_*/REX_* unset, one at a time. (1) go test -trimpath -p 2 -count=3 ./cmd/rhiza -run 'TestScripts/^run-status$' passes. (2) (failure case) With a temporary scratch edit, never committed: step first's run = "true" in the script's gate.toml changed to run = "sleep 0.5", the script at main abf1d4cc3 fails at run-status.txtar:38 with 'Current step: first' in stdout (as in PR #657's CI), and with the fix it passes on 3 runs; the edit reverted and the commit's Checks: block records both outcomes. (3) (failure case) Under the same scratch delay, a poll that greps only '^- Current step: hold$' fails at the 'running' assertion (bare status reads the earlier stopped run), showing both conditions are needed; reverted. (4) Every other line of the script is unchanged (git diff touches only the poll line), and go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/run-status' (the -interventions, -messages, -stalled and -tokens scripts too) passes.

## files

cmd/rhiza/testdata/script/run-status.txtar

## must_not

Change any Go code (internal/ or cmd/rhiza/*.go), including latestRun's fallback or run state reporting; quarantine, skip or add the script to any skip list; loosen or remove any assertion; raise the poll's 400-try bound or 0.025s interval; edit other scripts; use busy loops, CPU stress or parallel stress runs to reproduce

## statement

In cmd/rhiza/testdata/script/run-status.txtar the live-run poll (the exec sh -c loop at line 35) captures each bare `rhiza run status` output once and exits 0 only when that same output matches both '^Run .* running$' and '^- Current step: hold$', so the assertions at lines 36-39 run only after the live run reached its wait step

## Links

### Out

- decided-by [run-status-poll-waits-for-hold](../decision/run-status-poll-waits-for-hold.md)

### In

- verifies [run-status-live-poll-waits-for-hold-verified](../evidence/run-status-live-poll-waits-for-hold-verified.md)
- verifies [run-status-live-poll-waits-for-hold-qa-1](../evidence/run-status-live-poll-waits-for-hold-qa-1.md)

## History

- `01M7RHC6VMMBXVWS3BKKFF6F8N` Fix flaky test TestScripts/run-status: failed then passed on a rerun in PR #657 at head f5b2f2b99b740b4a9a18311f074ec9e4e2df94f2; reproduce it locally, repeatedly, then fix it or quarantine it
