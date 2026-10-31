# roadmap-umbrella-group

capability · `01M686TX1JY61Z3V1YSRWCQRCR` · parent [roadmap](../roadmap.md)

## acceptance

A new roadmap-umbrella.txtar with an umbrella intent, a milestone intent and change sets serving the umbrella, both, or nothing: roadmap --json shows group kinds milestone, backlog, umbrella in that order, an umbrella change set appears in no other group, and the human output heads the group with the intent's handle; describe-golden.txtar is regenerated

## statement

A proposed or accepted change set serving an intent with a non-empty umbrella field is listed only in that intent's umbrella group, after the milestone groups and the backlog

## Links

### Out

- decided-by [umbrella-intent-field](../../../../decision/umbrella-intent-field.md)

### In

- verifies [roadmap-umbrella-group-verified](../../../../evidence/roadmap-umbrella-group-verified.md)
- verifies [roadmap-umbrella-group-qa-1](../../../../evidence/roadmap-umbrella-group-qa-1.md)
- verifies [pr-145-umbrella-areas-grouped](../../../../evidence/pr-145-umbrella-areas-grouped.md)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
