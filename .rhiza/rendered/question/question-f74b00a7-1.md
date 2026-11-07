# question-f74b00a7-1

question · `01M6TBHMBQ5DZ7XVV6GW7Q1SWV`

## session

f74b00a7-a63c-4e88-b98a-01a1b4b71fc9

## statement

Change set 01M6TBD fixes two change sets with the same title slug sharing a worktree and branch. How should a change set's slug be made unique? Note that every caller (worktree, branch, ship target, abandon cleanup, run tabs) re-derives the slug from the title today, so the rule must give the same slug every time it runs.

## Links

### In

- answers [slug-unique-by-earliest-id](../decision/slug-unique-by-earliest-id.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
