# web-token-verified

evidence · `01M7BGJA6AKAB254F3SA6SG7CG` · actor [claude-code-ab4b9c05](../actor/claude-code-ab4b9c05.md)

## result

PASS

## statement

At 7f81e4e55 re-ran the Checks block: go test -trimpath ./internal/web -run TestWebToken (5 tests PASS: no token set refuses every POST with 403, wrong/missing/basic token refused 401, ?token= visit sets an HttpOnly SameSite=Strict cookie and redirects without the token, bearer accepted, ActionsOn), go test -trimpath ./internal/api -run TestWebTokenSetting PASS (daemon.set accepts rhiza.web.token and rejects a short one, replies 'set'; WebProjectHandler from web_root.go refuses POST 403 with none set, 401 with a wrong bearer, passes the right one, read per request; daemon.settings, daemon state rows and unauthenticated GET /api/daemon.settings and /api/daemon.status never carry the value), the wider api run, describe-golden -update (no diff), make self-check (exit 0 on two consecutive runs; one earlier run in the same chained command exited 1 and did not reproduce), rhiza ship generate and rhiza render (exit 0, no diff). The diff masks the value at its source (settingsOf), which covers daemon.status as the answered review question asked; it stays within the leaf and has no stubs or TODOs.

## Links

### Out

- verifies [web-token](../capability/rhiza/web-monitor/web-token.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
