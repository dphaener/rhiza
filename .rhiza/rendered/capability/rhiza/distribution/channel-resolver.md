# channel-resolver

capability · `01M76B7D55NHFJP1CNA4H5HWDH` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./internal/release/... against an httptest server on localhost: picks each of the four platforms, reports both channels newest versions, and rejects a bad sha256, a platform missing from the manifest and a malformed manifest; nothing reaches github.com

## statement

internal/release resolves a channel (release or latest) to its manifest, reports its newest version, picks the running os/arch archive, downloads it and rejects it unless its sha256 matches the manifest

## Links

### In

- verifies [channel-resolver-verified](../../../evidence/channel-resolver-verified.md)
- verifies [channel-resolver-qa-1](../../../evidence/channel-resolver-qa-1.md)
- verifies [channel-resolver-qa-2](../../../evidence/channel-resolver-qa-2.md)
- verifies [pr-479-channel-resolver](../../../evidence/pr-479-channel-resolver.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
