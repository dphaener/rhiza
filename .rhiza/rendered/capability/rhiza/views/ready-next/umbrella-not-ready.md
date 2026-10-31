# umbrella-not-ready

capability · `01M686TYF6800AVCAKJPHMAAM7` · parent [ready-next](../ready-next.md)

## acceptance

In roadmap-umbrella.txtar ready --json gives that reason and its next is a non-umbrella change set; a case in internal/api's daemon tests shows an umbrella group with leaves yields no candidate

## statement

An umbrella change set is never ready, with the reason umbrella: work is carved out of it, so ready's next and the daemon skip it

## Links

### Out

- decided-by [umbrella-intent-field](../../../../decision/umbrella-intent-field.md)

### In

- verifies [umbrella-not-ready-verified](../../../../evidence/umbrella-not-ready-verified.md)
- verifies [umbrella-not-ready-qa-1](../../../../evidence/umbrella-not-ready-qa-1.md)
- verifies [pr-145-umbrella-areas-grouped](../../../../evidence/pr-145-umbrella-areas-grouped.md)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
