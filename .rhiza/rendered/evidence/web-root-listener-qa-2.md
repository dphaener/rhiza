# web-root-listener-qa-2

evidence · `01M6SWYM8S1W507FAA9MACS1H9` · actor [claude-code-770ebe55](../actor/claude-code-770ebe55.md)

## result

PASS

## statement

QA claude-code-770ebe55 at 3056a0e2, hands-on with the built binary (go build -trimpath into the session scratchpad): a scratch root (rhiza root init) with alpha and beta initialised and bare a plain git checkout, each registered with rhiza project add; a hand-written failed run journal in alpha (ended 01:00) and in beta (ended 02:00); rhiza daemon set rhiza.daemon.paused true, rhiza.web.host 127.0.0.1, rhiza.web.port 18947 (written to the root's rhiza.toml [daemon]); rhiza daemon --permission-mode manual from alpha, then curl against http://127.0.0.1:18947. One listener answered for every project: / and /projects 200, /p/alpha/ and /p/beta/ 200, /p/alpha/runs and /p/beta/roadmap 200, /p/alpha 301 to /p/alpha/, /p/nope/ and /p/bare/ 404, /p/alpha/api/status and /p/beta/api/status 200, /static/ and /site.webmanifest 200 at the root. Running rhiza init in bare with the daemon up made /p/bare/ answer 200 with no restart. 90 concurrent requests to /, /inbox and /p/alpha/ during polls all answered and the daemon logged no error; go test -race over internal/web and the api web tests passed.

## Links

### Out

- verifies [web-root-listener](../capability/rhiza/web-root-listener.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
