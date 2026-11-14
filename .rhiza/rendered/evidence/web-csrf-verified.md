# web-csrf-verified

evidence · `01M7BGXC2XN3GJ9G3N39AWBQSR` · actor [claude-code-a86f05aa](../actor/claude-code-a86f05aa.md)

## result

PASS

## statement

At 6d49c35f5 re-ran the Checks block: go test -trimpath ./internal/web -run TestWebCSRF (exit 0, PASS), go test -trimpath ./internal/web (exit 0), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0); also go vet and gofmt -l on internal/web clean, go test ./internal/api -run Web passes. Read the diff: WithToken derives the CSRF value (HMAC of the token) and puts it in the request context via web.CSRF; a cookie-authenticated non-GET/HEAD request without an equal csrf form field gets 403 before reaching the handler, a bearer header needs none. TestWebCSRF drives a test handler behind the gate: GET with the cookie yields the CSRF value (empty without the token, never the token itself), a cookie POST with that value reaches the handler, missing/wrong/empty csrf are refused 403 and the handler is hit only once. Only internal/web plus the builder's actor record and rendered files changed; no stubs or TODOs.

## Links

### Out

- verifies [web-csrf](../capability/rhiza/web-monitor/web-csrf.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
