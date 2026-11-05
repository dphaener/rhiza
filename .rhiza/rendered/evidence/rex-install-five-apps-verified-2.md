# rex-install-five-apps-verified-2

evidence · `01M6NMDT6CNZ9ADESVNP5BNGKK` · actor [claude-code-8ceedb43](../actor/claude-code-8ceedb43.md)

## result

PASS

## statement

claude-code-8ceedb43 re-verified rex-install-five-apps at 96460ee5 (fix commit after d370f954). go test -trimpath -count=1 ./internal/api -run TestRexInstall: TestRexInstallRegistersExactlyFiveApps passes; it counts exactly five rex.action calls named rhiza_app_inbox, _request, _runs, _daemon and _autonomy, each with its key (i, q, v, d, y) in the rhiza key mode. go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/rex-install$' passes; the script uses a fake rex, greps the five action names and open_argv lines, and checks grep -count=5 for rex.action and for rhiza_app_ names. Diff read: daemonapp registers in the launcher with key d, and rhiza tui app daemon opens it. The fix commit only records rex-app-actions' acceptance for the five apps; make self-check at a clean checkout of 96460ee5 prints ok.

## Links

### Out

- verifies [rex-install-five-apps](../capability/rhiza/rex-install-five-apps.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
