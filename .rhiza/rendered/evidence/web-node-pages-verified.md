# web-node-pages-verified

evidence · `01M7BG50TREYR3KCMT49M02QEF` · actor [claude-code-527ff4b9](../actor/claude-code-527ff4b9.md)

## result

PASS

## statement

At e47960bb8 re-ran the Checks block: go test -trimpath ./internal/web -run TestNodePages -count=1 (exit 0), go test -trimpath ./internal/web/... (exit 0), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0, no tracked changes); go vet and gofmt clean. Read the diff: GET /nodes/{ref} (served under /p/<project>/ via scope, nodes mapped to the roadmap nav) renders kind, statement, an 'In your words' quote with where and the answered question, or 'Picked by you' wording with its question, grouped links, 'Where it came from' (sources and change sets) and History; TestNodePages covers a capability, a quoted decision, a worded decision with question, a question, a source and an intent, asserts neither the node's handle nor ID appears, and an unknown ref answers 404 naming the ref. Only internal/web, its templates/CSS and the builder's actor record changed; no TODOs or stubs; tests use the in-process API fixture only.

## Links

### Out

- verifies [web-node-pages](../capability/rhiza/web-monitor/web-node-pages.md)

## History

- `01M79NN43PGK6HDXTRY33VNC1E` Web model pages: every capability, decision and source link opens a designed page built from the model, with its statement, its why (decisions, sources, change sets) and its history
