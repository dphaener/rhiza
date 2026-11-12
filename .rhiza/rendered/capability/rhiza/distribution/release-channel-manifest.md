# release-channel-manifest

capability · `01M76B7CPKY732RZT6D4H4N8ND` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./internal/release/... ./tools/release/... builds a manifest from a dist folder with the four archives and checksums.txt and checks every field and that the sha256s match checksums.txt; the release-snapshot job writes the manifest and fails unless all four platforms are in it

## statement

Every v* release publishes manifest.json beside its archives: channel (release), version, commit, and for each os/arch the archive URL and sha256, so releases/latest/download/manifest.json always names the newest release

## Links

### In

- verifies [release-channel-manifest-verified](../../../evidence/release-channel-manifest-verified.md)
- verifies [release-channel-manifest-qa-1](../../../evidence/release-channel-manifest-qa-1.md)
- verifies [pr-476-release-channel-manifest-and](../../../evidence/pr-476-release-channel-manifest-and.md)

## History

- `01M76AZ94M8HWY6GCD06H8ES46` Release channel manifest and provenance: every v* release publishes a JSON manifest (channel, version, commit, per-platform archive URL and sha256) and GitHub build-provenance attestations alongside its archives and checksums
