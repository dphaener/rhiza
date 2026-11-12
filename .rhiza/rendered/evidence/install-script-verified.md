# install-script-verified

evidence · `01M781N2D4JA691DZC16X1WM87` · actor [claude-code-5cf7c92f](../actor/claude-code-5cf7c92f.md)

## result

PASS

## statement

At 7c384f535 I ran go test -trimpath -count=1 -v ./tools/install -run TestInstallScript: all subtests pass, namely syntax (sh -n install.sh, under dash), release channel and latest channel (each installs a fixture rhiza from a local base-URL-override channel and the installed binary prints its channel), tampered archive (fails with a sha256 mismatch, and no rhiza is installed), piped with a compact manifest, piped help, and bad channel. I checked that both attempt-1 findings are fixed. f1: the manifest is now split with tr, not sed. With uname faked to Darwin/arm64 and a hand-written mixed-layout manifest served over local HTTP, sh install.sh --dir installed the darwin/arm64 archive and it ran. f2: cat install.sh | sh -s -- --help prints the heredoc help. I read the diff: the fixture marshals the real release.Manifest, and the release, latest and release-snapshot workflows publish manifest.json at the paths the script fetches.

## Links

### Out

- verifies [install-script](../capability/rhiza/distribution/install-script.md)

## History

- `01M76B00QSZVGA0E6X85JE2BYY` install.sh: a curl | sh script that installs the newest release (or --channel latest) for this platform, verifying sha256 from the manifest, documented in the README
