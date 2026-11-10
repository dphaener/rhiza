# decision-draft-pr-is-a-hold

decision · `01M73G84DARF6DHW7YREZ73H44`

## reasoning

The GitHub timeline for PR #426 shows the owner (not Rhiza) ran convert_to_draft at 2026-11-10T19:28Z, after ship opened the PR as ready. That is a deliberate human signal not to merge, so marking it ready again would override the owner, and retrying cannot succeed: runShipMerge (internal/api/ship.go) calls gh pr merge, which fails with 'Pull Request is still a draft', and the merge step's on.fail sends the run to gate-rebase-merge, which loops until max_visits 6. That is the watchdog failure this change set fixes. A stop: line (ShipCheckResponse already has Stop) ends the run on an actionable reason, the way main-green stops, instead of burning visits or waiting forever on a PR nobody may mark ready. Checking in ship check as well puts the reason where the run first sees the PR.

## statement

A draft pull request is a hold, not a failure to retry. rhiza ship merge never merges a draft and never marks it ready itself (no gh pr ready): it merges nothing and ends the run with a stop: line naming the PR and saying it is a draft, to be marked ready for review before the run is resumed. It does not fail into gate-rebase-merge, so the merge -> gate-rebase-merge loop cannot reach max_visits. rhiza ship check reads isDraft too and reports a draft PR as that same stop instead of green, so the run stops before merge.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-b8f88971-1](../question/question-b8f88971-1.md)
- derived-from [draft-pr-merge-loop-plan](../source/draft-pr-merge-loop-plan.md)

### In

- decided-by [ship-draft-pr-stops](../capability/ship-through-rhiza/ship-draft-pr-stops.md)

## History

- `01M73G0WSCCBG7H96NDY1H4KQ6` Fix watchdog unmerged PR at step gate-rebase-merge: run 01M73CQ4D71CW038ZBNSQB1SJF; log: PR green at 0b28569 and unmerged after step gate-rebase-merge failed: step "gate-rebase-merge" went over max_visits 6
