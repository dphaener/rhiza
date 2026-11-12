# release-provenance-verified

evidence · `01M77FT9KTD31100776T46BEPS` · actor [claude-code-0170ab75](../actor/claude-code-0170ab75.md)

## result

PASS

## statement

Verifier read commit 39091a90a (diff: .github/workflows/release.yml plus the builder's actor record only). release.yml now grants workflow permissions id-token: write and attestations: write (with contents: write), and adds an 'Attest build provenance' step using actions/attest-build-provenance@v3 with subject-path: dist/*.tar.gz, placed after 'Run GoReleaser' (which writes dist/) and before the manifest publish. Parsed the YAML with python yaml.safe_load: permissions = {contents: write, id-token: write, attestations: write}; steps = checkout, Set up Go, Run GoReleaser, Attest build provenance, Publish the channel manifest. goreleaser check / snapshot: the commit does not change .goreleaser.yaml or the release-snapshot job's steps, and release-snapshot.yml's paths filter includes .github/workflows/release.yml, so the job runs on this PR. goreleaser could not be run locally (the cached v2.20.0 needs go >= 1.27.3, the host has 1.25, and there is no network), so the job's pass on the PR is left to PR CI. The gh attestation verify check is out of scope here (moved to v010-cut-and-attestation-evidence).

## Links

### Out

- verifies [release-provenance](../capability/rhiza/distribution/release-provenance.md)

## History

- `01M76AZ94M8HWY6GCD06H8ES46` Release channel manifest and provenance: every v* release publishes a JSON manifest (channel, version, commit, per-platform archive URL and sha256) and GitHub build-provenance attestations alongside its archives and checksums
