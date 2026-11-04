# m26-nebari-on-rhiza

intent · `01M6J303FW4PYT5E5CE39QNHJ8` · parent [milestone-26](../milestone-26.md)

## statement

Nebari's daemon builds and ships a Nebari change set on its own with Rhiza's workflows, the ship waiting on Nebari's deploy health check, with no afk or shipit left.

## Links

### Out

- decided-by [nebari-tooling-retired](../../decision/nebari-tooling-retired.md)

### In

- serves change set `01M6J375CH04RMRS4NN372NCN2` (Deploy check on ship: ship-change waits on a repo's deploy health URL after merge, from daemon settings, so Nebari can ship through Rhiza)

## History

- `01M6J2RTBPHTT1VFJYW2ZA5FRC` Milestone 26: an existing repo with built code comes into Rhiza, Nebari first, and its work then runs on Rhiza workflows
