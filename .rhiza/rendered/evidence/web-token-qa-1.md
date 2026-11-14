# web-token-qa-1

evidence · `01M7BK29SC4388S7B1Z3WRV0SA` · actor [claude-code-a17d514a](../actor/claude-code-a17d514a.md)

## result

PASS

## statement

QA cycle 1: built the branch to a scratch path and ran it as a scratch root daemon with the web UI on a local port. rhiza daemon set rhiza.web.token refused 'short' and a value with spaces (exit 2) and accepted a 26-character token, replying value 'set'; daemon settings (human and --json), daemon status --json, GET /api/daemon.settings, GET /api/daemon.status (root and /p/proj/) and the daemon state file never held the value, only 'set'. Without the token every POST got 401 (wrong or missing bearer/cookie); GET ?token=<token> set an HttpOnly SameSite=Strict cookie and 303-redirected without the token; a wrong ?token= set nothing. Pages showed no action forms without the cookie. With the token removed from rhiza.toml, every POST (bearer, cookie, none) got 403 'actions are off' and pages showed no forms; pages stayed readable. After rotating the token the old cookie got 401 and lost its forms.

## Links

### Out

- verifies [web-token](../capability/rhiza/web-monitor/web-token.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
