# m30-generated-reference

decision · `01M6W5BQB2SZW5TD310XD49GJJ`

## question

How is the CLI and config reference kept accurate?

## statement

Generate the reference from the code and fail make check when it is stale

## wording

Generated, CI-checked: a rhiza docs command renders docs/reference/*.md from the op registry and settings schema; make check fails when it's stale, picked by the owner

## Links

### Out

- derived-from [m30-planning](../source/m30-planning.md)

### In

- decided-by [m30-reference](../intent/milestone-30/m30-reference.md)
- decided-by [docs-reference-cli](../capability/rhiza/documentation/docs-reference-cli.md)
- decided-by [docs-reference-config](../capability/rhiza/documentation/docs-reference-config.md)
- decided-by [docs-reference-workflow](../capability/rhiza/documentation/docs-reference-workflow.md)
- decided-by [docs-reference-stale](../capability/rhiza/documentation/docs-reference-stale.md)
- decided-by [docs-commands-exist](../capability/rhiza/documentation/docs-commands-exist.md)

## History

- `01M6W5A2R4BXDGQGHNWGZNB3XY` Milestone 30: Rhiza is ready to share: install, README and docs, generated reference, public under MIT
