# old-apps-deleted-verified

evidence · `01M6NK4MVVT0J9P01S0Z106V2B` · actor [claude-code-593b7e06](../actor/claude-code-593b7e06.md)

## result

PASS

## statement

At b9b5c5ab: go build ./... and go vet ./... pass; internal/tui keeps only autonomyapp, daemonapp, inboxapp, launcher, modalopen, modalshell, requestapp, runsapp, and git grep finds no reference to changesapp, healthapp, modelapp, readapps, roadmapapp or writeapps; the Makefile's BUDGET_PKGS no longer lists internal/tui/modelapp. A scratch build of rhiza tui app roadmap|model|health|changes|show|live --once each printed unknown app and exited 4, and the new tui-old-apps-gone testscript passes. The command testscripts health, health-machine-part, roadmap, roadmap-milestones, roadmap-paged, roadmap-umbrella, show, why, why-paged, slice, slice-budget, slice-change, slice-errors, change, change-title, lifecycle and apply pass and are unchanged in the diff; text-refs-display (moved from tui app show/slice/why to the commands), rex-install, the remaining tui-* scripts, TestMakeTestTimeout and the internal/tui and quiettest packages pass too.

## Links

### Out

- verifies [old-apps-deleted](../capability/rhiza/old-apps-deleted.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
