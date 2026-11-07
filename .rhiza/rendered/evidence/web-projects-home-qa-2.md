# web-projects-home-qa-2

evidence · `01M6SWYMBAT1MPWVQ172EMY76P` · actor [claude-code-770ebe55](../actor/claude-code-770ebe55.md)

## result

PASS

## statement

QA claude-code-770ebe55 at 3056a0e2, hands-on with the built binary (go build -trimpath into the session scratchpad): a scratch root (rhiza root init) with alpha and beta initialised and bare a plain git checkout, each registered with rhiza project add; a hand-written failed run journal in alpha (ended 01:00) and in beta (ended 02:00); rhiza daemon set rhiza.daemon.paused true, rhiza.web.host 127.0.0.1, rhiza.web.port 18947 (written to the root's rhiza.toml [daemon]); rhiza daemon --permission-mode manual from alpha, then curl against http://127.0.0.1:18947. The web UI now opens on the projects home: GET / (and /projects) listed alpha (none live, 1 in inbox, unhealthy), bare (not initialised) and beta (none live, 1 in inbox, unhealthy), linking /p/alpha/, /p/bare/ and /p/beta/. The old unscoped pages redirect into the home project: /runs 302 to /p/alpha/runs, /roadmap to /p/alpha/roadmap, /runs?x=1 keeps its query. After rhiza init in bare, bare showed none live, 0 in inbox.

## Links

### Out

- verifies [web-projects-home](../capability/rhiza/web-projects-home.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
