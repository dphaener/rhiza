# run-status-flake-plan

source · `01M7RHFPDRW6QBFY77DPNW4E1N`

## quotes

```json
[]
```

## raw_sha256

8ed8d069d96de3eac8864e3525f58d4e739594ea61bbe4d7c961f4228629d608

## session

074f1103-3564-4552-976e-01a1f1165f3b

## summary

Autonomous plan run, no owner. PR #657's CI (run 44423965760, attempt 1) failed TestScripts/run-status at run-status.txtar:38: the poll at line 35 exits once bare `rhiza run status` reads 'running', but the live run's current step was still 'first', so '^- Current step: hold$' had no match. Reproduced deterministically by changing step first's run to 'sleep 0.5' in scratch. Polling for 'Current step: hold' alone is not enough: until the live run's journal exists, bare status falls back (latestRun in internal/api/run.go) to the earlier stopped run, also at hold, so line 37's 'running' then fails (reproduced). Requiring both 'running' and 'Current step: hold' in one status read passed twice under the same delay. Ruling: fix the script's poll, no quarantine, no Go change.

## title

run-status flake, autonomous plan

## Links

### In

- derived-from [run-status-poll-waits-for-hold](../decision/run-status-poll-waits-for-hold.md)

## History

- `01M7RHC6VMMBXVWS3BKKFF6F8N` Fix flaky test TestScripts/run-status: failed then passed on a rerun in PR #657 at head f5b2f2b99b740b4a9a18311f074ec9e4e2df94f2; reproduce it locally, repeatedly, then fix it or quarantine it
