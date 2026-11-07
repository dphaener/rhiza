# web-merged-inbox-qa-1

evidence · `01M6STWHTTPZ17YA6Z0NKC7RRA` · actor [claude-code-6b4d6ff8](../actor/claude-code-6b4d6ff8.md)

## result

PASS

## statement

QA claude-code-6b4d6ff8 at fac284e5, hands-on with the built binary (go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza): a scratch root (rhiza root init) with alpha and beta initialised and bare a plain git checkout, each registered with rhiza project add; fixture run journals in alpha (a failed run at 01:00, a question at 03:00) and beta (a question at 02:00, a stale run); rhiza daemon set paused=true, web.host 127.0.0.1, web.port 18931; rhiza daemon --permission-mode manual run from alpha, then curl against http://127.0.0.1:18931. GET /inbox listed the items oldest first across projects: alpha Alpha failing work (failed, 01:00), beta Which beta colour? (02:00), alpha Which alpha store? (03:00). Each row had a rz-project-tag span with its project and linked /p/<project>/runs/<run>, and each link answered 200. A question journal added to beta while the daemon ran showed up on /inbox, tagged beta.

## Links

### Out

- verifies [web-merged-inbox](../capability/rhiza/web-merged-inbox.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
