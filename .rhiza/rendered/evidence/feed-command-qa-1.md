# feed-command-qa-1

evidence · `01M6BEFE41WQ9CR3Y0P1WJBS73` · actor [claude-code-d8740dd0](../actor/claude-code-d8740dd0.md)

## result

PASS

## statement

QA at 52027dba with a binary built from HEAD, on this repo's real .rhiza: rhiza feed --limit 3 listed the oldest landed record change sets in land order with ID, kind, actor, title and leaf count, and printed a next: line with a cursor; --json gave next_cursor and land_entry/landed_at per item; 62 of the landed change sets carry pr (e.g. The review workflow pr 45, Abandon clears worktree pr 164); --cursor nonsense and --since yesterday each exited 2 with next: rhiza describe feed; --since 2026-10-31T00:00:00Z started at that day's lands; feed --limit 5 --json and api call feed {"limit":5} were byte-identical (cmp). In a scratch clone a newly landed record appeared in feed --since at its land time. go test feed, describe-golden and internal/api passed.

## Links

### Out

- verifies [feed-command](../capability/rhiza/views/change-feed/feed-command.md)

## History

- `01M6B3NCDEX88PC3NFRW16YTVZ` rhiza feed lists what landed in land order with PR numbers, and a Rex changes app shows what landed since the owner last looked
