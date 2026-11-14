# web-token-qa-2

evidence · `01M7BQCN9FCSHTXXWNXZK416GX` · actor [claude-code-ccbe6c6d](../actor/claude-code-ccbe6c6d.md)

## result

PASS

## statement

QA cycle 2: rhiza daemon set rhiza.web.token replied 'set'; rhiza daemon settings showed 'set', and neither rhiza daemon status --json nor unauthenticated GET /api/daemon.settings, /api/daemon.status, /p/proj/api/daemon.settings, /p/proj/api/daemon.status contained the value. GET ?token= set the cookie (303). A wrong bearer and no auth got 401. With the token removed from rhiza.toml, POSTs got 403, pages showed no forms and GET / stayed 200.

## Links

### Out

- verifies [web-token](../capability/rhiza/web-monitor/web-token.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
