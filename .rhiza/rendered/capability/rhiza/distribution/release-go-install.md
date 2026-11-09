# release-go-install

capability · `01M6W5GZC2YS53M8KZWFQK4E00` · parent [distribution](../distribution.md)

## acceptance

a CI step runs go install ./cmd/rhiza into a fresh GOBIN with -mod=readonly and runs rhiza version; go.mod has no replace directives

## statement

go install github.com/dphaener/rhiza/cmd/rhiza@<version> builds a working rhiza: go.mod's module path matches the repository and nothing needs a replace directive or a private dependency

## Links

### Out

- decided-by [m30-install-path](../../../decision/m30-install-path.md)

### In

- verifies [release-go-install-verified](../../../evidence/release-go-install-verified.md)
- verifies [release-go-install-qa-1](../../../evidence/release-go-install-qa-1.md)
- verifies [pr-407-release-pipeline](../../../evidence/pr-407-release-pipeline.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
