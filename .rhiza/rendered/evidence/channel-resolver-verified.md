# channel-resolver-verified

evidence · `01M77J2Y55YEW2TZX5T3FSC56D` · actor [claude-code-2d0ad57a](../actor/claude-code-2d0ad57a.md)

## result

PASS

## statement

Verifier at 9916356b3: go vet and go test -trimpath -count=1 -v ./internal/release/... pass (10 tests), also with HTTP(S)_PROXY pointed at a dead port. Resolver tests run only against an httptest server on localhost (Resolver.URLs and Client set to the fixture; no github.com URL fetched): TestDownloadEachPlatform downloads and checks the archive for all four Platforms (linux/darwin x amd64/arm64); TestNewestBothChannels reports v0.1.0 for release and v0.1.1-12-gabc for latest and rejects an unknown channel; TestDownloadRejectsBadSHA rejects a tampered archive on sha256; TestDownloadRejectsMissingPlatform rejects linux/riscv64 and a platform dropped from the manifest; TestRejectsMalformedManifest rejects non-JSON, empty, wrong-channel, short-sha and missing-url manifests. Diff read: resolve.go bounds reads, checks HTTP status, compares lowercase hex sha256 before returning the archive. gofmt clean; commit trailers and staged actor record correct.

## Links

### Out

- verifies [channel-resolver](../capability/rhiza/distribution/channel-resolver.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
