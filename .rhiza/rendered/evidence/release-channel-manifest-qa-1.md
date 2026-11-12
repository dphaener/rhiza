# release-channel-manifest-qa-1

evidence · `01M77GAJVJ9F5ZG9V9JKAG75S6` · actor [claude-code-f1ace9df](../actor/claude-code-f1ace9df.md)

## result

PASS

## statement

QA at 55bf2f1d0: cross-compiled the real rhiza binary for linux/darwin x amd64/arm64 with the .goreleaser.yaml ldflags, packed rhiza_0.1.0_<os>_<arch>.tar.gz per its name_template, wrote checksums.txt with sha256sum, then ran release.yml's step `go run ./tools/release -dist dist -tag v0.1.0 -commit <sha> -repo dphaener/rhiza`: exit 0, dist/manifest.json has channel=release, version=v0.1.0, the full commit, and 4 platforms sorted by os/arch, each URL https://github.com/dphaener/rhiza/releases/download/v0.1.0/<archive> and sha256 equal to checksums.txt; the extracted linux_amd64 binary prints `rhiza v0.1.0`. Re-ran the release-snapshot job's new step on a GoReleaser-shaped snapshot dist (metadata.json version 0.0.1-SNAPSHOT-55bf2f1): exit 0 and the jq length check = 4; with darwin_arm64 removed it fails ('missing platforms: darwin/arm64'), with a tampered archive it fails (sha256 mismatch), without -tag it fails. go test -count=1 -trimpath ./internal/release/... ./tools/release/... passes; gofmt and vet clean. GoReleaser itself not run (not installed here; PR CI runs it).

## Links

### Out

- verifies [release-channel-manifest](../capability/rhiza/distribution/release-channel-manifest.md)

## History

- `01M76AZ94M8HWY6GCD06H8ES46` Release channel manifest and provenance: every v* release publishes a JSON manifest (channel, version, commit, per-platform archive URL and sha256) and GitHub build-provenance attestations alongside its archives and checksums
