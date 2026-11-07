# unique-slug-plan-2026-11-07

source · `01M6TBRAH2H9J5J6KYPHMFX6C8`

## quotes

```json
[]
```

## raw_sha256

f42ba751fb2ba7ac99c46faa9a67931d1277d7829fbba3a87f8d782416b9fbfb

## session

f74b00a7-a63c-4e88-b98a-01a1b4b71fc9

## summary

Autonomous plan run for change set 01M6TBD; there was no owner. The planner read internal/session/names.go and its callers in internal/api (session.go, run.go, run_leaves.go, run_checks.go, ship.go, ship_openpr.go), where each one re-derives the branch and worktree from session.Slug(title). It asked one question, how to make slugs unique. The arbiter answered with Decision slug-unique-by-earliest-id: a change set whose title slug matches an earlier change set's (lower ID, any state) gets a suffix of the first 12 characters of its ID, lowercased, and one function over the model gives the slug to every caller. The planner settled the PR guard itself: ship tracks only a pull request whose body names this change set, using the 'Change set <id>' line that open-pr already writes, and refuses any other. The change set serves m29-one-daemon, as the sibling governor fixes do.

## title

Unique change set slugs, autonomous plan

## Links

### In

- derived-from [ship-pr-owned-by-body](../decision/ship-pr-owned-by-body.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
