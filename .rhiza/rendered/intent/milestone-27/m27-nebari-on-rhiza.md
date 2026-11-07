# m27-nebari-on-rhiza

intent · `01M6J303FW4PYT5E5CE39QNHJ8` · parent [milestone-27](../milestone-27.md)

## statement

Nebari's daemon builds and ships a Nebari change set on its own with Rhiza's workflows, the ship waiting on Nebari's deploy health check, with no afk or shipit left.

## Links

### Out

- decided-by [nebari-tooling-retired](../../decision/nebari-tooling-retired.md)

### In

- serves change set `01M6J375CH04RMRS4NN372NCN2` (Deploy check on ship: ship-change waits on a repo's deploy health URL after merge, from daemon settings, so Nebari can ship through Rhiza)
- serves change set `01M6TB6N9BYGB9S6ZF3HS4EDWM` (Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests)

## History

- `01M6J2RTBPHTT1VFJYW2ZA5FRC` Milestone 26: an existing repo with built code comes into Rhiza, Nebari first, and its work then runs on Rhiza workflows
- `01M6MYK80HN9E4DG8QREZH38DW` Milestone 26: Rhiza runs are fully managed from five small single-purpose Rex apps built on rexuikit, each a small modal that opens, does one job and closes; the Nebari import moves to milestone 27
