# old-apps-deleted-qa-2

evidence · `01M6NNZSQ77AV6DTQTVA2RK2MP` · actor [claude-code-37af0711](../actor/claude-code-37af0711.md)

## result

PASS

## statement

Built the branch head (go build -trimpath) and ran go vet ./... clean. With the built binary, rhiza tui app roadmap, model, health, changes, show and live --once each printed 'unknown app' and exited 4; inbox, runs and autonomy --once printed JSON and exited 0, request and daemon gave their no --once usage error (exit 2). The wrapped CLI commands still work: roadmap, slice and change ls exit 0, health prints its parts. The slice shows all 14 old-app capabilities (roadmap-app-*, model-tree-app, model-app-search, modelapp-timing-alone, changes-app, health-app, read-apps-match-cli, write-apps-previewed, runs-sessions-apps) retired, and show no longer finds them. TestScripts/tui-old-apps-gone passes.

## Links

### Out

- verifies [old-apps-deleted](../capability/rhiza/old-apps-deleted.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
