# v010-cut-and-attestation-evidence

question · `01M76C5KP25JNKAG2A8RX0P656`

## statement

The owner cuts v0.1.0 by hand once the release manifest lands (the owner picked 'By hand after A'), so no leaf cuts it and no builder can verify a real attestation. After the owner pushes the v0.1.0 tag, who runs gh attestation verify on one of its archives and records it as Evidence on release-provenance, closing m33-release-channel's 'v0.1.0 is published' promise: the owner, or the orchestrator when it sees the tag?

## Links

### Out

- relates-to [release-provenance](../capability/rhiza/distribution/release-provenance.md) — The attestation check moved out of this leaf's acceptance (plan review round 2, f1): it can only run after v0.1.0 is cut by hand

### In

- answers [orchestrator-cuts-v010](../decision/orchestrator-cuts-v010.md)

## History

- `01M76AZ94M8HWY6GCD06H8ES46` Release channel manifest and provenance: every v* release publishes a JSON manifest (channel, version, commit, per-platform archive URL and sha256) and GitHub build-provenance attestations alongside its archives and checksums
