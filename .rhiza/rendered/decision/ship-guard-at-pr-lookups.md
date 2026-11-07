# ship-guard-at-pr-lookups

decision · `01M6TCGVZFQEPRYXS60YJGP7ST`

## reasoning

The finding is right: internal/api registers no ship wait op, and the wait-pr steps in plan.toml and ship-change.toml are workflow wait steps whose until is 'rhiza ship check --settled'. Reading internal/api, the only places that look a PR up by branch on GitHub are runShipOpenPR (ship_openpr.go:125-150, gh pr view <branch>) and ship check without --main (ship.go:481, gh pr view <branch>), which records PR, URL and head into ship state (ship.go:558). Ship merge (ship.go:796-822) and land (ship_close.go:195-217) read st.PR from ship state and never search by branch; check --main reads the recorded merge commit. So guarding open-pr and check keeps every other op off a foreign PR, since a foreign PR can never reach ship state. In the reported failure the merged #345 was picked up by exactly this branch lookup. Naming check --settled explicitly keeps the wait-pr path covered by the acceptance without inventing an op a builder cannot test; adding guards to merge or rebase would duplicate the check against state that is already guarded.

## statement

ship-tracks-own-pr names the two ops that look a pull request up by branch, ship open-pr and ship check (with and without --settled, which is what the workflows' wait-pr steps run), not a 'wait' op: each refuses a PR on the branch, open, merged or closed, whose body's 'Change set <id>' line is missing or names another change set, with a conflict naming that PR, and writes nothing to ship state. Ship merge, land and close-out need no guard of their own: they act only on the PR recorded in ship state, which only open-pr and check write. Its acceptance tests open-pr, check and check --settled with a fake gh.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6tbdm-r3-f1](../question/review-01m6tbdm-r3-f1.md)

### In

- decided-by [ship-tracks-own-pr](../capability/rhiza/orchestration/shipping/ship-tracks-own-pr.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
