# latest-channel-prerelease-verified

evidence · `01M77J0PG27D6G2HHSJKPEB438` · actor [claude-code-a583de19](../actor/claude-code-a583de19.md)

## result

PASS

## statement

Verified at bbdef16 (attempt 2; both attempt-1 findings fixed). Read latest.yml: on push to branches [main] only; concurrency group 'latest' with cancel-in-progress true; gate step skips publishing when tag latest exists, differs from GITHUB_SHA and is a descendant of it (git merge-base --is-ancestor). Ran the gate logic in a scratch repo: no tag -> publish, latest at a descendant -> skip, same commit -> publish, newer commit -> publish. Reproduced release-snapshot.yml in a scratch clone with goreleaser v2.20.1: RHIZA_BUILD_VERSION=$(go run ./tools/release version) = v0.0.1-main.480+bbdef16; goreleaser check ok (also with the env unset); release --snapshot --clean built rhiza_0.0.1-main.480+bbdef16_{linux,darwin}_{amd64,arm64}.tar.gz, the four-archive ls loop passed; go run ./tools/release -dist dist -tag $V -version ${V#v} -channel latest -release-tag latest exited 0, manifest has channel latest, version v0.0.1-main.480+bbdef16, 4 platforms with /releases/download/latest/ URLs; the built binary prints 'rhiza v0.0.1-main.480+bbdef16' (no vv). go test -trimpath ./internal/release/... ./tools/release/... ok. Caveat not covered by the acceptance: asset names contain '+' (semver build metadata); GitHub may rename '+' in uploaded asset names, which would break the manifest URLs; worth checking on the first real publish.

## Links

### Out

- verifies [latest-channel-prerelease](../capability/rhiza/distribution/latest-channel-prerelease.md)

## History

- `01M76AZHGY605TXF68M27W4SZV` Latest channel: every merge to main builds the four archives as the next patch semver prerelease (vX.Y.Z+1-main.N+sha) and replaces one rolling latest GitHub prerelease with them, their checksums, manifest and attestations
