# install-script

capability · `01M76B7DBYAAJRESBTJC5PHPB4` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/install -run TestInstallScript runs the repo-root install.sh against a local fixture channel (base URL override): it installs a working rhiza for the host platform, --channel latest picks the latest manifest, and a tampered archive fails with nothing installed; the same test runs sh -n install.sh, so make check covers the script's syntax with no extra tool

## statement

install.sh, run as curl -fsSL <raw URL> | sh, installs the newest release (or --channel latest) for the running linux or darwin on amd64 or arm64 into ~/.local/bin or --dir, verifying the archive sha256 against the channel manifest and failing on a mismatch

## Links

### In

- verifies [install-script-verified](../../../evidence/install-script-verified.md)
- verifies [install-script-qa-1](../../../evidence/install-script-qa-1.md)
- verifies [pr-491-install-sh](../../../evidence/pr-491-install-sh.md)

## History

- `01M76B00QSZVGA0E6X85JE2BYY` install.sh: a curl | sh script that installs the newest release (or --channel latest) for this platform, verifying sha256 from the manifest, documented in the README
