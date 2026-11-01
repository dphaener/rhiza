# wait-settles-conflicting-pr

capability · `01M6A0NG1VRQKBM20M8CSVTX3B` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A ship-check-settles-conflict.txtar with a fake gh whose pr view reports mergeable CONFLICTING until a new head is pushed: with origin/main conflicting only under .rhiza/rendered, ship check --settled exits non-zero, origin's branch then contains origin/main and the ship record's pushed is the new head; with a conflict in a code file it exits 75, prints a stop: line, raises a rebase-conflict question and leaves the branch unchanged

## statement

rhiza ship check --settled reads the PR's mergeable state and, when GitHub reports it CONFLICTING, rebases the branch onto origin/main as rhiza ship rebase does, settling rendered-only conflicts by rendering again and recording the pushed head, then reports not settled; a conflict in code raises the rebase-conflict question, leaves the branch as it was and exits 75 with a stop: line

## Links

### Out

- decided-by [wait-settles-conflict-rule](../../../../../decision/wait-settles-conflict-rule.md)

### In

- verifies [wait-settles-conflicting-pr-verified](../../../../../evidence/wait-settles-conflicting-pr-verified.md)
- verifies [wait-settles-conflicting-pr-qa-1](../../../../../evidence/wait-settles-conflicting-pr-qa-1.md)
- verifies [pr-153-wait-settles-conflicting-pr](../../../../../evidence/pr-153-wait-settles-conflicting-pr.md)

## History

- `01M6A0JJDJ0FXD58K4WV88FV0A` A ship-change wait settles a PR that conflicts after main moves by rebasing it itself, and stops for the owner on a code conflict
