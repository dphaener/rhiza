# web-projects-home-qa-1

evidence · `01M6STWHP4PAQVTBT7BP00K6JD` · actor [claude-code-6b4d6ff8](../actor/claude-code-6b4d6ff8.md)

## result

PASS

## statement

QA claude-code-6b4d6ff8 at fac284e5, hands-on with the built binary (go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza): a scratch root (rhiza root init) with alpha and beta initialised and bare a plain git checkout, each registered with rhiza project add; fixture run journals in alpha (a failed run at 01:00, a question at 03:00) and beta (a question at 02:00, a stale run); rhiza daemon set paused=true, web.host 127.0.0.1, web.port 18931; rhiza daemon --permission-mode manual run from alpha, then curl against http://127.0.0.1:18931. GET /projects listed alpha (1 waiting on you · 2 in inbox), bare (not initialised, no counts or health) and beta (2 waiting on you · 1 in inbox). Each project had its own health pill: unknown at first, unhealthy after the first poll, which matches rhiza health in that project (bad check, as the model was not rendered). Rows linked /p/alpha/, /p/bare/ and /p/beta/. The home is served at /projects and linked from every page's nav; / itself is still the home project's overview (see the QA findings).

## Links

### Out

- verifies [web-projects-home](../capability/rhiza/web-projects-home.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
