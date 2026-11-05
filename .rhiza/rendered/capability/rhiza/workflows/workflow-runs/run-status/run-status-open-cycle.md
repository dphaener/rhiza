# run-status-open-cycle

capability · `01M6ATBX8NH5KPKEAY7MMZ7JBC` · parent [run-status](../run-status.md)

## acceptance

run-status.txtar with a journal where a check's cycle 1 failed and cycle 2 is open prints 'open in cycle 2 (cycle 1 failed)' and --json has an empty outcome for it; a run summary test heads it open

## statement

A review check whose latest cycle is open reads as open, not as its previous cycle's outcome: run status text says 'open in cycle N (cycle N-1 failed)', its JSON outcome is empty and run summary heads it open

## Links

### Out

- decided-by [open-cycle-reads-open](../../../../../decision/open-cycle-reads-open.md)

### In

- verifies [run-status-open-cycle-verified](../../../../../evidence/run-status-open-cycle-verified.md)
- verifies [run-status-open-cycle-qa-1](../../../../../evidence/run-status-open-cycle-qa-1.md)
- verifies [pr-163-open-review-cycle-label](../../../../../evidence/pr-163-open-review-cycle-label.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6ATASW7ZS2DPMJM2Q12GHM3` A review check whose latest cycle is open reads as open in run status, the live app and run summary
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
