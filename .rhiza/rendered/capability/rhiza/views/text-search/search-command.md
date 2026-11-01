# search-command

capability · `01M6BDAFV0SK190M10KJ76PMFM` · parent [text-search](../text-search.md)

## acceptance

A search.txtar records three nodes and gets them in the expected ranked order; --kind narrows; --limit 1 gives a next_cursor whose page holds the rest; a bad cursor exits 2; a retired node is not found; a node updated in an unlanded change set is found by its new text only under --view what-if:<cs>; search --json and api call search print the same bytes; describe-golden.txtar is regenerated; the index bumps its schema version and rebuilds

## statement

rhiza search <text> lists live nodes whose handle or field text matches, ranked, with handle, ID, kind, path, matched field and a snippet, narrowed by --kind, read in any view, and paged with --limit and a cursor

## Links

### Out

- decided-by [search-fts-asis-scan-whatif](../../../../decision/search-fts-asis-scan-whatif.md)

### In

- verifies [search-command-verified](../../../../evidence/search-command-verified.md)
- verifies [search-command-qa-1](../../../../evidence/search-command-qa-1.md)
- verifies [pr-170-text-search](../../../../evidence/pr-170-text-search.md)

## History

- `01M6BD9GSN00SK8K209A759RT0` rhiza search finds model nodes by their text, and the model app filters with it
