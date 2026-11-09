# release-goreleaser

capability · `01M6W5GZ1GA5NPS1AP0QYYSKTA` · parent [distribution](../distribution.md)

## acceptance

goreleaser check passes on .goreleaser.yaml; a CI step runs goreleaser release --snapshot --clean on pull requests that touch it and produces all four archives; the release job's trigger is only v* tags

## statement

A v* tag runs a GitHub Actions release job that builds rhiza with GoReleaser for linux and darwin on amd64 and arm64 and publishes the archives and a checksums file to a GitHub release

## Links

### Out

- decided-by [m30-goreleaser](../../../decision/m30-goreleaser.md)
- decided-by [m30-install-path](../../../decision/m30-install-path.md)

### In

- verifies [release-goreleaser-verified](../../../evidence/release-goreleaser-verified.md)
- verifies [release-goreleaser-qa-1](../../../evidence/release-goreleaser-qa-1.md)
- verifies [pr-407-release-pipeline](../../../evidence/pr-407-release-pipeline.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
