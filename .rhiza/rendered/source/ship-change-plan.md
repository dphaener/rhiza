# ship-change-plan

source · `01M5Z8WDQW7G4Q29WE96XWJEH9`

## quotes

```json
[]
```

## raw_sha256

4c1bfb7f530b44b63648b3254cf132495d2d88cd0d27b4a8ba9a7ce0e4f0db69

## session

cb8bd934-a417-493a-aebe-01a17e247ca8

## summary

The plan step of a plan run on change set 8 of milestone 3 (ship-change on Rhiza's workflows), 2026-10-27. Starting from the settled milestone-3 decisions (four-workflows, two-owner-gates, engine-is-conductor, autonomous-mode, review-workflow-shape), the build and review change sets' planned leaves (the leaves and checks step kinds, rhiza run summary, review's sign-off with rhiza run approve and request-changes) and today's ship-change.toml (a /shipit:go agent step, then wait, check, rebase, land, merge and close-out steps), the planner put two topics to the owner through AskUserQuestion, each with a recommendation; the owner picked the recommended option every time. Settled: (1) ship-change runs build and review through a new `workflow` step kind that runs another workflow's steps inside the same run and journal, so repo overrides of build or review apply when shipping; (2) the PR opens before the owner's sign-off: build, then review's checks without review's own sign-off, then push and open the PR with rhiza run summary as its body, wait for CI, then ship-change's own sign-off, then rebase, land and merge as today; (3) request-changes at that sign-off runs one fix cycle as in review, re-runs only the failed checks, pushes, updates the PR body, waits for CI and returns to the sign-off; (4) a failing CI check fails the run as today, to be fixed and resumed; (5) the shipit path is replaced, not kept, and a change set without leaves is refused at run start, pointing to the plan workflow. No answer was typed, so there are no owner quotes; the picks are captured as Decisions by session interview; the leaves link the landed milestone-3 decisions.

## title

ship-change on Rhiza's workflows, planned with the owner

## Links

### In

- derived-from [ship-through-rhiza](../capability/ship-through-rhiza.md)
- derived-from [workflow-step](../capability/ship-through-rhiza/workflow-step.md)
- derived-from [ship-open-pr](../capability/ship-through-rhiza/ship-open-pr.md)
- derived-from [ship-signoff](../capability/ship-through-rhiza/ship-signoff.md)
- derived-from [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)
- derived-from [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)
- derived-from [answer-cb8bd934-1](../decision/answer-cb8bd934-1.md)
- derived-from [answer-cb8bd934-2](../decision/answer-cb8bd934-2.md)
- derived-from [answer-cb8bd934-3](../decision/answer-cb8bd934-3.md)
- derived-from [answer-cb8bd934-4](../decision/answer-cb8bd934-4.md)
- derived-from [answer-cb8bd934-5](../decision/answer-cb8bd934-5.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
