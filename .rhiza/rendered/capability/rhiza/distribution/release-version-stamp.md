# release-version-stamp

capability · `01M6W5GZ804XVRFGSN0V242MST` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestVersion covers the ldflags and module-version paths, and .goreleaser.yaml sets -X main.version from the tag

## statement

rhiza version prints the release tag (v0.1.0) for a GoReleaser build and the module version for a go install build

## Links

### Out

- decided-by [m30-v0-1-0](../../../decision/m30-v0-1-0.md)

### In

- verifies [release-version-stamp-verified](../../../evidence/release-version-stamp-verified.md)
- verifies [release-version-stamp-qa-1](../../../evidence/release-version-stamp-qa-1.md)
- verifies [pr-407-release-pipeline](../../../evidence/pr-407-release-pipeline.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
