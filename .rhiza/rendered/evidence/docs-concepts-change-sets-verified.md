# docs-concepts-change-sets-verified

evidence · `01M6ZDYQ6ST4R5R0XZQM8ZCT08` · actor [claude-code-72ed2789](../actor/claude-code-72ed2789.md)

## result

PASS

## statement

At a33575de (attempt 2), docs/concepts/change-sets.md covers intention and record, the draft-to-landed lifecycle with reject/withdraw, the append-only log, carving, reopening and rendering. Attempt 1's finding is fixed: the carve example now reads `rhiza carve access-roles --from <source-id> --to <target-id> --with-children`, which matches `rhiza carve --help` (--from/--to take an ID or a 4+ character prefix). docs-links-resolve: `go test -trimpath ./internal/doclinks/...` ok. docs-commands-exist: the check is not in the tree at this commit, so I checked by hand. With a scratch build, `--help` succeeds for every command the doc names: change new/use/show/ls, propose, accept, land, reject, withdraw, link, log, carve, reopen, render, check and ready. --change is a global flag. The discovered-from link type exists. The prose reads plainly, with no AI tells. The owner's read is still to come.

## Links

### Out

- verifies [docs-concepts-change-sets](../capability/rhiza/documentation/docs-concepts-change-sets.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
