# ship-draft-pr-stops-qa-1

evidence · `01M73JG4MJTSGB3T094F3RQ3Z2` · actor [claude-code-8a6947ea](../actor/claude-code-8a6947ea.md)

## result

PASS

## statement

Built rhiza from the branch head (go build ./cmd/rhiza) and drove a real ship-change run with it in a scratch Rhiza root, using run-ship-merge-draft.txtar's fake gh/rex/claude and its draft.sh step that turns PR #7 draft after rebased-merge-green. rhiza run approve --foreground stopped the run (Run ... stopped) at merge visit 1 with 'stop: pull request #7 is a draft: mark it ready for review, then resume the run'; run status showed Current step merge, no gate-rebase-merge/rebase-merge revisit, and gh.log had no pr merge and no pr ready. Directly, rhiza ship merge and ship merge --json exited 75 with that stop line; rhiza ship check exited 75 with the same stop and --json showed green false, stop set, pr 7. After removing the draft flag, rhiza run resume --foreground continued merge visit 1, made exactly one pr merge 7 --squash --match-head-commit <head>, never pr ready, and the run succeeded with the change set landed. Scoped tests: go test ./internal/api -run 'TestShipCheck|TestShipMerge' and TestScripts/(run-ship-merge-draft|run-ship-change|describe-golden) pass.

## Links

### Out

- verifies [ship-draft-pr-stops](../capability/ship-through-rhiza/ship-draft-pr-stops.md)

## History

- `01M73G0WSCCBG7H96NDY1H4KQ6` Fix watchdog unmerged PR at step gate-rebase-merge: run 01M73CQ4D71CW038ZBNSQB1SJF; log: PR green at 0b28569 and unmerged after step gate-rebase-merge failed: step "gate-rebase-merge" went over max_visits 6
