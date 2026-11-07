# ship-tracks-own-pr-qa-2

evidence · `01M6TJZDR2QV5Z9WXATCBQDKNJ` · actor [claude-code-78ae5b68](../actor/claude-code-78ae5b68.md)

## result

PASS

## statement

At 404f4b15 I built rhiza (go build -trimpath) into a scratch path. In a scratch repo with a bare origin, a fake gh and RHIZA_* unset apart from a planted review run, I gave the later of two same-titled change sets' branch a foreign commit on origin and a local commit. ship open-pr refused with 'pull request #345 ... does not name change set <id> ...' for a body naming the other change set and for a body with no Change set line, in OPEN, MERGED and CLOSED states. Each time gh saw only 'pr view'. HEAD and origin's branch did not move, the foreign commit was never fetched, and no ship state was written. ship check and ship check --settled refused the same six cases the same way. A foreign body arriving after ship state existed left the state file byte-identical (md5). With origin's branch cleared, open-pr with no PR pushed and ran gh pr create. With the PR's body naming the run's change set it ran gh pr edit 345 (created=false) and wrote .git/rhiza/ship/<id>.json with pr 345, the suffixed branch and the pushed head. ship check then tracked it as 'PR #345 at 3619007: not settled'.

## Links

### Out

- verifies [ship-tracks-own-pr](../capability/rhiza/orchestration/shipping/ship-tracks-own-pr.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
