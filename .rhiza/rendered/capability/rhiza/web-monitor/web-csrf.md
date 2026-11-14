# web-csrf

capability · `01M7B49JW37Q5ECBXE03TZF2EQ` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run TestWebCSRF shows, through a test handler behind the gate, that a cookie-authenticated POST with the right CSRF field passes the gate and one with a missing or wrong CSRF field is refused, and that the CSRF value the gate hands pages for their forms is the one it accepts.

## statement

Every web UI POST authenticated by the token cookie must carry a CSRF field that the page's action forms include; a cookie POST with a missing or wrong CSRF field is refused, whatever action it names.

## Links

### Out

- depends-on [web-token](web-token.md)
- decided-by [web-safe-actions](../../../decision/web-safe-actions.md)

### In

- depends-on [web-action-answer](web-action-answer.md)
- depends-on [web-action-pause](web-action-pause.md)
- verifies [web-csrf-verified](../../../evidence/web-csrf-verified.md)
- verifies [web-csrf-qa-1](../../../evidence/web-csrf-qa-1.md)
- verifies [web-csrf-qa-2](../../../evidence/web-csrf-qa-2.md)
- verifies [pr-574-safe-actions-in-the-web-ui](../../../evidence/pr-574-safe-actions-in-the-web-ui.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
