# release-provenance-qa-1

evidence · `01M77GAJYRJMT0945B3RWE1QCW` · actor [claude-code-f1ace9df](../actor/claude-code-f1ace9df.md)

## result

PASS

## statement

QA at 55bf2f1d0: parsed .github/workflows/release.yml with yaml.safe_load: permissions are contents/id-token/attestations all write; steps are checkout, Set up Go, Run GoReleaser (release --clean), Attest build provenance (actions/attest-build-provenance@v3, subject-path dist/*.tar.gz), Publish the channel manifest. gh api confirms the actions/attest-build-provenance v3 tag exists. The glob dist/*.tar.gz matches exactly the four archives the manifest step produced in a GoReleaser-shaped dist built from real cross-compiled binaries. release-snapshot.yml parses and its path filter covers release.yml, internal/release and tools/release, so goreleaser check and the snapshot run on the PR. Not run: goreleaser (not installed) and a real attestation (needs a v* tag; the orchestrator verifies it after cutting v0.1.0 per orchestrator-cuts-v010).

## Links

### Out

- verifies [release-provenance](../capability/rhiza/distribution/release-provenance.md)

## History

- `01M76AZ94M8HWY6GCD06H8ES46` Release channel manifest and provenance: every v* release publishes a JSON manifest (channel, version, commit, per-platform archive URL and sha256) and GitHub build-provenance attestations alongside its archives and checksums
