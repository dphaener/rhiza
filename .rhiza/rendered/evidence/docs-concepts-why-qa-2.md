# docs-concepts-why-qa-2

evidence · `01M6ZF5T3J9GNCVVZ6GVE04H7D` · actor [claude-code-1b418df2](../actor/claude-code-1b418df2.md)

## result

PASS

## statement

QA cycle 2, scratch build at f10165fc in a fresh scratch repo: ran the Questions example from docs/concepts/why.md verbatim (add question csv-delimiter --record, link csv-delimiter blocks export-csv --record, change new, set export-csv, question answer csv-delimiter). The answer was written as decision answer-csv-delimiter with an answers link. The stated rule also holds: an unlanded question in the current change set is refused, a landed question blocking an untouched node is refused, and landed questions blocking an ancestor of a touched node or the change set itself are answered. why and slice export-csv print as described.

## Links

### Out

- verifies [docs-concepts-why](../capability/rhiza/documentation/docs-concepts-why.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
