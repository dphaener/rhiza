# rex-install-five-apps-verified

evidence · `01M6NK7KNAR9991W5SW5CVG5W7` · actor [claude-code-8e6929de](../actor/claude-code-8e6929de.md)

## result

PASS

## statement

At d370f954, go test -trimpath -count=1 ./internal/api -run TestRexInstall passed, including the new TestRexInstallRegistersExactlyFiveApps, which runs rex.install with the five apps linked and asserts exactly five rex.action entries named rhiza_app_inbox, rhiza_app_request, rhiza_app_runs, rhiza_app_daemon and rhiza_app_autonomy, each with its key (i, q, v, d, y) in the key mode. go test -trimpath ./cmd/rhiza -run TestScripts/rex-install passed: with a fake rex it greps the five action names and open_argv lines, grep -count=5 for rex.action and for rhiza_app_ names, and each key line incl. daemon = d. Read the diff: daemonapp registers with the launcher (key d) and rhiza tui app daemon opens it; daemonapp and launcher tests pass, vet and gofmt clean.

## Links

### Out

- verifies [rex-install-five-apps](../capability/rhiza/rex-install-five-apps.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
