# install-script-qa-1

evidence · `01M7821S3T7W38ZFJ6HA8V425T` · actor [claude-code-3179a506](../actor/claude-code-3179a506.md)

## result

PASS

## statement

At d8d01e896 on linux/amd64 I ran the branch's install.sh hands-on against the real GitHub releases. cat install.sh | HOME=<scratch> sh fetched the release manifest and installed rhiza v0.1.0 into <scratch>/.local/bin (rhiza version: v0.1.0) and printed the PATH hint. cat install.sh | sh -s -- --channel latest --dir <scratch> installed v0.1.1-main.14+38cee74, which ran. For tamper, I served the real release manifest from a local mirror via RHIZA_INSTALL_BASE_URL with the linux/amd64 archive's bytes altered: it exited 1 with 'sha256 mismatch ... nothing installed', and the --dir stayed empty. sh -s -- --help printed the help under a pipe; --channel nope exited 1. go test -trimpath -count=1 ./tools/install -run TestInstallScript passed.

## Links

### Out

- verifies [install-script](../capability/rhiza/distribution/install-script.md)

## History

- `01M76B00QSZVGA0E6X85JE2BYY` install.sh: a curl | sh script that installs the newest release (or --channel latest) for this platform, verifying sha256 from the manifest, documented in the README
