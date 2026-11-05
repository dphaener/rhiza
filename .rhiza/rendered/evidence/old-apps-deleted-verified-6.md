# old-apps-deleted-verified-6

evidence · `01M6NQ6S6QB73V9Z5T2F7C1BSQ` · actor [claude-code-1a012060](../actor/claude-code-1a012060.md)

## result

PASS

## statement

Verified by claude-code-1a012060 at 77527265. I ran go build ./... and go vet ./..., and both were clean. internal/tui now holds only autonomyapp, daemonapp, inboxapp, launcher, modalopen, modalshell, requestapp and runsapp. git grep finds no reference to changesapp, healthapp, modelapp, readapps, roadmapapp or writeapps in Go or the Makefile, and BUDGET_PKGS no longer lists internal/tui/modelapp. A scratch build of rhiza tui app with roadmap, model, health, changes, show, live, read, write and requests each exits 4 with rhiza: unknown app. The testscript tui-old-apps-gone passes, along with the unchanged scripts roadmap, roadmap-milestones, roadmap-paged, roadmap-umbrella, show, why, why-paged, slice, slice-change, slice-errors, health, health-machine-part, change and change-title. None of these scripts changed since the merge-base. rex-install and the tui-* scripts also pass. go test -trimpath passes for internal/tui/..., internal/api and internal/quiettest.

## Links

### Out

- verifies [old-apps-deleted](../capability/rhiza/old-apps-deleted.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
