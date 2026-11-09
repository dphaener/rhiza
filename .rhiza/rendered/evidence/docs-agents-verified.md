# docs-agents-verified

evidence · `01M6ZHEB0MWJJFM8YCWDPCDEEZ` · actor [claude-code-beae8c26](../actor/claude-code-beae8c26.md)

## result

PASS

## statement

At 06da5bad, read docs/agents.md and the diff: it covers installing the guide and skills (rhiza guide, guide ls, guide install --agent claude, --project, --force), the read-slice/work/apply/commit/land loop, and the Change:/Actor:/Leaf: trailers with Co-Authored-By in one paragraph; the rules section matches `rhiza guide`. docs-links-resolve: go test -trimpath -count=1 ./internal/doclinks/... passes over the real README and docs/ (new links to quickstart.md, concepts/change-sets.md, concepts/why.md, concepts/model.md, agents.md resolve). docs-commands-exist has no automated check in the tree yet, so checked by hand with a scratch build: every rhiza command the page names (guide, guide ls, guide build, guide install, slice, ready, apply, land, render, check, accept) answers --help, and guide install's flags match the page. TestScripts/quickstart still passes after the quickstart.md edit. Prose reads plain and unpadded, consistent with doc-prose-humanized; the owner's read is left to the owner.

## Links

### Out

- verifies [docs-agents](../capability/rhiza/documentation/docs-agents.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
