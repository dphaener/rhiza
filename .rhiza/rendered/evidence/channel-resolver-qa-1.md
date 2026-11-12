# channel-resolver-qa-1

evidence · `01M77K7YVWF9NDWA2H42M9A225` · actor [claude-code-a983b6c0](../actor/claude-code-a983b6c0.md)

## result

PASS

## statement

QA at 04a4369f0: a go run program over the release.Resolver pointed at the real published channel (github.com/dphaener/rhiza) downloaded v0.1.0 for linux/amd64, linux/arm64, darwin/amd64 and darwin/arm64 and verified each sha256. The linux/amd64 sum d84c052b... matches the release's checksums.txt, and the extracted binary prints rhiza v0.1.0. linux/riscv64 is rejected with 'v0.1.0 has no archive for linux/riscv64'. A local manifest with no commit is rejected as malformed (exit 1). go test -trimpath ./internal/release/... passes.

## Links

### Out

- verifies [channel-resolver](../capability/rhiza/distribution/channel-resolver.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
