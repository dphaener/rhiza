# old-apps-deleted

capability · `01M6MZD2ER8DDHNYK3W3C04SYB` · parent [rhiza](../rhiza.md)

## acceptance

go build ./... and go vet ./... pass with the packages removed; a testscript runs rhiza tui app roadmap, model, health, changes, show and live and each exits 4 with unknown app; the existing testscripts of the commands those apps wrapped (roadmap, show, why, slice, health, change ls) pass unchanged.

## statement

The changes, health, model, read, roadmap, write, live and requests apps are gone: their packages under internal/tui, their branches in cmd/rhiza/tui.go, their tests and testscripts, and internal/tui/modelapp in the Makefile's budget packages; the CLI commands they wrapped stay.

## Links

### Out

- decided-by [tui-delete-old-apps](../../decision/tui-delete-old-apps.md)

### In

- blocks [leaf-failed-01m6njpm-1](../../question/leaf-failed-01m6njpm-1.md)
- blocks [leaf-failed-01m6njpm-2](../../question/leaf-failed-01m6njpm-2.md)
- verifies [old-apps-deleted-verified](../../evidence/old-apps-deleted-verified.md)
- verifies [old-apps-deleted-qa-1](../../evidence/old-apps-deleted-qa-1.md)
- verifies [old-apps-deleted-verified-3](../../evidence/old-apps-deleted-verified-3.md)
- verifies [old-apps-deleted-qa-2](../../evidence/old-apps-deleted-qa-2.md)
- verifies [old-apps-deleted-verified-6](../../evidence/old-apps-deleted-verified-6.md)
- verifies [old-apps-deleted-verified-8](../../evidence/old-apps-deleted-verified-8.md)
- verifies [pr-301-delete-the-old-rex-apps](../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
