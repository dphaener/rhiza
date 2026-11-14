# web-links-resolve

capability · `01M79Q8NS7D874K7TVBX23DXGA` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run TestLinksResolve crawls Home, the roadmap, a change set page, a run page and a node page of a two-project fixture root and gets 200 for every internal link.

## statement

Every link the web UI renders inside a project resolves to a page that exists in that project.

## Links

### Out

- decided-by [web-model-pages-designed](../../../decision/web-model-pages-designed.md)

### In

- verifies [web-links-resolve-verified](../../../evidence/web-links-resolve-verified.md)
- verifies [web-links-resolve-qa-1](../../../evidence/web-links-resolve-qa-1.md)
- verifies [pr-573-web-model-pages](../../../evidence/pr-573-web-model-pages.md)

## History

- `01M79NN43PGK6HDXTRY33VNC1E` Web model pages: every capability, decision and source link opens a designed page built from the model, with its statement, its why (decisions, sources, change sets) and its history
