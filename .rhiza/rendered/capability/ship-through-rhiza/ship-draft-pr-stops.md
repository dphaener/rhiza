# ship-draft-pr-stops

capability · `01M73G9RGRGPCPWM9M1RVJDAF5` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShipCheck|TestShipMerge' passes with new tests where a fake gh reports isDraft true: ship check on a green draft PR exits 75 with a stop: line naming the PR as a draft and records no checked head; ship merge on a recorded head whose PR is now a draft exits 75 with that stop: line and never calls gh pr merge or gh pr ready; a ready PR still checks green and merges as before. go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-merge-draft' passes: a ship-change run whose PR turns draft after sign-off stops at merge (not failed, no visit past 1 at gate-rebase-merge from merge), and after the fake PR is marked ready, rhiza run resume merges it. describe-golden is regenerated if a response type changes.

## statement

When a ship-change run's pull request is a draft, rhiza ship check and rhiza ship merge stop the run (exit 75) with a stop: line naming the PR as a draft to mark ready before resuming; they never report it green, never merge it, never mark it ready (no gh pr ready) and never fail merge into gate-rebase-merge, so the merge -> gate-rebase-merge loop cannot reach max_visits.

## Links

### Out

- decided-by [decision-draft-pr-is-a-hold](../../decision/decision-draft-pr-is-a-hold.md)
- derived-from [draft-pr-merge-loop-plan](../../source/draft-pr-merge-loop-plan.md)

### In

- verifies [ship-draft-pr-stops-verified](../../evidence/ship-draft-pr-stops-verified.md)
- verifies [ship-draft-pr-stops-qa-1](../../evidence/ship-draft-pr-stops-qa-1.md)
- verifies [pr-432-fix-watchdog-unmerged-pr-at-01m73g0wsccb](../../evidence/pr-432-fix-watchdog-unmerged-pr-at-01m73g0wsccb.md)

## History

- `01M73G0WSCCBG7H96NDY1H4KQ6` Fix watchdog unmerged PR at step gate-rebase-merge: run 01M73CQ4D71CW038ZBNSQB1SJF; log: PR green at 0b28569 and unmerged after step gate-rebase-merge failed: step "gate-rebase-merge" went over max_visits 6
