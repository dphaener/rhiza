# run-status-open-cycle

capability · `01M6ATBX8NH5KPKEAY7MMZ7JBC` · parent [run-status](../run-status.md)

## acceptance

run-status.txtar with a journal where a check's cycle 1 failed and cycle 2 is open prints 'open in cycle 2 (cycle 1 failed)' and --json has an empty outcome for it; a Go test of the live app's checks summary counts it running; a run summary test heads it open

## statement

A review check whose latest cycle is open reads as open, not as its previous cycle's outcome: run status text says 'open in cycle N (cycle N-1 failed)', its JSON outcome is empty, the live app counts it running and run summary heads it open

## Links

### Out

- decided-by [open-cycle-reads-open](../../../../../decision/open-cycle-reads-open.md)

### In

- verifies [run-status-open-cycle-verified](../../../../../evidence/run-status-open-cycle-verified.md)
- verifies [run-status-open-cycle-qa-1](../../../../../evidence/run-status-open-cycle-qa-1.md)
- verifies [pr-163-open-review-cycle-label](../../../../../evidence/pr-163-open-review-cycle-label.md)

## History

- `01M6ATASW7ZS2DPMJM2Q12GHM3` A review check whose latest cycle is open reads as open in run status, the live app and run summary
