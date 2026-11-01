# feed-command-verified

evidence · `01M6BD2YKXEFH99N76JVVSHFY8` · actor [claude-code-75f68d60](../actor/claude-code-75f68d60.md)

## result

PASS

## statement

At b8a3cbbf: go test ./cmd/rhiza -run 'TestScripts/(feed|describe-golden)$' -count=1 passes. feed.txtar lands alpha/beta/gamma record change sets and sees them in land order; --limit 2 emits next_cursor and the cursor page holds beta and gamma with no next; bad cursor and bad limit give code usage; --since 2026-10-26 drops alpha, 2026-10-28 gives Nothing landed, a non-RFC-3339 since is usage; Evidence with pr=.../pull/42 shows "pr": 42; a reopened and re-landed change set appears once at its 2026-10-30 land; feed --limit 2 --json and api call feed {"limit":2} cmp equal. describe-golden.txtar includes feed and passes. Built binary by hand: rhiza feed --cursor nonsense exits 2 with next: rhiza describe feed. make fmt-check vet test pass; make self-check passes with RHIZA_ROOT/CHANGE/RUN unset. Read internal/api/feed.go: ordered by latest land entry ID, cursor is base64url of that entry ID.

## Links

### Out

- verifies [feed-command](../capability/rhiza/views/change-feed/feed-command.md)

## History

- `01M6B3NCDEX88PC3NFRW16YTVZ` rhiza feed lists what landed in land order with PR numbers, and a Rex changes app shows what landed since the owner last looked
