# version-channel-check-verified

evidence · `01M77K1EZ4PVJXEXHHXAPSM7KC` · actor [claude-code-177b63d6](../actor/claude-code-177b63d6.md)

## result

PASS

## statement

Verifier at 103ce69c9: go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(version-channel|describe-golden|version|usage)$' PASS. version-channel.txtar points RHIZA_TEST_CHANNEL_BASE_URL at an httptest file-server fixture (no network) and covers newer (v0.1.0 vs release v0.2.0), same (v0.2.0), older (v0.3.0), a dev build on latest, an unknown channel (usage error), a missing --channel value, an unreachable channel (127.0.0.1:1 errors cleanly), plain version unchanged, and asserts no .rhiza was written. describe-golden passes unchanged: no request/response type changed (only the version frame command's help summary, which describe does not emit), so the golden is current. internal/release tests PASS; internal/api -run 'Frame|Check|Help|Usage' PASS; go vet and gofmt clean. The diff reads the manifest via release.Resolver.Newest and prints channel/newest/newer (JSON via canon.Marshal) with no writes.

## Links

### Out

- verifies [version-channel-check](../capability/rhiza/distribution/version-channel-check.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
