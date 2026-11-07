# unique-change-set-slug-qa-1

evidence · `01M6THM2WZAZS428G8B92HM2D3` · actor [claude-code-1bbb24c8](../actor/claude-code-1bbb24c8.md)

## result

PASS

## statement

Built rhiza from the branch head and, in a scratch repo, made two change sets titled "Fix the usage governor's" (01M6THGEMJ…, 01M6THGETK…): the earlier got slug fix-the-usage-governor-s, the later fix-the-usage-governor-s-01m6thgetkzw (12 lowercased ID characters); a plan --new dry run for a third same-titled change set predicted fix-the-usage-governor-s-01m6thh4w48d

## Links

### Out

- verifies [unique-change-set-slug](../capability/rhiza/runtime/rex-sessions/unique-change-set-slug.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
