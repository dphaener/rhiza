# feed-hides-actor-records

capability · `01M6BH9A4779HN5ZTFQGW718ZT` · parent [feed-command](../feed-command.md)

## acceptance

feed.txtar: by default an actor-only record is hidden while a planning record and an intention are shown, and with --all all three are shown; with --limit the pages are full and the next-page line carries --all; describe-golden.txtar is regenerated

## statement

rhiza feed leaves out landed change sets whose ops only create or update actor nodes unless --all is given, filtering before paging so pages stay full and cursors stable

## Links

### Out

- decided-by [feed-hides-actor-only](../../../../../decision/feed-hides-actor-only.md)

### In

- verifies [feed-hides-actor-records-verified](../../../../../evidence/feed-hides-actor-records-verified.md)
- verifies [feed-hides-actor-records-qa-1](../../../../../evidence/feed-hides-actor-records-qa-1.md)
- verifies [pr-176-feed-hides-actor-records](../../../../../evidence/pr-176-feed-hides-actor-records.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6BH7XPH07NSDTBKE2RGX949` rhiza feed and the changes app leave out actor-only change sets unless --all is given
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
