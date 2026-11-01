# changes-app

capability · `01M6B3PG5BQ496EZCTRDFMH0T2` · parent [change-feed](../change-feed.md)

## acceptance

A tui-changes-app.txtar: --once lists the unseen landed change sets; after marking them seen a second --once lists none, and a change set landed later appears; unit tests on a fake runner cover seen IDs and ordering; rex-install.txtar shows the rhiza_app_changes action on key g and no other key moved

## statement

A Rex app, changes, on key g of the rhiza key mode, shows landed change sets the viewer has not seen, newest first, with their PR numbers, and marks them seen in .git/rhiza/feed.json; --once prints the unseen list as JSON

## Links

### Out

- decided-by [changes-app-remembers-seen](../../../../decision/changes-app-remembers-seen.md)

### In

- verifies [changes-app-verified](../../../../evidence/changes-app-verified.md)
- verifies [changes-app-qa-1](../../../../evidence/changes-app-qa-1.md)
- verifies [pr-167-change-feed-and-changes-app](../../../../evidence/pr-167-change-feed-and-changes-app.md)

## History

- `01M6B3NCDEX88PC3NFRW16YTVZ` rhiza feed lists what landed in land order with PR numbers, and a Rex changes app shows what landed since the owner last looked
