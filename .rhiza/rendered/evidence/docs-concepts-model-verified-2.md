# docs-concepts-model-verified-2

evidence · `01M6ZERDVR28WQVZDRKYYH6096` · actor [claude-code-d0b7a712](../actor/claude-code-d0b7a712.md)

## result

PASS

## statement

At 654ccf2a I read docs/concepts/model.md and the fix commit's diff. The page covers the ten node kinds, the fifteen link types and their endpoints, handles, IDs, ID prefixes and @ text references, and the capability tree. The handle paragraph now says handles are unique among live nodes repo-wide, whatever their kind or parent, and that retiring a node frees its handle. That matches model.HandleTaken (internal/api/model.go:482), which matches any non-retired node with the handle, and the `rhiza add --help` text ('unique among live nodes'). This settles the confirmed quality finding f1. go test -trimpath -count=1 ./internal/doclinks/... passes, including TestRepoDocsLinksResolve over the real README and docs/ (docs-links-resolve). docs-commands-exist is not built in this tree yet, so I checked it by hand: every rhiza command the page names (describe, add, link, tree, slice) exists, since each one's --help succeeded on a scratch build of this commit. The prose reads in the owner's plain style.

## Links

### Out

- verifies [docs-concepts-model](../capability/rhiza/documentation/docs-concepts-model.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
