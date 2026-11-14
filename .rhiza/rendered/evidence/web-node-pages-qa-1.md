# web-node-pages-qa-1

evidence · `01M7BH4AHXGZQC7WKN05533QAK` · actor [claude-code-4ec4fe59](../actor/claude-code-4ec4fe59.md)

## result

PASS

## statement

QA at 236ec084e: built the branch with go build -trimpath and served the real api.WebProjectHandler and api.WebRootHandler (a scratch root registering this checkout) on 127.0.0.1 over this repo's own model, with no daemon. GET /nodes/<ref> and /p/rhiza/nodes/<ref>, by handle and by full ID, gave 200 for a capability (web-node-pages, the root rhiza), a constraint, a decision with an owner quote (cli-rules: 'In your words' with where), a decision with wording and question (web-model-pages-designed: 'Picked by you', 'Answering: ...'), a question and a source (title plus summary). Each showed kind, statement, grouped links (Decided by, What it decided, Constrains, Evidence), Where it came from (sources, change sets) and History, with no handle or ID in the page text (IDs only inside hrefs). Unknown refs (no-such-thing, a change set ID, a short ID, an upper-case handle, a%2Fb, %20, an escaped script tag) answered 404 with 'Nothing found for <ref>', HTML-escaped. Unscoped /nodes/x redirected 302 to /p/rhiza/nodes/x, Roadmap nav was marked current and the 404 crumb was scoped.

## Links

### Out

- verifies [web-node-pages](../capability/rhiza/web-monitor/web-node-pages.md)

## History

- `01M79NN43PGK6HDXTRY33VNC1E` Web model pages: every capability, decision and source link opens a designed page built from the model, with its statement, its why (decisions, sources, change sets) and its history
