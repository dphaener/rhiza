# rex-install-five-apps-qa-1

evidence · `01M6NKHV4M98CVVHVA11W40ZAA` · actor [claude-code-c4ccbed3](../actor/claude-code-c4ccbed3.md)

## result

PASS

## statement

Ran the built rhiza rex install with a fake rex on PATH and REX_CONFIG in a scratch dir. It exited 0 and reloaded the config once. rhiza.lua holds exactly 5 rex.action entries: rhiza_app_autonomy, rhiza_app_daemon, rhiza_app_inbox, rhiza_app_request and rhiza_app_runs, on keys y, d, i, q and v, with no old app names. go test ./internal/api -run TestRexInstallRegistersExactlyFiveApps and the rex-install testscript pass. tui app daemon --once exits 2 with a usage error; daemon is a form app.

## Links

### Out

- verifies [rex-install-five-apps](../capability/rhiza/rex-install-five-apps.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
