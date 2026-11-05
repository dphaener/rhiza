# changes-app-qa-1

evidence · `01M6BEFFNJCP8PMRHWNXV0W1M0` · actor [claude-code-d8740dd0](../actor/claude-code-d8740dd0.md)

## result

PASS

## statement

QA at 52027dba with a binary built from HEAD, in a scratch clone of this repo: the first rhiza tui app changes --once listed all 1230 landed change sets newest first (PR numbers present, e.g. pr 164), wrote .git/rhiza/feed.json with sorted seen IDs, and a second --once printed change_sets []. From a git worktree of the clone, a record landed there showed up alone in --once and was gone on the next run, so the worktree shares the clone's feed.json. Then I ran the interactive app in a real Rex window: it rendered '1 landed since you last looked' with the newly landed qa-probe-three row (time, ID, title) and a 'j/k scroll · q quit' hint. Pressing q closed it, and a later --once listed none, with the ID now in feed.json. go test tui-changes-app, rex-install and internal/tui/changesapp passed.

## Links

### Out

- verifies changes-app (retired)

## History

- `01M6B3NCDEX88PC3NFRW16YTVZ` rhiza feed lists what landed in land order with PR numbers, and a Rex changes app shows what landed since the owner last looked
