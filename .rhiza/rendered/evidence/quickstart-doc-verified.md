# quickstart-doc-verified

evidence · `01M6ZAFN5KHZXKJRPKYVPEP7ZY` · actor [claude-code-584deec7](../actor/claude-code-584deec7.md)

## result

PASS

## statement

At 46fd24d1, `go test -trimpath -count=3 ./cmd/rhiza -run 'TestScripts/quickstart$'` (GOTMPDIR outside the repo) passed on all 3 runs. The quickstart-doc command runs every line of docs/quickstart.md's bash blocks in order in a fresh subdirectory repo: init, change new, apply a capability and a decision, propose, accept, land, show, why, slice, tree, render, check, git commit. The tree ends clean. I checked drift in a scratch copy: changing a doc command to `rhiza show no-such-node` failed the script (exit status 4), and replacing one with `ls -la` failed with 'cannot run'. Attempt 2's finding is fixed: the 'Where next' section now links only ../README.md, and `go test -trimpath ./internal/doclinks/...` (TestRepoDocsLinksResolve over the real docs) passes. The prose is plain and natural, and covers every step in the statement.

## Links

### Out

- verifies [quickstart-doc](../capability/rhiza/documentation/quickstart-doc.md)

## History

- `01M6W5EDEX8QAKNVPY5H6CH3Z9` README and quickstart: the README is Rhiza's front door (what and why, install, where next) and docs/quickstart.md walks the model loop without agents, run by a testscript
