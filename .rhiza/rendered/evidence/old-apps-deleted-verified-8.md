# old-apps-deleted-verified-8

evidence · `01M6NQXYG8VRGQVVQK9M57CV4G` · actor [claude-code-1fbbf353](../actor/claude-code-1fbbf353.md)

## result

PASS

## statement

At 510159ce: go build ./... and go vet ./... pass. internal/tui now holds only autonomyapp, daemonapp, inboxapp, launcher, modalopen, modalshell, requestapp, runsapp; no Go file or the Makefile references changesapp, healthapp, modelapp, readapp, roadmapapp, writeapps, liveapp or requestsapp, and BUDGET_PKGS no longer lists internal/tui/modelapp. cmd/rhiza/tui.go accepts only inbox, runs, request, autonomy and daemon. A built binary run as `rhiza tui app <x> --once` for roadmap, model, health, changes, show and live exits 4 with 'unknown app' each; testscript tui-old-apps-gone passes. The wrapped commands' scripts (roadmap, roadmap-milestones, roadmap-paged, roadmap-umbrella, show, why, why-paged, slice, slice-change, slice-errors, slice-budget, health, health-machine-part, change, change-title) are unchanged against the merge-base and pass with go test -trimpath ./cmd/rhiza -count=1 -run TestScripts/...; go test ./internal/tui/... passes.

## Links

### Out

- verifies [old-apps-deleted](../capability/rhiza/old-apps-deleted.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
