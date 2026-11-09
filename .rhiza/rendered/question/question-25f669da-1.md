# question-25f669da-1

question · `01M6YCK5FPCMN68J174FVPHQ23`

## session

25f669da-d907-4a8d-8044-01a1bcc76cfc

## statement

Change set 'Stop using CI as a debugger' serves no intent yet, and the daemon builds only change sets that serve a milestone. Which should it serve? (milestone-25 is 'the backlog is empty: the daemon builds and ships on its own'; its sibling 'Ship-change fixes red PR CI' served it.) / When does Rhiza count a repo's CI as unstable (which turns on the full local gate before each push and shows red in rhiza health)? / What is the 'full local gate' a ship-change run runs before pushing while CI is unstable? / When rhiza ship check sees a red PR that passes on its one rerun (a flake), what happens to the run and the flaky test?

## Links

### In

- answers [ci-unstable-gate-rules](../decision/ci-unstable-gate-rules.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
