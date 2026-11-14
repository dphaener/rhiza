# web-token

capability · `01M79Q8PX68N5R6XWD5MHV6FKJ` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run 'TestWebTokenNoneSetRefusesEveryPost|TestWebTokenWrongOrMissingRefused|TestWebTokenCookieAfterQueryVisit|TestWebTokenBearerAccepted|TestWebTokenActionsOn' shows every POST refused when no token is set, a POST refused with a wrong or missing token, the cookie set after a ?token= visit, and a bearer header accepted; go test -trimpath ./internal/api -run TestWebTokenSetting shows daemon.set accepts rhiza.web.token, the daemon's web handler (web_root.go) gates POSTs on the token it reads from that setting, and neither daemon.settings nor an unauthenticated GET /api/daemon.settings through WebProjectHandler returns the token's value.

## statement

A web token (rhiza.web.token, set with rhiza daemon set and read by the daemon's web handler) turns on the web UI's actions: a browser presents it once with ?token= and keeps it in an HttpOnly cookie, a script sends it as a bearer header; the token's value is never shown, so daemon.settings and the web UI's read-only API (GET /api/daemon.settings) report only whether it is set; without a token set, pages show no action controls and every POST is refused, and pages stay readable as today.

## Links

### Out

- decided-by [web-safe-actions](../../../decision/web-safe-actions.md)

### In

- depends-on [web-csrf](web-csrf.md)
- depends-on [web-action-answer](web-action-answer.md)
- depends-on [web-action-pause](web-action-pause.md)
- blocks [review-01m7b3wjc309vgwck1ef6ms5gb-r3-f1](../../../question/review-01m7b3wjc309vgwck1ef6ms5gb-r3-f1.md)
- verifies [web-token-verified](../../../evidence/web-token-verified.md)
- verifies [web-token-qa-1](../../../evidence/web-token-qa-1.md)
- verifies [web-token-qa-2](../../../evidence/web-token-qa-2.md)
- verifies [pr-574-safe-actions-in-the-web-ui](../../../evidence/pr-574-safe-actions-in-the-web-ui.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
