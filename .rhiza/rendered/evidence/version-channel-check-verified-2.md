# version-channel-check-verified-2

evidence · `01M77KNVVMHB3GHAWMBP8SNMKV` · actor [claude-code-e838eb92](../actor/claude-code-e838eb92.md)

## result

PASS

## statement

At 940233379 (attempt 2) ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(version-channel|describe-golden)$|TestVersionNewer': all PASS. version-channel.txtar serves the channel manifests from a file-server fixture via RHIZA_TEST_CHANNEL_BASE_URL (no network) and covers newer, same, older, a dev build, an unknown channel (usage error), an unreachable channel (error, no crash) and that nothing is written; describe-golden passes against the current output. Read the fix diff: versionNewer now compares pre-release identifiers per semver (numeric ones numerically) and treats a -N-g<sha> describe suffix as N commits after its version; the new TestVersionNewer table covers rc.10 vs rc.2, main.10 vs main.9, v0.2.1-12-gdef vs v0.2.1-5-gabc, describe vs its tag and -dirty, resolving the confirmed f1 findings of the quality, qa and model checks. internal/release tests pass; gofmt and go vet clean on cmd/rhiza.

## Links

### Out

- verifies [version-channel-check](../capability/rhiza/distribution/version-channel-check.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
