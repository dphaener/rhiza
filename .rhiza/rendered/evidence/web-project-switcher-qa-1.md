# web-project-switcher-qa-1

evidence · `01M6STWHXG5PEV02TET88XBCCJ` · actor [claude-code-6b4d6ff8](../actor/claude-code-6b4d6ff8.md)

## result

PASS

## statement

QA claude-code-6b4d6ff8 at fac284e5, hands-on with the built binary (go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza): a scratch root (rhiza root init) with alpha and beta initialised and bare a plain git checkout, each registered with rhiza project add; fixture run journals in alpha (a failed run at 01:00, a question at 03:00) and beta (a question at 02:00, a stale run); rhiza daemon set paused=true, web.host 127.0.0.1, web.port 18931; rhiza daemon --permission-mode manual run from alpha, then curl against http://127.0.0.1:18931. /p/alpha/, /p/alpha/roadmap, /p/alpha/runs, /p/beta/runs/<run> and /p/beta/changes/<cs> each showed only their own project's change sets and runs. Each header had nav.rz-switcher with All projects and every initialised project. The links kept the page kind: roadmap went to /p/beta/roadmap, runs to /p/beta/runs, and a run or change set detail page went to the other project's runs or roadmap list. The current project was marked aria-current. Another project's run or change set ID under /p/alpha/ gave 404, and so did /p/nope/events. body data-rz-base=/p/beta. With /p/alpha/events and /p/beta/events open, a new beta journal sent inbox and run events only on beta's stream; alpha's got only its daemon health event.

## Links

### Out

- verifies [web-project-switcher](../capability/rhiza/web-project-switcher.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
