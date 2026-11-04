# log-paged

capability · `01M6H3M6MXZQMD7EMF9FPNGPE0` · parent [paged-lists](../paged-lists.md)

## acceptance

A testscript logs a node touched by three change sets: --limit 2 gives two change sets, a next_cursor and a 'next: rhiza log <ref> --limit 2 --cursor <c>' line, and that cursor's page holds the third with no next_cursor; --all gives all three; --limit 0 gives the default page as ls does; --limit -1, --limit 1001 and a bad cursor exit 2; log --json and api call log print the same bytes; describe-golden.txtar is regenerated

## statement

rhiza log pages its change sets in apply order with --limit (1 to 1000, default 100) and an opaque --cursor, returns next_cursor while more remain and ends its text with a next: line, and --all returns every change set; target and view come on every page

## Links

### Out

- decided-by [decision-paging-log-roadmap-why](../../../../decision/decision-paging-log-roadmap-why.md)

### In

- verifies [log-paged-verified](../../../../evidence/log-paged-verified.md)
- verifies [log-paged-qa-1](../../../../evidence/log-paged-qa-1.md)
- verifies [pr-228-log](../../../../evidence/pr-228-log.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
