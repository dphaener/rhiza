# latest-channel-prerelease-qa-1

evidence · `01M77JCWBN9ZJW6CF0YJRV6BAY` · actor [claude-code-f8adba29](../actor/claude-code-f8adba29.md)

## result

PASS

## statement

QA at c6d4316, run as the Latest workflow does in a scratch clone, with goreleaser v2.20.1. RHIZA_BUILD_VERSION=v0.0.1-main.482+c6d4316; goreleaser check ok; release --snapshot --clean built the four rhiza_0.0.1-main.482+c6d4316_{linux,darwin}_{amd64,arm64}.tar.gz plus checksums.txt, and the workflow's four-archive loop passed. tools/release -channel latest -release-tag latest wrote manifest.json with channel latest, that version, the commit and four /releases/download/latest/ URLs whose sha256s match checksums.txt. The extracted linux_amd64 binary prints 'rhiza v0.0.1-main.482+c6d4316'. Tag path unbroken: with v0.1.0 tagged and RHIZA_BUILD_VERSION unset, goreleaser release --skip=publish built rhiza_0.1.0 archives and the binary prints v0.1.0. Ran the gate step's shell: no latest tag, publish; latest at a newer descendant on an old run, skip; same commit, publish; latest older, publish. Read latest.yml: push to main only, concurrency group latest with cancel-in-progress, and the release is created with --prerelease --latest=false. The earlier '+' caveat is resolved: github.com keeps '+' in asset names (bell-sw/Liberica's latest release has bellsoft-jdk27.0.1.0.3+5-*), and a download URL with a literal '+' returns 200. Not run: a real publish to GitHub, attestation, or gh release create.

## Links

### Out

- verifies [latest-channel-prerelease](../capability/rhiza/distribution/latest-channel-prerelease.md)

## History

- `01M76AZHGY605TXF68M27W4SZV` Latest channel: every merge to main builds the four archives as the next patch semver prerelease (vX.Y.Z+1-main.N+sha) and replaces one rolling latest GitHub prerelease with them, their checksums, manifest and attestations
