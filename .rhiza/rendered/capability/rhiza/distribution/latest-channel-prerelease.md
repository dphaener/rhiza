# latest-channel-prerelease

capability · `01M76B7CXJ4H445ZBV963ADV16` · parent [distribution](../distribution.md)

## acceptance

The workflow triggers on push to main only, runs in one concurrency group that cancels older runs, and skips publishing when the latest tag already points at a descendant of its commit; the release-snapshot job builds a latest-channel snapshot with the computed version and requires the four archives and the manifest

## statement

Every push to main runs a Latest workflow that builds the four archives with that version and replaces the rolling latest GitHub prerelease (its tag moved to the commit) with them, checksums.txt, manifest.json (channel latest) and provenance attestations; it is never GitHub latest release, and an older run never overwrites a newer commit

## Links

### Out

- decided-by [m33-release-requirements](../../../decision/m33-release-requirements.md)

### In

- verifies [latest-channel-prerelease-verified](../../../evidence/latest-channel-prerelease-verified.md)
- verifies [latest-channel-prerelease-qa-1](../../../evidence/latest-channel-prerelease-qa-1.md)
- verifies [latest-channel-prerelease-verified-2](../../../evidence/latest-channel-prerelease-verified-2.md)
- verifies [pr-478-latest-channel](../../../evidence/pr-478-latest-channel.md)

## History

- `01M76AZHGY605TXF68M27W4SZV` Latest channel: every merge to main builds the four archives as the next patch semver prerelease (vX.Y.Z+1-main.N+sha) and replaces one rolling latest GitHub prerelease with them, their checksums, manifest and attestations
