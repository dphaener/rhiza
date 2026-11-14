# review-01m7b3wjc309vgwck1ef6ms5gb-r3-f1

question · `01M7B50GWMGJ7JK3R718NR4KZK`

## statement

Plan review round 3 left finding f1 (contradiction on web-token): web-token says the token's value is never shown, but it names only daemon.settings and GET /api/daemon.settings as places to hide it. The value also leaks through daemon.status. The daemon's poll copies every setting row into its state (conf.stateRows(), internal/api/daemon_settings.go:509, stored at internal/api/daemon_projects.go:168,241). daemon.status is a ReadOnly operation (internal/api/daemon.go:1620-1621) whose response returns those rows with their Value (Settings field, internal/api/daemon.go:1717, filled at 2057-2058). ReadOnlyHandler serves it with no authentication at GET /api/daemon.status (internal/api/web_api.go:18-30). Anyone who can load a page can read the token from there and send it as a bearer header. Neither the statement nor the acceptance covers daemon.status or the daemon state, so a builder who masks only daemon.settings passes while the token stays exposed. The critic confirmed it: Verified in the code. The daemon's poll copies every setting row with its Value into state (conf.stateRows(), internal/api/daemon_settings.go:509, stored at daemon_projects.go:168,241). daemon.status is registered ReadOnly (internal/api/daemon.go:1620-1621) and returns those rows' Value in Settings (daemon.go:1717, filled at 2057-2058). ReadOnlyHandler serves any ReadOnly op over unauthenticated GET /api/<op> (web_api.go:18-30). web-token's statement and acceptance name only daemon.settings and GET /api/daemon.settings, so a builder who masks only that op passes TestWebTokenSetting while GET /api/daemon.status still returns the token that gates every action. rhiza check does not report this. Major stands. How should it be settled?

## Links

### Out

- blocks [web-token](../capability/rhiza/web-monitor/web-token.md)

### In

- answers [web-token-masked-at-source](../decision/web-token-masked-at-source.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
