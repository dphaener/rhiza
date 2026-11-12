# latest-build-version

capability · `01M76B7CWZHSPFD2WJYGWGZ0HH` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/release/... computes the version with no tag, with a tag N commits back and with a tag at HEAD, and a semver comparison orders it after the tag, before the next patch, and after the previous latest build

## statement

A latest-channel build reports the next patch of the last v* tag as a semver prerelease: vX.Y.(Z+1)-main.<commits since the tag>+<short sha> (v0.0.1-main.N+sha before any tag), which sorts after the tag and before the next patch

## Links

### In

- verifies [latest-build-version-verified](../../../evidence/latest-build-version-verified.md)
- verifies [latest-build-version-qa-1](../../../evidence/latest-build-version-qa-1.md)
- verifies [pr-478-latest-channel](../../../evidence/pr-478-latest-channel.md)

## History

- `01M76AZHGY605TXF68M27W4SZV` Latest channel: every merge to main builds the four archives as the next patch semver prerelease (vX.Y.Z+1-main.N+sha) and replaces one rolling latest GitHub prerelease with them, their checksums, manifest and attestations
