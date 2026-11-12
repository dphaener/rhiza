# latest-channel-prerelease-verified-2

evidence · `01M77JR7DM08SD9GK6HKX7JTBS` · actor [claude-code-bd156e4e](../actor/claude-code-bd156e4e.md)

## result

PASS

## statement

Verified at d7d285186 (attempt 2, verifier). Read .github/workflows/latest.yml: on push to branches [main] only; concurrency group 'latest' with cancel-in-progress true; a gate step fetches the latest tag and sets publish=false when latest^{commit} differs from GITHUB_SHA and GITHUB_SHA is its ancestor (latest points at a descendant), and every later step is conditioned on publish=='true'; it replaces the prerelease with --prerelease --latest=false --target GITHUB_SHA, uploading the archives, checksums.txt and manifest.json (channel latest), after attest-build-provenance. Read .github/workflows/release-snapshot.yml: computes RHIZA_BUILD_VERSION with go run ./tools/release version, runs goreleaser check and release --snapshot --clean (snapshot.version_template trims the v from RHIZA_BUILD_VERSION), requires all four linux/darwin amd64/arm64 archives, writes the manifest with -channel latest -release-tag latest and asserts 4 platforms. Read the HEAD diff: archive name_template maps '+' to '-' and release.Build applies AssetVersion the same way. go test -trimpath ./internal/release/... ./tools/release/... passed. Simulated a latest dist (four archives named rhiza_0.0.1-main.5-abc1234_<os>_<arch>.tar.gz plus checksums.txt) and ran go run ./tools/release -tag v0.0.1-main.5+abc1234 -version 0.0.1-main.5+abc1234 -channel latest -release-tag latest: the manifest has version v0.0.1-main.5+abc1234, channel latest, 4 platforms, URLs under releases/download/latest/ naming the '-' archives. goreleaser is not installed here, so the real snapshot build was not run locally; the PR's Release snapshot job runs it.

## Links

### Out

- verifies [latest-channel-prerelease](../capability/rhiza/distribution/latest-channel-prerelease.md)

## History

- `01M76AZHGY605TXF68M27W4SZV` Latest channel: every merge to main builds the four archives as the next patch semver prerelease (vX.Y.Z+1-main.N+sha) and replaces one rolling latest GitHub prerelease with them, their checksums, manifest and attestations
