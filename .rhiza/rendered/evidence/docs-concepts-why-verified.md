# docs-concepts-why-verified

evidence · `01M6ZDSMEKG72NH3RTV0CBD9PW` · actor [claude-code-1d1eb393](../actor/claude-code-1d1eb393.md)

## result

PASS

## statement

At e3cb612c, docs/concepts/why.md covers each point the leaf asks for: decisions and verbatim owner quotes (decided-by, supersedes, accept/reject --reason/--rule), questions (blocks, rhiza question answer), sources and scrubbing (source add with --raw/--raw-sha256/--raw-session, hash-only, derived-from), evidence (verifies, observes) and reading back with rhiza why and rhiza slice. I checked every command, flag and link type it names against `rhiza <cmd> --help`, and they all match. docs-links-resolve: go test -trimpath -count=1 ./internal/doclinks/... passes, including TestRepoDocsLinksResolve over the real docs. docs-commands-exist is not landed in this tree, so I built its commit ffea5bfb into the scratchpad and ran its rhiza check here. It reports no docs-command failure for why.md (the only one is README's `rhiza version`, which is outside this leaf). As a negative control I added `rhiza no-such-command` to a copy of why.md, and that check flagged it. The prose reads plain and human, with no AI tells I could find. The owner's read is left to the owner.

## Links

### Out

- verifies [docs-concepts-why](../capability/rhiza/documentation/docs-concepts-why.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
