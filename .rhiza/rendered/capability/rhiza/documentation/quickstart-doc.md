# quickstart-doc

capability · `01M6W5H11AZV166CP6M6R0F9SX` · parent [documentation](../documentation.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/quickstart runs every command in docs/quickstart.md's shell blocks, in order, in a fresh repo, and fails if any fails or the doc and script drift. Prose written with the humanize skill (constraint doc-prose-humanized).

## statement

docs/quickstart.md walks a new user through the model loop without agents: init in a fresh repo, add a capability and a decision, open, propose, accept and land a change set, render, and read slice and why

## Links

### Out

- decided-by [m30-quickstart-model-loop](../../../decision/m30-quickstart-model-loop.md)
- decided-by [m30-humanize-rule](../../../decision/m30-humanize-rule.md)

### In

- verifies [quickstart-doc-verified](../../../evidence/quickstart-doc-verified.md)
- verifies [quickstart-doc-qa-1](../../../evidence/quickstart-doc-qa-1.md)
- verifies [pr-408-readme-and-quickstart](../../../evidence/pr-408-readme-and-quickstart.md)

## History

- `01M6W5EDEX8QAKNVPY5H6CH3Z9` README and quickstart: the README is Rhiza's front door (what and why, install, where next) and docs/quickstart.md walks the model loop without agents, run by a testscript
