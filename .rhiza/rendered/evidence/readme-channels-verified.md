# readme-channels-verified

evidence · `01M781TT8J4TD581GH7VCDE3CQ` · actor [claude-code-aacf15e7](../actor/claude-code-aacf15e7.md)

## result

PASS

## statement

At 0b0801d5b, I read the README diff from the merge-base. It names install.sh (curl ... install.sh | sh, --dir, --channel latest), both manifest URLs, releases/latest/download/manifest.json for release and releases/download/latest/manifest.json for latest, which match install.sh lines 47-48. It keeps go install and covers verification with sha256sum/shasum -a 256 and `gh attestation verify rhiza_<version>_<os>_<arch>.tar.gz --repo dphaener/rhiza`. Then I made a clean detached worktree at that commit and ran `rhiza check` there, with a binary built from the commit and RHIZA_ROOT/CHANGE/RUN unset. It printed ok and exited 0, with only pre-existing 'unclaimed' warnings, and no render failures. That fixes attempt 1's finding about the stale rendered actor claude-code-a6fbcd94.

## Links

### Out

- verifies [readme-channels](../capability/rhiza/distribution/readme-channels.md)

## History

- `01M76B00QSZVGA0E6X85JE2BYY` install.sh: a curl | sh script that installs the newest release (or --channel latest) for this platform, verifying sha256 from the manifest, documented in the README
