# m30-goreleaser

decision · `01M6W5BQP21T6REHE624ZGDMHK`

## question

How are releases built? (version stamping via -ldflags already exists in the Makefile)

## statement

Build releases with GoReleaser from a GitHub Actions job on v* tags

## wording

GoReleaser on tag: a .goreleaser.yaml and a GitHub Actions job on v* tags building linux/darwin amd64+arm64 with checksums, attached to a GitHub release, picked by the owner

## Links

### Out

- derived-from [m30-planning](../source/m30-planning.md)

### In

- decided-by [m30-install](../intent/milestone-30/m30-install.md)
- decided-by [release-goreleaser](../capability/rhiza/distribution/release-goreleaser.md)

## History

- `01M6W5A2R4BXDGQGHNWGZNB3XY` Milestone 30: Rhiza is ready to share: install, README and docs, generated reference, public under MIT
