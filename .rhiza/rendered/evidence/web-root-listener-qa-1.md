# web-root-listener-qa-1

evidence · `01M6STWHHKVW4Y9HGA7G2XXD1G` · actor [claude-code-6b4d6ff8](../actor/claude-code-6b4d6ff8.md)

## result

PASS

## statement

QA claude-code-6b4d6ff8 at fac284e5, hands-on with the built binary (go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza): a scratch root (rhiza root init) with alpha and beta initialised and bare a plain git checkout, each registered with rhiza project add; fixture run journals in alpha (a failed run at 01:00, a question at 03:00) and beta (a question at 02:00, a stale run); rhiza daemon set paused=true, web.host 127.0.0.1, web.port 18931; rhiza daemon --permission-mode manual run from alpha, then curl against http://127.0.0.1:18931. One listener answered for every project: / 200 (home alpha), /p/alpha/ and /p/beta/ 200, /p/alpha/runs, /p/beta/roadmap 200, /p/alpha 301, /p/nope/ and /p/bare/ 404. Running rhiza init in bare with the daemon still up made /p/bare/ answer 200 with no restart, and the switcher then listed bare.

## Links

### Out

- verifies [web-root-listener](../capability/rhiza/web-root-listener.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
