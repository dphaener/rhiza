# web-action-pause

capability · `01M79Q8Q1X2SDBPR5CNADZJAS1` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run 'TestWebPauseFormsCarryCSRFAndAreAbsentWithoutToken|TestWebPausePostsAndRefusesMissingCSRF|TestWebActionsOnly' shows the pause and resume forms carry the CSRF field and are absent when no token is set, pauses and resumes the daemon and a project through POST with the token and the CSRF field, refuses a pause or resume POST missing the CSRF field, and gets 404 or 405 for a POST to any other path, retry, abandon and merge included; go test -trimpath ./internal/api -run TestWebPauseAction shows a POST through WebProjectHandler and WebRootHandler writes rhiza.daemon.paused and the named project's rhiza.project.paused through daemon.set.

## statement

From the web UI the owner pauses or resumes the daemon or one project (rhiza.daemon.paused, rhiza.project.paused), and no other action exists: every other POST path is refused.

## Links

### Out

- depends-on [web-token](web-token.md)
- depends-on [web-csrf](web-csrf.md)
- decided-by [web-safe-actions](../../../decision/web-safe-actions.md)

### In

- verifies [web-action-pause-verified](../../../evidence/web-action-pause-verified.md)
- verifies [web-action-pause-qa-1](../../../evidence/web-action-pause-qa-1.md)
- verifies [web-action-pause-qa-2](../../../evidence/web-action-pause-qa-2.md)
- verifies [pr-574-safe-actions-in-the-web-ui](../../../evidence/pr-574-safe-actions-in-the-web-ui.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
