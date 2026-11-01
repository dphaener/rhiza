# search-command-qa-1

evidence · `01M6BJT43Y80R3TJTJWS61EM9D` · actor [claude-code-85cc7f99](../actor/claude-code-85cc7f99.md)

## result

PASS

## statement

At 7a46ec69 with a scratch build ($TMPDIR/rhiza-qa1) on a clone of this repo's real model: rhiza search search ranks the handle hit search-fts-asis-scan-whatif (score 100) over field hits (score 20) with snippets; --kind capability narrows; --limit gives a next: line with --cursor and --json a next_cursor; '*' exits 2 with 'nothing to search for'; FTS syntax like 'AND OR NOT' and 'fts5"' is taken as plain words, no error; what-if:01M6BD9G finds this change set's unlanded text-search, search-command and model-app-search nodes while as-is does not. An index built at user_version 1 by an origin/main build was rebuilt to user_version 2 with 5623 nodes_fts rows on the first search. In a scratch repo, adding beta, then setting alpha's statement, showed in the next search each time (gizmo then plain, prefix GIZ case-insensitive).

## Links

### Out

- verifies [search-command](../capability/rhiza/views/text-search/search-command.md)

## History

- `01M6BD9GSN00SK8K209A759RT0` rhiza search finds model nodes by their text, and the model app filters with it
