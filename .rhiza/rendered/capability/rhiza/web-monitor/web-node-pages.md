# web-node-pages

capability · `01M79Q8NR43ZG57MNBDSN6KVJS` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run TestNodePages renders a capability, a decision with a quote, a decision with wording and question, a question, a source and an intent, each with its groups of links and history, and checks that no page shows the node's handle or ID; an unknown ref answers 404 with a page naming the ref.

## statement

Every model node has a page at /p/<project>/nodes/<ref>, as in pages/model-page.html: its kind and statement; for a decision, the owner's quote shown as 'In your words' with the question it answered, or the agent's wording marked as picked by the owner; what it decided or links to, grouped; where it came from (sources and change sets); and its history. Unlike the design's History panel, the page never shows the node's handle or ID (owner-surfaces-no-ids).

## Links

### Out

- decided-by [web-model-pages-designed](../../../decision/web-model-pages-designed.md)

### In

- verifies [web-node-pages-verified](../../../evidence/web-node-pages-verified.md)
- verifies [web-node-pages-qa-1](../../../evidence/web-node-pages-qa-1.md)
- verifies [pr-573-web-model-pages](../../../evidence/pr-573-web-model-pages.md)

## History

- `01M79NN43PGK6HDXTRY33VNC1E` Web model pages: every capability, decision and source link opens a designed page built from the model, with its statement, its why (decisions, sources, change sets) and its history
