# roadmap-app-umbrella-row

capability · `01M686TZYQZXDVR5B3880M6C21` · parent [roadmap-app-list](../roadmap-app-list.md)

## acceptance

A test in internal/tui/roadmapapp with a fake runner returning an umbrella group shows one row while collapsed and the change set rows after enter; --once output is unchanged

## statement

The roadmap app shows each umbrella group as one collapsed row naming how many change sets it holds, and enter expands or collapses it

## Links

### Out

- decided-by [umbrella-intent-field](../../../decision/umbrella-intent-field.md)

### In

- verifies [roadmap-app-umbrella-row-verified](../../../evidence/roadmap-app-umbrella-row-verified.md)
- verifies [roadmap-app-umbrella-row-qa-1](../../../evidence/roadmap-app-umbrella-row-qa-1.md)
- verifies [pr-145-umbrella-areas-grouped](../../../evidence/pr-145-umbrella-areas-grouped.md)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
