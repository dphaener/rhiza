# stall-reason-on-run

decision · `01M7QBVD6XKWV353PWZJQDKXTH`

## statement

The watchdog's stall recovery states its reason on the run: each retry and each give-up names the stalled session, its role (a spawned child's role and parent, or the step's own session) and the time it last reported, in the run's journal as a warning entry that rhiza run status shows, in the charged attempt's action, and in the hold's detail and inbox item once the budget is spent; the give-up stays the existing hold and inbox item rather than a new failed state, as failed-run retry belongs to the m40-retry change set

## wording

(assumed, would have escalated) the conductor's plan of an autonomous run, from the intent's 'failed with a reason'

## Links

### Out

- derived-from [m40-stalls-conduct-plan](../source/m40-stalls-conduct-plan.md)

### In

- decided-by [stall-reason-journaled](../capability/rhiza/orchestration/daemon/daemon-watchdog/stall-reason-journaled.md)

## History

- `01M7NVMG0EJFRJV03G025C6366` The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch
