# old-apps-deleted-qa-1

evidence · `01M6NKHNAW9V5P7YGY5TX2E564` · actor [claude-code-c4ccbed3](../actor/claude-code-c4ccbed3.md)

## result

PASS

## statement

Built the branch head with go build -trimpath, then ran go vet ./... (clean). In a scratch repo, rhiza tui app roadmap, model, health, changes, show, live, why, slice, requests, read and write --once each exited 4 with 'unknown app'. inbox and runs --once exited 0. roadmap, health, change ls, show, why and slice ran on this repo's model (health exits 1 only for real reported problems). internal/tui now holds only autonomyapp, daemonapp, inboxapp, launcher, modalopen, modalshell, requestapp and runsapp, and no Go file, Makefile or txtar names the old packages. The health, roadmap, show, why, slice, tui-* and rex-install scripts pass (25), and the scripts of the wrapped commands are unchanged.

## Links

### Out

- verifies [old-apps-deleted](../capability/rhiza/old-apps-deleted.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
