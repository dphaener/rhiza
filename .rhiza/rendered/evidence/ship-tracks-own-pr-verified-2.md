# ship-tracks-own-pr-verified-2

evidence · `01M6TJQWXPQNNSKF3G5BVD2SAH` · actor [claude-code-ff19fa3c](../actor/claude-code-ff19fa3c.md)

## result

PASS

## statement

At 1213886b I ran go test -trimpath -count=1 ./internal/api/ -run TestShip, and every ship test passed with the fake gh. TestShipOpenPRRefusesForeignPR covers a PR body that names another change set and a body with no 'Change set' line, each OPEN, MERGED and CLOSED, with the foreign commit both ahead on origin and absent from it. In every case it gets exit 3 with a conflict naming #7, no pr create or pr edit, no ship state, HEAD unmoved (nothing pulled in) and origin's branch unmoved (nothing pushed). TestShipCheckRefusesForeignPR covers ship check with {} and with settled:true on the same bodies plus an empty body, in all three states: each is refused with a conflict naming #345 and the state file is not written. TestShipOpenPR and TestShipCheckGreenThenMerge pass, so a PR that names the run's change set is still tracked. Reading the diff, open-pr now looks the PR up and calls requirePROwned before shipFollowPR and shipPush, and ship check calls requirePROwned before it writes state. As a mutation check, I restored the old ship_openpr.go in a scratch worktree, and the new test failed on HEAD moving and on origin's branch moving.

## Links

### Out

- verifies [ship-tracks-own-pr](../capability/rhiza/orchestration/shipping/ship-tracks-own-pr.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
