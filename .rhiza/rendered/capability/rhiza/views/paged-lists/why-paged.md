# why-paged

capability · `01M6H3MJKJC9D134K5JAMQHNSP` · parent [paged-lists](../paged-lists.md)

## acceptance

A test walks a node with decisions, sources and change sets and pages it with a --limit that splits the sequence across list boundaries: each page carries the node and view, the concatenated pages equal the --all answer in order, the last page has no next_cursor, and a bad cursor or out-of-range limit exits 2; why --json and api call why print the same bytes; describe-golden.txtar is regenerated

## statement

rhiza why computes its whole walk, then pages decisions, then sources, then change sets as one sequence with --limit counting entries across all three, an opaque --cursor naming the list and offset, and --all; node and view come on every page

## Links

### Out

- decided-by [decision-paging-log-roadmap-why](../../../../decision/decision-paging-log-roadmap-why.md)

### In

- verifies [why-paged-verified](../../../../evidence/why-paged-verified.md)
- verifies [why-paged-qa-1](../../../../evidence/why-paged-qa-1.md)
- verifies [pr-228-log](../../../../evidence/pr-228-log.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
