# search-command-verified

evidence · `01M6BJ8EB404GS2XQ8GTMAHF8F` · actor [claude-code-fdd64d60](../actor/claude-code-fdd64d60.md)

## result

PASS

## statement

At 5e82132d, go test ./cmd/rhiza -run 'TestScripts/(search|describe-golden)$' passes: search.txtar records nodes and gets disk, store, cache in ranked order; --kind narrows; --limit 1 prints a next_cursor whose page holds store and cache; a bad cursor exits 2; a retired node is not found; a node set in an unlanded change set is found by its new text under what-if and not under as-is; search --json and api call search cmp equal; describe-golden.txtar includes search. By hand with a scratch build, --view as-is and to-be miss the new text and find the old, while what-if finds the new. SchemaVersion went from 1 to 2 with an nodes_fts table, and an index downgraded to user_version 1 with nodes_fts dropped was rebuilt to version 2 on the next search. make check exits 0.

## Links

### Out

- verifies [search-command](../capability/rhiza/views/text-search/search-command.md)

## History

- `01M6BD9GSN00SK8K209A759RT0` rhiza search finds model nodes by their text, and the model app filters with it
