# slug-callers-agree-qa-2

evidence · `01M6TJZDSDF8JS484R3PRB9CJR` · actor [claude-code-78ae5b68](../actor/claude-code-78ae5b68.md)

## result

PASS

## statement

At 404f4b15, with the built binary in a scratch Rhiza root and project, two change sets titled "Fix the usage governor's" (01M6TJTQRV..., then 01M6TJTQVQ..., sharing an 8-character prefix) got the following from rhiza session start: worktrees repo.fix-the-usage-governor-s and repo.fix-the-usage-governor-s-01m6tjtqvq15, and branches cs/fix-the-usage-governor-s and cs/fix-the-usage-governor-s-01m6tjtqvq15. run start --dry-run on the later one predicted the suffixed worktree. A --new plan dry run with the same title predicted repo.fix-the-usage-governor-s-01m6tjvdcnnk. ship open-pr and ship check on the later one targeted the suffixed branch. rhiza run abandon of a clean run on the later change set removed only the suffixed worktree and deleted only cs/fix-the-usage-governor-s-01m6tjtqvq15; the earlier one's worktree and branch stayed. Scoped go tests for internal/session, internal/api (Ship|Slug|Abandon) and the change-worktree-unique-slug, ship-follows-pr-head and run-ship-change scripts passed.

## Links

### Out

- verifies [slug-callers-agree](../capability/rhiza/runtime/rex-sessions/slug-callers-agree.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
