# ship-signoff

capability · `01M5Z8XRQTCA20KGMAN4KKZJ5D` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

Testscripts with fake agents and a fake gh find a ship-change run waiting at sign-off after green CI with run status naming the PR URL and run summary, moving on to land only after `rhiza run approve`; `run request-changes` with one finding naming a leaf starts a fixer and that leaf's verifier, re-runs only the check it names as a new cycle of review's checks step journaled under the review workflow step of the same run, calls ship open-pr again, waits for CI and waits at sign-off again; both commands still exit 2 on a run not waiting at a sign-off step; an autonomous run reaches land with no wait and an entry with rule autonomous; a failing CI check ends the run failed naming the check, as today.

## statement

ship-change waits for the owner's sign-off once the PR's CI is green, `rhiza run status` naming the PR's URL and `rhiza run summary`, and this leaf widens `rhiza run approve` and `rhiza run request-changes <file>` from a review run waiting at sign-off to any run waiting at a sign-off step, ship-change's included: `rhiza run approve` moves on to rebase, land and merge, and `rhiza run request-changes <file>` journals the owner's findings as confirmed and runs one fix cycle outside the cap as review does (a fixer per named leaf, then its verifier), then goes back into the review workflow step at review's checks step, which re-runs only the failed checks with review's sign-off still skipped, after which ship open-pr pushes and replaces the PR body, the run waits for CI and returns to ship-change's sign-off; an autonomous run skips the wait and journals the skip with rule autonomous.

## Links

### Out

- depends-on [ship-open-pr](ship-open-pr.md)
- decided-by [two-owner-gates](../../decision/two-owner-gates.md)
- decided-by [autonomous-mode](../../decision/autonomous-mode.md)
- decided-by [answer-cb8bd934-2](../../decision/answer-cb8bd934-2.md)
- decided-by [answer-cb8bd934-3](../../decision/answer-cb8bd934-3.md)
- decided-by [answer-cb8bd934-4](../../decision/answer-cb8bd934-4.md)
- derived-from [ship-change-plan](../../source/ship-change-plan.md)

### In

- depends-on [ship-change-file](ship-change-file.md)
- verifies [ship-signoff-verified](../../evidence/ship-signoff-verified.md)
- verifies [ship-signoff-qa-1](../../evidence/ship-signoff-qa-1.md)
- verifies [ship-signoff-qa-3](../../evidence/ship-signoff-qa-3.md)
- verifies [ship-signoff-qa-4](../../evidence/ship-signoff-qa-4.md)
- verifies [ship-signoff-qa-5](../../evidence/ship-signoff-qa-5.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
