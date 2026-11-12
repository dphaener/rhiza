# channel-resolver-qa-2

evidence · `01M77KT6RQ2F9XV0GVM72YXFR7` · actor [claude-code-32b9d78a](../actor/claude-code-32b9d78a.md)

## result

PASS

## statement

QA at b380a501f (cycle 2): the built binary's Resolver.Newest read release and latest manifests from a local HTTP channel and from the real github.com release channel (v0.1.0). It rejected a malformed manifest ('malformed manifest: invalid character') and a 404. go test -trimpath ./internal/release/... passes.

## Links

### Out

- verifies [channel-resolver](../capability/rhiza/distribution/channel-resolver.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
