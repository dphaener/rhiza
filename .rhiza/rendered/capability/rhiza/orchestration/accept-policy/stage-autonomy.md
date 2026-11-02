# stage-autonomy

capability · `01M6CS1ZKGPVYZ3HAS2EP7Q03Y` · parent [accept-policy](../accept-policy.md)

## acceptance

Testscripts with a fake rex and agents: ship owner leaves a ready change set unstarted with skip reason ship-owner; signoff owner starts run start without --autonomous; accept owner leaves the plan's drafts proposed after an autonomous plan gate and inbox --json lists an accept item for each, while accept auto accepts them

## statement

The ship, signoff and accept settings govern the daemon and plan runs: ship owner stops the daemon starting ship runs, signoff owner starts them non-autonomous so they stop at sign-off, accept owner leaves an owner-requested plan's change sets proposed at its gate and the inbox lists each as waiting for acceptance

## Links

### Out

- decided-by [stage-autonomy-meanings](../../../../decision/stage-autonomy-meanings.md)

### In

- verifies [stage-autonomy-verified](../../../../evidence/stage-autonomy-verified.md)
- verifies [stage-autonomy-qa-1](../../../../evidence/stage-autonomy-qa-1.md)
- verifies [stage-autonomy-qa-3](../../../../evidence/stage-autonomy-qa-3.md)
- verifies [pr-187-live-daemon-settings](../../../../evidence/pr-187-live-daemon-settings.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
