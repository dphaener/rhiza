# version-channel-check

capability · `01M76B7D5PY6PAKYK14QQCKZE1` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/version-channel, with the channel base URL pointed at a fixture (a test override, no network), covers newer, same and older, an unknown channel and an unreachable channel; describe-golden is regenerated

## statement

rhiza version --channel release|latest prints the newest version on that channel and whether it is newer than the running build, and writes nothing

## Links

### In

- verifies [version-channel-check-verified](../../../evidence/version-channel-check-verified.md)
- verifies [version-channel-check-qa-1](../../../evidence/version-channel-check-qa-1.md)
- verifies [version-channel-check-verified-2](../../../evidence/version-channel-check-verified-2.md)
- verifies [version-channel-check-qa-2](../../../evidence/version-channel-check-qa-2.md)
- verifies [pr-479-channel-resolver](../../../evidence/pr-479-channel-resolver.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
