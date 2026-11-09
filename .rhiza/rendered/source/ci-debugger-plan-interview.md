# ci-debugger-plan-interview

source · `01M6YCQ68SZB48PD6ESRDXY20F`

## quotes

```json
[
  {
    "by": "owner",
    "text": "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci",
    "where": "change set intent, 2026-11-08"
  },
  {
    "by": "owner",
    "text": "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense",
    "where": "change set intent, 2026-11-08"
  }
]
```

## raw_sha256

c074af1b79018f8c2b240c7eab18f7be9a33f06f33859cdc20ab8cea7c1638e3

## session

25f669da-d907-4a8d-8044-01a1bcc76cfc

## summary

Autonomous plan run for change set 01M6YCDW1VM9R8KZPGQR330MAY. No owner was present; the planner read the intent (which quotes the owner's 2026-11-08 words), milestone-25, the ship-change workflow, rhiza ship check's rerun budget, the watchdog's CI rerun and bug filing, and rhiza health. Existing code already reruns a head's failed jobs at most once (ship check and watchdog) and sends a red head with failing tests to one ci-fix fixer that re-runs only those tests. Gaps found: flakes (red that passes on rerun) are not tracked or filed (the watchdog explicitly files no bug for a CI flake), rhiza health shows no CI stability, there is no full local gate before a push, and ci-fix pushes without reproducing locally. Four questions went to the arbiter as one; Decision ci-unstable-gate-rules answered them: serve milestone-25; CI is unstable while the last 20 completed CI runs (a project setting) hold a flake or main's latest run is red; a project-configured local gate command (make check here) runs before every push while unstable, with a health warning when none is set; a flake does not hold the run but is recorded per test and filed once per test as a bug; a head still red after its one rerun is reproduced and fixed or quarantined locally before the next push.

## title

Stop using CI as a debugger: autonomous plan interview

## Links

### In

- derived-from [ci-unstable-gate-rules](../decision/ci-unstable-gate-rules.md)
- derived-from [no-ci-as-debugger](../decision/no-ci-as-debugger.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
