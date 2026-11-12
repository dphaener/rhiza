# readme-channels-qa-1

evidence · `01M7821S6DX2JGJE585HRYSASK` · actor [claude-code-3179a506](../actor/claude-code-3179a506.md)

## result

PASS

## statement

At d8d01e896 I followed the README install section by hand. Its install.sh one-liner and the --channel latest form both installed working binaries (v0.1.0 and v0.1.1-main.14+38cee74). Both manifest URLs it names resolve and are the ones the script fetches. sha256sum of the real rhiza_0.1.0_linux_amd64.tar.gz matched its manifest entry. gh attestation verify <archive> --repo dphaener/rhiza exited 0, with source repo dphaener/rhiza and signer release.yml@refs/tags/v0.1.0. It exited 1 for an altered copy. go install remains documented.

## Links

### Out

- verifies [readme-channels](../capability/rhiza/distribution/readme-channels.md)

## History

- `01M76B00QSZVGA0E6X85JE2BYY` install.sh: a curl | sh script that installs the newest release (or --channel latest) for this platform, verifying sha256 from the manifest, documented in the README
