# orchestrator-cuts-v010

decision · `01M76NPNRH6ZE9WNWKBYH5Z8PX`

## question

The owner cuts v0.1.0 by hand once the release manifest lands (the owner picked 'By hand after A'), so no leaf cuts it and no builder can verify a real attestation. After the owner pushes the v0.1.0 tag, who runs gh attestation verify on one of its archives and records it as Evidence on release-provenance, closing m33-release-channel's 'v0.1.0 is published' promise: the owner, or the orchestrator when it sees the tag?

## statement

The orchestrator cuts v0.1.0 once the manifest-and-provenance change set lands, then runs gh attestation verify on one of its archives and records the Evidence on @release-provenance

## wording

Orchestrator cuts and verifies: once the manifest-and-provenance change set lands, the orchestrator tags v0.1.0, watches the release workflow, runs gh attestation verify on an archive and records the Evidence; picked by the owner

## Links

### Out

- answers [v010-cut-and-attestation-evidence](../question/v010-cut-and-attestation-evidence.md)

## History

- `01M76AZ94M8HWY6GCD06H8ES46` Release channel manifest and provenance: every v* release publishes a JSON manifest (channel, version, commit, per-platform archive URL and sha256) and GitHub build-provenance attestations alongside its archives and checksums
