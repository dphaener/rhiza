# ship-pr-owned-by-body

decision · `01M6TBRVXWBXEPRVDGGZ0DRPK4`

## reasoning

The intent requires that a run never commit to or track another change set's PR. The body line is already written by planSummary and reviewSummary, so ownership is checkable from gh pr view without new state. Unique slugs stop new collisions, and this guard catches the branches that already collide.

## statement

Ship tracks a pull request only when its body's 'Change set <id>' line, which open-pr already writes, names the run's change set; any other PR on the branch, open, merged or closed, is refused and never recorded in ship state

## Links

### Out

- derived-from [unique-slug-plan-2026-11-07](../source/unique-slug-plan-2026-11-07.md)

### In

- decided-by [ship-tracks-own-pr](../capability/rhiza/orchestration/shipping/ship-tracks-own-pr.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
