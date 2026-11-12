# release-provenance

capability · `01M76B7CPZ26NZ03MPCK9S55PA` · parent [distribution](../distribution.md)

## acceptance

release.yml runs actions/attest-build-provenance over dist/*.tar.gz with id-token: write and attestations: write; goreleaser check and the release-snapshot job pass on the PR

## statement

The release job publishes a GitHub build-provenance attestation for every archive, so gh attestation verify <archive> --repo dphaener/rhiza succeeds on a downloaded release

## Links

### In

- relates-to [v010-cut-and-attestation-evidence](../../../question/v010-cut-and-attestation-evidence.md) — The attestation check moved out of this leaf's acceptance (plan review round 2, f1): it can only run after v0.1.0 is cut by hand
- verifies [release-provenance-verified](../../../evidence/release-provenance-verified.md)
- verifies [release-provenance-qa-1](../../../evidence/release-provenance-qa-1.md)
- verifies [pr-476-release-channel-manifest-and](../../../evidence/pr-476-release-channel-manifest-and.md)
- verifies [v010-attestation-verified](../../../evidence/v010-attestation-verified.md)

## History

- `01M76AZ94M8HWY6GCD06H8ES46` Release channel manifest and provenance: every v* release publishes a JSON manifest (channel, version, commit, per-platform archive URL and sha256) and GitHub build-provenance attestations alongside its archives and checksums
