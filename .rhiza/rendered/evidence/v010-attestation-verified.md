# v010-attestation-verified

evidence · `01M77JSGB3SM2V70783BW8HRE8` · actor [claude-code-orchestrator](../actor/claude-code-orchestrator.md)

## command

gh attestation verify rhiza_0.1.0_linux_amd64.tar.gz --repo dphaener/rhiza

## observed

2026-11-12, release https://github.com/dphaener/rhiza/releases/tag/v0.1.0, workflow run https://github.com/dphaener/rhiza/actions/runs/43194702318

## result

exit 0; SLSA provenance v1 (https://slsa.dev/provenance/v1) covering all four archives, signed by https://github.com/dphaener/rhiza/.github/workflows/release.yml@refs/tags/v0.1.0 at commit 00726452d61b13defb574d462ae3c2f14f8522ad; archive sha256 d84c052b37edbc0b6abbdd9c1208723eb8ed6454dd3c455c0dbcf9215fa431e9 matches manifest.json; the extracted binary reports rhiza v0.1.0

## statement

gh attestation verify rhiza_0.1.0_linux_amd64.tar.gz --repo dphaener/rhiza succeeds on the downloaded v0.1.0 release archive

## status

verified

## Links

### Out

- verifies [release-provenance](../capability/rhiza/distribution/release-provenance.md)

## History

- `01M77JSGDC8XGRZ8Q76E2PBER2` apply /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/ev-v010.json
