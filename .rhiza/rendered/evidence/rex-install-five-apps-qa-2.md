# rex-install-five-apps-qa-2

evidence · `01M6NP01TV2M23MSNAFRY33GJB` · actor [claude-code-37af0711](../actor/claude-code-37af0711.md)

## result

PASS

## statement

Ran the built binary's rhiza rex install with a fake rex on PATH and REX_CONFIG in a scratch folder: it wrote rhiza.lua and the require line and logged one 'config reload'. rhiza.lua holds exactly 5 rex.action entries, rhiza_app_autonomy, daemon, inbox, request and runs, with one key each (y, d, i, q, v) and no old app. rex-app-actions' acceptance now names the five apps. TestScripts/rex-install and internal/api -run RexInstall pass.

## Links

### Out

- verifies [rex-install-five-apps](../capability/rhiza/rex-install-five-apps.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
