# feed-command

capability · `01M6B3PEPKV2KGQRZWN4PG4BYT` · parent [change-feed](../change-feed.md)

## acceptance

A feed.txtar lands three record change sets and sees them in land order; --limit 2 gives a next_cursor and the next page holds the third; a bad cursor exits 2; --since filters by land time; a change set reopened and landed again appears once at its new land; a change set with PR Evidence (pr field) shows its PR number; feed --json and api call feed print the same bytes; describe-golden.txtar is regenerated

## statement

rhiza feed lists landed change sets ordered by their latest land entry, each with ID, title, kind, leaf count, land time, actor and PR number when PR Evidence holds one, paged with --limit and a cursor, and filtered with --since <RFC 3339 time>

## Links

### Out

- decided-by [feed-lists-landed-by-land-entry](../../../../decision/feed-lists-landed-by-land-entry.md)

### In

- verifies [feed-command-verified](../../../../evidence/feed-command-verified.md)
- verifies [feed-command-qa-1](../../../../evidence/feed-command-qa-1.md)
- verifies [pr-167-change-feed-and-changes-app](../../../../evidence/pr-167-change-feed-and-changes-app.md)

## History

- `01M6B3NCDEX88PC3NFRW16YTVZ` rhiza feed lists what landed in land order with PR numbers, and a Rex changes app shows what landed since the owner last looked
