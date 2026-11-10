# draft-pr-merge-loop-plan

source · `01M73G9F36A10TYCW9C8ZSAPWH`

## quotes

```json
[]
```

## raw_sha256

e3788d692b4888adddfe97d3406118736e26153350e875ae093e9ca01bfb0221

## session

b8f88971-6f9a-4c0c-b545-01a1c700dbf9

## summary

Autonomous plan of the watchdog bug from run 01M73CQ4D71CW038ZBNSQB1SJF. The planner read the run log: every merge visit failed with gh's 'Pull Request is still a draft' on PR #426, and merge's on.fail sent the run back to gate-rebase-merge until it went over max_visits 6; the watchdog then retried twice from the same step. The PR timeline shows it was converted to draft at 2026-11-10T19:28Z, after ship opened it ready, while rhiza ship check (which reads no isDraft) still reported it green. The planner raised one question (stop for the owner, or mark the PR ready and merge); with no owner present, the arbiter answered with Decision decision-draft-pr-is-a-hold: a draft PR is a hold, ship merge and ship check end the run with a stop: line naming the draft PR and never mark it ready. The change set serves m4-runs-hold-up and has one leaf, ship-draft-pr-stops.

## title

Draft PR merge loop, planned autonomously

## Links

### In

- derived-from [decision-draft-pr-is-a-hold](../decision/decision-draft-pr-is-a-hold.md)
- derived-from [ship-draft-pr-stops](../capability/ship-through-rhiza/ship-draft-pr-stops.md)

## History

- `01M73G0WSCCBG7H96NDY1H4KQ6` Fix watchdog unmerged PR at step gate-rebase-merge: run 01M73CQ4D71CW038ZBNSQB1SJF; log: PR green at 0b28569 and unmerged after step gate-rebase-merge failed: step "gate-rebase-merge" went over max_visits 6
