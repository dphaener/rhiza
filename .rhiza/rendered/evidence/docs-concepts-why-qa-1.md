# docs-concepts-why-qa-1

evidence · `01M6ZEA82PGMB4N1RZ4B532039` · actor [claude-code-ca72ee86](../actor/claude-code-ca72ee86.md)

## result

FAIL

## statement

QA cycle 1, scratch build in a fresh scratch repo: ran every why.md example in order in one change set. apply decision.json, add question, link blocks, source add --raw-sha256, add evidence, link verifies, rhiza why (decisions, sources, change sets shown), rhiza slice and --view as-is (refused for an unlanded node, as stated) all work. But `rhiza question answer csv-delimiter ...` fails as the page sequences it ('not an open question blocking this change set'), even with the question also linked blocks to the change set. It succeeds only for a question that is already landed and blocks a node the current change set touches. accept/reject take --reason and --rule, and the source add raw flags exist.

## Links

### Out

- verifies [docs-concepts-why](../capability/rhiza/documentation/docs-concepts-why.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
