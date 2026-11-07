# web-project-switcher-qa-2

evidence · `01M6SWYMGAD8WVAFV5VWMX24NK` · actor [claude-code-770ebe55](../actor/claude-code-770ebe55.md)

## result

PASS

## statement

QA claude-code-770ebe55 at 3056a0e2, hands-on with the built binary (go build -trimpath into the session scratchpad): a scratch root (rhiza root init) with alpha and beta initialised and bare a plain git checkout, each registered with rhiza project add; a hand-written failed run journal in alpha (ended 01:00) and in beta (ended 02:00); rhiza daemon set rhiza.daemon.paused true, rhiza.web.host 127.0.0.1, rhiza.web.port 18947 (written to the root's rhiza.toml [daemon]); rhiza daemon --permission-mode manual from alpha, then curl against http://127.0.0.1:18947. Every page kind under /p/alpha/ and /p/beta/ (overview, runs, roadmap, run detail, change set) carried data-rz-base=/p/<project> and a switcher of All projects, alpha and beta with its own project aria-current, the other project's link keeping the page kind (runs to runs, roadmap to roadmap, run detail and change set to that project's runs and roadmap). Data stayed scoped: alpha's run is absent from /p/beta/runs and /p/beta/runs/<alpha run> is 404; alpha's change set is 200 under /p/alpha/changes/ and 404 under /p/beta/changes/. /p/alpha/events answered 200 text/event-stream, /p/nope/events 404. The switcher listed bare as soon as it was initialised.

## Links

### Out

- verifies [web-project-switcher](../capability/rhiza/web-project-switcher.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
