# web-links-resolve-qa-1

evidence · `01M7BH4AQSMPDS81CFE1PFWGDY` · actor [claude-code-4ec4fe59](../actor/claude-code-4ec4fe59.md)

## result

PASS

## statement

QA at 236ec084e: with the branch's root web handler over a scratch root registering this checkout, a breadth-first crawl from /p/rhiza/, /p/rhiza/roadmap and four node pages followed every href, src and data-rz-fragment for 120 pages: all 200, none left /p/rhiza/. An unscoped crawl of the project handler (213 pages: home, roadmap, runs, change sets, run fragments, node pages) gave no 404; its only failures were the harness's own timeout, and each such URL re-fetched 200. go test -trimpath ./internal/web -run 'TestLinksResolve|TestNodePages' -count=1 ok.

## Links

### Out

- verifies [web-links-resolve](../capability/rhiza/web-monitor/web-links-resolve.md)

## History

- `01M79NN43PGK6HDXTRY33VNC1E` Web model pages: every capability, decision and source link opens a designed page built from the model, with its statement, its why (decisions, sources, change sets) and its history
