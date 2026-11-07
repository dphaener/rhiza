# unique-change-set-slug-qa-2

evidence · `01M6TJZDTVM284XENYWQM18J60` · actor [claude-code-78ae5b68](../actor/claude-code-78ae5b68.md)

## result

PASS

## statement

At 404f4b15 by hand with the built binary: of two same-titled change sets, the earlier (lower ID) kept the bare slug fix-the-usage-governor-s. The later got fix-the-usage-governor-s-01m6tjtqvq15, the first 12 characters of its ID lowercased, not its display short ID. The slug held across session start, run start dry runs, ship and abandon, so it is the same on repeated calls.

## Links

### Out

- verifies [unique-change-set-slug](../capability/rhiza/runtime/rex-sessions/unique-change-set-slug.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
