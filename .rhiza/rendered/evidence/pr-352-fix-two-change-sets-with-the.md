# pr-352-fix-two-change-sets-with-the

evidence · `01M6TMNB1YG8HQ2YMZD8R8JQBG` · actor [dphaener](../actor/dphaener.md)

## head

a6e219f5a1744016e38bac296515f8c0e2bf19c9

## pr

https://github.com/dphaener/rhiza/pull/352

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #352 (https://github.com/dphaener/rhiza/pull/352) passed CI at a6e219f

## status

verified

## Links

### Out

- verifies [unique-change-set-slug](../capability/rhiza/runtime/rex-sessions/unique-change-set-slug.md)
- verifies [slug-callers-agree](../capability/rhiza/runtime/rex-sessions/slug-callers-agree.md)
- verifies [ship-tracks-own-pr](../capability/rhiza/orchestration/shipping/ship-tracks-own-pr.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
