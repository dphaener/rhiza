# docs-concepts-model-qa-2

evidence · `01M6ZF5T1NG8XKTTYJENX4FVQW` · actor [claude-code-1b418df2](../actor/claude-code-1b418df2.md)

## result

PASS

## statement

QA cycle 2, scratch build at f10165fc in a fresh scratch repo: checked the fixed handles paragraph. add capability export-csv under a second parent and add decision reports (a capability's handle) are both refused with 'handle ... is taken by node ...'. After retire tmp-h --reason gone, add capability tmp-h succeeds, so a retired node frees its handle. Matches docs/concepts/model.md as written. go test -trimpath ./internal/doclinks/... ok.

## Links

### Out

- verifies [docs-concepts-model](../capability/rhiza/documentation/docs-concepts-model.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
