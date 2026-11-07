# review-01m6tbdm-r3-f1

question · `01M6TCFJY75K5P3Z54PQTH334Y`

## statement

Plan review round 3 left finding f1 (contradiction on ship-tracks-own-pr): The statement and acceptance name three ship operations, 'open-pr, check and wait', but there is no ship wait op: internal/api registers ship.open-pr, ship.check, ship.merge, ship.rebase, ship.land, ship.close-out, ship.command, ship.deploy, ship.sh. The wait-pr steps in internal/workflow/builtin/plan.toml and ship-change.toml are workflow wait steps that run 'rhiza ship check --settled'. The acceptance therefore names a test of an op that does not exist; it should name ship open-pr and ship check (including --settled), or say which other ops (merge, rebase) also read the PR. The critic confirmed it: internal/api registers ship.open-pr (ship_openpr.go:23), ship.check and ship.merge (ship.go:40,48), ship.rebase, ship.land, ship.close-out, ship.command and ship.deploy, and no ship wait op exists anywhere in internal/ or cmd/. The wait-pr steps in plan.toml (lines 79, 99) are workflow 'wait' steps whose until is 'rhiza ship check --settled'. So ship-tracks-own-pr's statement and acceptance name an op ('wait') that a builder cannot test; they should name ship open-pr and ship check (incl. --settled), plus any other op that reads the PR. How should it be settled?

## Links

### Out

- blocks [ship-tracks-own-pr](../capability/rhiza/orchestration/shipping/ship-tracks-own-pr.md)

### In

- answers [ship-guard-at-pr-lookups](../decision/ship-guard-at-pr-lookups.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
