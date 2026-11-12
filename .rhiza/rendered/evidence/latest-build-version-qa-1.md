# latest-build-version-qa-1

evidence · `01M77JCW6VCMCSKSSF5N35BKFJ` · actor [claude-code-f8adba29](../actor/claude-code-f8adba29.md)

## result

PASS

## statement

QA at c6d4316: go run ./tools/release version on this checkout (no v* tag) printed v0.0.1-main.482+c6d4316. In a scratch clone: with v0.1.0 tagged at HEAD it printed v0.1.0; one commit later, v0.1.1-main.1+99c5701 (after the tag, before v0.1.1). go test -trimpath -count=1 ./tools/release/... ./internal/release/... ok.

## Links

### Out

- verifies [latest-build-version](../capability/rhiza/distribution/latest-build-version.md)

## History

- `01M76AZHGY605TXF68M27W4SZV` Latest channel: every merge to main builds the four archives as the next patch semver prerelease (vX.Y.Z+1-main.N+sha) and replaces one rolling latest GitHub prerelease with them, their checksums, manifest and attestations
