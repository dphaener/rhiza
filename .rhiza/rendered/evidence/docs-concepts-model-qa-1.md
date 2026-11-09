# docs-concepts-model-qa-1

evidence · `01M6ZEA7W608B7HG9NJQWYC48A` · actor [claude-code-ca72ee86](../actor/claude-code-ca72ee86.md)

## result

FAIL

## statement

QA cycle 1, scratch build of the branch head in a fresh scratch repo: ran the model.md examples (add capabilities reports/export-csv, rhiza link export-csv decided-by csv-streaming, rhiza tree, rhiza describe). The ten kinds and the fifteen link types with their endpoints match internal/check/endpoints.go; a wrong-endpoint link (decision verifies capability) is refused; a reference to reports is stored by ID, someone@example.com is left alone, a doubled at-sign becomes a literal one. One claim is wrong: line 64 says handles are unique among siblings, but `rhiza add capability export-csv` under a different parent, and `rhiza add decision reports`, are both refused with 'handle ... is taken' (handles are unique among all live nodes). go test ./internal/doclinks ok.

## Links

### Out

- verifies [docs-concepts-model](../capability/rhiza/documentation/docs-concepts-model.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
