# web-merged-inbox-qa-2

evidence · `01M6SWYMDKVS15A38B1T3AYW0J` · actor [claude-code-770ebe55](../actor/claude-code-770ebe55.md)

## result

PASS

## statement

QA claude-code-770ebe55 at 3056a0e2, hands-on with the built binary (go build -trimpath into the session scratchpad): a scratch root (rhiza root init) with alpha and beta initialised and bare a plain git checkout, each registered with rhiza project add; a hand-written failed run journal in alpha (ended 01:00) and in beta (ended 02:00); rhiza daemon set rhiza.daemon.paused true, rhiza.web.host 127.0.0.1, rhiza.web.port 18947 (written to the root's rhiza.toml [daemon]); rhiza daemon --permission-mode manual from alpha, then curl against http://127.0.0.1:18947. GET /inbox showed '2 waiting on you': alpha's failed run (01:00) above beta's (02:00), oldest first, each row tagged alpha or beta and linking /p/alpha/runs/<run> and /p/beta/runs/<run>, both of which answer 200 with that project's run.

## Links

### Out

- verifies [web-merged-inbox](../capability/rhiza/web-merged-inbox.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
