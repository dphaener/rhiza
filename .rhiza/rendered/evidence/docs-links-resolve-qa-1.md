# docs-links-resolve-qa-1

evidence · `01M6ZATPM1X5Y61C44AMBX27GY` · actor [claude-code-d4acd8b3](../actor/claude-code-d4acd8b3.md)

## result

PASS

## statement

QA at b92987bf: go test -trimpath -count=1 ./internal/doclinks/... passes over the real README.md and docs/. With README's quickstart link changed to docs/quickstart-nope.md it fails with README.md:23 docs/quickstart-nope.md does not exist; with the quickstart's ../README.md#install changed to #installing it fails with docs/quickstart.md:3 no heading #installing in README.md. Both edits reverted afterwards.

## Links

### Out

- verifies [docs-links-resolve](../capability/rhiza/documentation/docs-links-resolve.md)

## History

- `01M6W5EDEX8QAKNVPY5H6CH3Z9` README and quickstart: the README is Rhiza's front door (what and why, install, where next) and docs/quickstart.md walks the model loop without agents, run by a testscript
