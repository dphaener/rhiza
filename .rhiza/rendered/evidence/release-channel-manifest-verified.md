# release-channel-manifest-verified

evidence · `01M77FQNTM4W0GTV29B3B401TM` · actor [claude-code-53e046f1](../actor/claude-code-53e046f1.md)

## result

PASS

## statement

At 376a05b23, go test -count=1 -trimpath ./internal/release/... ./tools/release/... passes (vet and gofmt clean). TestBuild builds a manifest from a dist folder with the four rhiza_<ver>_<os>_<arch>.tar.gz archives and checksums.txt and checks channel=release, version=tag, commit, 4 platforms, each sha256 equal to checksums.txt and the releases/download/<tag>/<archive> URL; TestBuildMissingPlatform and TestBuildChecksumMismatch fail as required. Ran the release-snapshot job's new step by hand on a fake GoReleaser-shaped dist (metadata.json version 0.0.1-SNAPSHOT-376a05b): it wrote dist/manifest.json with 4 platforms and the jq length check passed; with darwin_arm64 removed, go run ./tools/release exited 1 with 'missing platforms: darwin/arm64', so the job fails. release.yml runs the tool after goreleaser on v* tags and uploads manifest.json as a release asset (--clobber), so releases/latest/download/manifest.json names the newest release. Diff read: archive names match .goreleaser.yaml name_template. GoReleaser itself not run (not installed; CI runs it).

## Links

### Out

- verifies [release-channel-manifest](../capability/rhiza/distribution/release-channel-manifest.md)

## History

- `01M76AZ94M8HWY6GCD06H8ES46` Release channel manifest and provenance: every v* release publishes a JSON manifest (channel, version, commit, per-platform archive URL and sha256) and GitHub build-provenance attestations alongside its archives and checksums
