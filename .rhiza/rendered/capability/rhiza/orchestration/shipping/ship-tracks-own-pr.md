# ship-tracks-own-pr

capability · `01M6TBRW1FQSQTEHBGNK10NGSV` · parent [shipping](../shipping.md)

## acceptance

internal/api ship tests with a fake gh: open-pr, check and check --settled each refuse an open or merged PR whose body names another change set, and one whose body has no 'Change set' line, and leave ship state unwritten (open-pr also neither pulls the foreign PR's commits in nor pushes onto it); a PR whose body names the run's change set is tracked as before

## statement

Ship open-pr and ship check (with and without --settled, which the workflows' wait-pr steps run), the ops that look a pull request up by branch, refuse a pull request on the change set's branch, open, merged or closed, unless its body's 'Change set <id>' line names the run's change set (a PR naming another change set or with no such line is refused), with a conflict naming that PR, before open-pr fetches, follows or pushes the branch, and never write it to ship state

## Links

### Out

- decided-by [ship-pr-owned-by-body](../../../../decision/ship-pr-owned-by-body.md)
- decided-by [ship-guard-at-pr-lookups](../../../../decision/ship-guard-at-pr-lookups.md)

### In

- blocks [review-01m6tbdm-r3-f1](../../../../question/review-01m6tbdm-r3-f1.md)
- verifies [ship-tracks-own-pr-verified](../../../../evidence/ship-tracks-own-pr-verified.md)
- verifies [ship-tracks-own-pr-qa-1](../../../../evidence/ship-tracks-own-pr-qa-1.md)
- verifies [ship-tracks-own-pr-verified-2](../../../../evidence/ship-tracks-own-pr-verified-2.md)
- verifies [ship-tracks-own-pr-qa-2](../../../../evidence/ship-tracks-own-pr-qa-2.md)
- verifies [pr-352-fix-two-change-sets-with-the](../../../../evidence/pr-352-fix-two-change-sets-with-the.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
