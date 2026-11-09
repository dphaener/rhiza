# release-goreleaser-verified

evidence · `01M6Z6YZYF5EYBVV2SFMET0R2B` · actor [claude-code-2402d5bf](../actor/claude-code-2402d5bf.md)

## result

PASS

## statement

At 18512e7e, with GoReleaser v2.20.0 (go install, scratch dir) in a detached worktree: goreleaser check on .goreleaser.yaml exited 0 (1 configuration file validated); goreleaser release --snapshot --clean exited 0 and produced rhiza_0.0.0-SNAPSHOT-18512e7e_{linux,darwin}_{amd64,arm64}.tar.gz plus checksums.txt, and the workflow's own four-archive ls loop passed; the extracted linux_amd64 binary ran 'rhiza version' -> 'rhiza 0.0.0-SNAPSHOT-18512e7e' (ldflags wire main.version). Read the diff: .github/workflows/release-snapshot.yml runs check, release --snapshot --clean and the four-archive check on pull_request to main with paths limited to .goreleaser.yaml and the two release workflows; .github/workflows/release.yml triggers only on push tags ['v*'] (YAML parse confirms 'on' has only push), contents: write, runs goreleaser release --clean with GITHUB_TOKEN publishing to dphaener/rhiza; /dist/ is gitignored. Runner and action versions match ci.yml.

## Links

### Out

- verifies [release-goreleaser](../capability/rhiza/distribution/release-goreleaser.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
